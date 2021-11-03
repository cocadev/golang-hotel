package jac

import (
	"encoding/xml"
	"fmt"
	"net/url"
	repository "roomres/repository"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"
	roomresutils "../roomres/utils"

	//. "github.com/ahmetb/go-linq"
	"github.com/go-errors/errors"
)

/*
	{"Rooms":[{"adults":2,"children":0}],
	"Details":true,
	"ProviderId":6,
	"CheckIn":"2017-08-11",
	"CheckOut":"2017-08-12",
	"Lat":null,
	"Lon":null,
	"MinStarRating":null,
	"MaxStarRating":null,
	"HotelIds":["4566", "5472"],
	"CurrencyCode":"AUD"}
*/

type JacClient struct {
	HotelsPerBatch  int
	NumberOfThreads int
	SmartTimeout    int

	UserName        string
	SiteId          string
	ApiKey          string
	SecretKey       string
	ProfileCurrency string

	SearchEndPoint                          string
	PreBookEndPoint                         string
	BookEndPoint                            string
	BookingCancellationEndPoint             string
	BookingCancellationConfirmationEndPoint string
	PropertyDetailsEndpoint                 string

	CommissionProvider hbecommon.ICommissionProvider
	RepositoryFactory  *repository.RepositoryFactory

	Log roomresutils.ILog
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings, logging roomresutils.ILog) *JacClient {

	jacClient := &JacClient{

		Log: logging,

		SearchEndPoint:                          hotelProviderSettings.SearchEndPoint,
		PreBookEndPoint:                         hotelProviderSettings.PreBookEndPoint,
		BookEndPoint:                            hotelProviderSettings.BookingConfirmationEndPoint,
		BookingCancellationEndPoint:             hotelProviderSettings.BookingCancellationEndPoint,
		BookingCancellationConfirmationEndPoint: hotelProviderSettings.BookingCancellationConfirmationEndPoint,
		PropertyDetailsEndpoint:                 hotelProviderSettings.PropertyDetailsEndpoint,

		UserName:           hotelProviderSettings.UserName,
		SiteId:             hotelProviderSettings.SiteId,
		ApiKey:             hotelProviderSettings.ApiKey,
		SecretKey:          hotelProviderSettings.SecretKey,
		HotelsPerBatch:     hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads:    hotelProviderSettings.NumberOfThreads,
		SmartTimeout:       hotelProviderSettings.SmartTimeout,
		CommissionProvider: hotelProviderSettings.CommissionProvider,
		RepositoryFactory:  hotelProviderSettings.RepositoryFactory,
	}

	if strings.ToUpper(hotelProviderSettings.ProfileCurrency) == "AUD" {
		jacClient.ProfileCurrency = "5"
	} else if strings.ToUpper(hotelProviderSettings.ProfileCurrency) == "USD" {
		jacClient.ProfileCurrency = "2"
	}

	return jacClient
}

func (m *JacClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func (m *JacClient) CreateHttpSearchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate",
	}
}

func (m *JacClient) NewLoginDetails() *LoginDetails {

	return &LoginDetails{
		Login:          m.UserName,
		Password:       m.SecretKey,
		Locale:         "",
		AgentReference: "",
		CurrencyId:     m.ProfileCurrency, //"5", //AUD
	}
}

func (m *JacClient) NewSearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *SearchRequest {

	searchRequest := &SearchRequest{}

	searchRequest.LoginDetails = m.NewLoginDetails()

	return searchRequest
}

func (m *JacClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {

	hotelIdTransformer := NewHotelIdTransformer()

	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	if CheckMaxNumberOfRooms(hotelSearchRequest) && CheckMaxNumberOfPax(hotelSearchRequest) && hotelSearchRequest.CheckTreshold(12) {

		if len(hotelSearchRequest.HotelIds) == 1 && len(hotelSearchRequest.SpecificRoomRefs) > 0 {

			mapping_searchrequest_to_hbe_searchresponse(hotelSearchRequest, hotelSearchResponse, hotelIdTransformer)

		} else {
			request := m.NewSearchRequest(hotelSearchRequest)
			mapping_hbe_to_searchrequest(hotelSearchRequest, request, hotelIdTransformer)

			requests := m.SplitSearchRequests(request, hotelIdTransformer)
			response := m.SendSearchRequests(requests)

			AssignRequestedRoomsToSearchResult(request, response)

			mapping_searchresponse_to_hbe(
				response,
				hotelSearchRequest,
				hotelSearchResponse,
				hotelIdTransformer,
				m.Log.StartLogScope(roomresutils.LogScopeRef{}),
			)
		}

		if hotelSearchRequest.Details && len(hotelSearchRequest.SpecificRoomRefs) > 0 {

			m.PopulateDetails(hotelSearchRequest, hotelSearchResponse, hotelIdTransformer)
		}
	}

	return hotelSearchResponse
}

func (m *JacClient) RetreiveHotelContentItems(hotelIds []string) []*repository.ProviderHotelContent {

	hotelContentRepository := m.RepositoryFactory.CreateProviderHotelContentRepository()

	hotelContentItems := hotelContentRepository.GetHotelContentByProviderLanguageHotels(JacProviderId, JacStandardLanguage, hotelIds)

	return hotelContentItems
}

type HotelDetailsRequest struct {
	Index    int
	CheckIn  string
	Duration int
	Hotel    *hbecommon.Hotel
	RoomType *hbecommon.RoomType
}

func (m *JacClient) PopulateDetails(
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	hotelIdTransformer *HotelIdTransformer) {

	hotelDetailsRequests := []*HotelDetailsRequest{}
	uniqueHotelDetailsRequests := []*hbecommon.Hotel{}

	for _, hbeHotel := range hbeHotelSearchResponse.Hotels {

		for i, hbeRoomType := range hbeHotel.RoomTypes {

			hotelDetailsRequests = append(hotelDetailsRequests, &HotelDetailsRequest{
				CheckIn:  hotelSearchRequest.CheckIn,
				Duration: hotelSearchRequest.GetLos(),
				Hotel:    hbeHotel,
				RoomType: hbeRoomType,
				Index:    i,
			})
		}

		uniqueHotelDetailsRequests = append(uniqueHotelDetailsRequests, hbeHotel)
	}

	closableChannel := roomresutils.NewClosableChannel()

	//var requestChannel chan *HotelDetailsRequest = make(chan *HotelDetailsRequest)

	var responses []*HotelDetailsRequest

	var i int = 0
	var maxBatches int = m.NumberOfThreads * 2

	if len(hotelDetailsRequests) < maxBatches {
		maxBatches = len(hotelDetailsRequests)
	}

	for i < maxBatches {
		go m.SendHotelDetailsRequestChannel(closableChannel, hotelDetailsRequests[i])
		i++
	}

	/* get property details  */
	//propertyDetailsRequest := m.NewPropertyDetailsRequest()
	//propertyDetailsRequest.PropertyId = uniqueHotelDetailsRequests[0].CustomTag
	//propertyDetailsResponse := m.SendPropertyDetailsRequest(propertyDetailsRequest)

	propertyDetailsDescription := ""

	hotelContentItems := m.RetreiveHotelContentItems([]string{hotelIdTransformer.ExtractJacHotelId(uniqueHotelDetailsRequests[0].HotelId)})
	if len(hotelContentItems) > 0 && hotelContentItems[0].HotelContentObject != nil {
		propertyDetailsDescription = hotelContentItems[0].HotelContentObject.Remark
	}
	/* ********************* */

	for len(responses) < len(hotelDetailsRequests) {

		exit := false

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*HotelDetailsRequest))
		case <-time.After(time.Second * 30):
			exit = true
			closableChannel.Close()
			break
		}

		if exit {
			break
		}

		if i < len(hotelDetailsRequests) {
			go m.SendHotelDetailsRequestChannel(closableChannel, hotelDetailsRequests[i])
			i++
		}
	}

	descriptionItem := GetPropertyDescriptionItem(propertyDetailsDescription, "##Extra Information")

	for _, hotelDetailRequest := range hotelDetailsRequests {

		hotelDetailRequest.RoomType.Notes += "<br />" + descriptionItem
	}
}

func (m *JacClient) SplitSearchRequests(
	combinedSearchRequest *SearchRequest,
	hotelIdTransformer *HotelIdTransformer,
) (searchRequests []*SearchRequest) {

	searchRequestsByRegion := m.SplitSearchRequestsByRegion(combinedSearchRequest, hotelIdTransformer)

	for _, searchRequest := range searchRequestsByRegion {

		searchRequestsByBatch := m.SplitSearchRequestsByBatchSize(searchRequest, hotelIdTransformer)

		searchRequests = append(searchRequests, searchRequestsByBatch...)
	}

	return searchRequests
}

func (m *JacClient) SplitSearchRequestsByBatchSize(
	combinedSearchRequest *SearchRequest,
	hotelIdTransformer *HotelIdTransformer,
) (searchRequests []*SearchRequest) {

	if len(combinedSearchRequest.SearchDetails.PropertyReferenceIds) > m.HotelsPerBatch {

		searchRequest := combinedSearchRequest.Clone()
		searchRequest.SearchDetails.PropertyReferenceIds = []*PropertyReferenceId{}

		for _, propertyReferenceId := range combinedSearchRequest.SearchDetails.PropertyReferenceIds {

			if len(searchRequest.SearchDetails.PropertyReferenceIds) >= m.HotelsPerBatch {
				searchRequests = append(searchRequests, searchRequest)

				searchRequest = combinedSearchRequest.Clone()
				searchRequest.SearchDetails.PropertyReferenceIds = []*PropertyReferenceId{}
			}

			searchRequest.SearchDetails.PropertyReferenceIds = append(
				searchRequest.SearchDetails.PropertyReferenceIds,
				&PropertyReferenceId{ReferenceId: propertyReferenceId.ReferenceId})
		}

		searchRequests = append(searchRequests, searchRequest)

	} else {
		searchRequests = append(searchRequests, combinedSearchRequest)
	}

	return searchRequests
}

func (m *JacClient) SplitSearchRequestsByRegion(
	combinedSearchRequest *SearchRequest,
	hotelIdTransformer *HotelIdTransformer,
) (searchRequests []*SearchRequest) {

	if len(combinedSearchRequest.SearchDetails.PropertyReferenceIds) > 1 {

		searchRequestsByRegion := map[string]*SearchRequest{}

		for _, propertyReferenceId := range combinedSearchRequest.SearchDetails.PropertyReferenceIds {

			if info, ok := hotelIdTransformer.HotelIdInfoItems[propertyReferenceId.ReferenceId]; ok {

				if _, ok = searchRequestsByRegion[info.Region]; !ok {
					searchRequestsByRegion[info.Region] = combinedSearchRequest.Clone()
					searchRequestsByRegion[info.Region].SearchDetails.PropertyReferenceIds =
						[]*PropertyReferenceId{}
				}

				searchRequestsByRegion[info.Region].SearchDetails.PropertyReferenceIds =
					append(searchRequestsByRegion[info.Region].SearchDetails.PropertyReferenceIds,
						&PropertyReferenceId{ReferenceId: propertyReferenceId.ReferenceId})
			}
		}

		for _, searchRequest := range searchRequestsByRegion {
			searchRequests = append(searchRequests, searchRequest)
		}

	} else {

		searchRequests = append(searchRequests, combinedSearchRequest)
	}

	return searchRequests
}

func (m *JacClient) SendSearchRequests(searchRequests []*SearchRequest) *SearchResponse {

	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := roomresutils.NewClosableChannel()

	var responses []*SearchResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(searchRequests) < maxBatches {
		maxBatches = len(searchRequests)
	}

	for i < maxBatches {
		go m.SendSearchRequestChannel(closableChannel, searchRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.SmartTimeout

	if len(searchRequests) > 0 && searchRequests[0].Details {
		smartTimeout = 30
	}

	for len(responses) < len(searchRequests) {

		remaining := float64(smartTimeout) - totalTime.Seconds()
		time1 := time.Now()

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*SearchResponse))
		case <-time.After(time.Duration(float64(time.Second) * remaining)):
			break
		}

		time2 := time.Now()
		duration := time2.Sub(time1)

		totalTime = time.Duration(float64(time.Second) * (totalTime.Seconds() + duration.Seconds()))

		if totalTime.Seconds() >= float64(smartTimeout) {
			closableChannel.Close()
			break
		}

		if i < len(searchRequests) {
			go m.SendSearchRequestChannel(closableChannel, searchRequests[i])
			i++
		}
	}

	if len(responses) < len(searchRequests) {
		m.Log.LogEvent(roomresutils.EventTypeInfo2,
			"JAC SendSearchRequest ",
			fmt.Sprintf("%d - %d", len(searchRequests), len(responses)),
		)
	}

	combinedSearchResponse := &SearchResponse{PropertyResults: []*PropertyResult{}}

	for _, response := range responses {

		if response != nil && response.PropertyResults != nil {
			combinedSearchResponse.PropertyResults =
				append(
					combinedSearchResponse.PropertyResults,
					response.PropertyResults...,
				)
		}
	}

	return combinedSearchResponse
}

func (m *JacClient) SendSearchRequestChannel(closableChannel *roomresutils.ClosableChannel, searchRequest *SearchRequest) {

	defer func() {
		if err := recover(); err != nil {

			m.Log.LogEvent(roomresutils.EventTypeError,
				"JAC SendSearchRequest ",
				fmt.Sprintf("%s", errors.Wrap(err, 2).ErrorStack()),
			)

			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	searchResponse := m.SendSearchRequest(searchRequest)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}

func (m *JacClient) SendSearchRequest(searchRequest *SearchRequest) *SearchResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(searchRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - serialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo3, "JAC SendSearchRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
		})

	requestData = nil

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "JAC SendSearchRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - request", httpResponse.Err.Error())
	}

	if len(httpResponse.ResponseBody) > 500000 {

		var refs []string

		if searchRequest.SearchDetails != nil {
			for _, ref := range searchRequest.SearchDetails.PropertyReferenceIds {

				refs = append(refs, ref.ReferenceId)
			}
		}

		logScope.LogEvent(roomresutils.EventTypeInfo2, "JAC SendSearchRequest - LARGE RESPONSE DETECTED ", fmt.Sprintf("%d - %v", len(httpResponse.ResponseBody), refs))
	}

	var searchResponse *SearchResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - deserialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendSearchRequest - deserialization", errResponse.Error())
	}

	return searchResponse
}

func (m *JacClient) NewPreBookRequest() *PreBookRequest {

	preBookRequest := &PreBookRequest{}

	preBookRequest.LoginDetails = m.NewLoginDetails()

	return preBookRequest
}

func (m *JacClient) SendHotelDetailsRequestChannel(closableChannel *roomresutils.ClosableChannel, hotelDetailsRequest *HotelDetailsRequest) {

	defer func() {
		if err := recover(); err != nil {

			m.Log.LogEvent(roomresutils.EventTypeError,
				"JAC HotelDetailsRequest ",
				fmt.Sprintf("%s", errors.Wrap(err, 2).ErrorStack()),
			)

			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	preBookRequest := m.NewPreBookRequest()

	var combinedRoomRef CombinedRoomRef

	DecodeCombinedRoomRef(hotelDetailsRequest.RoomType.Ref, &combinedRoomRef)

	mapping_hbe_to_prebookrequest(hotelDetailsRequest.CheckIn, hotelDetailsRequest.Duration, &combinedRoomRef, preBookRequest)

	preBookRequest.Index = hotelDetailsRequest.Index

	preBookResponse := m.SendPreBookRequest(preBookRequest, false)

	mapping_prebookresponse_to_hbe(m.CommissionProvider, preBookResponse, hotelDetailsRequest.RoomType)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- hotelDetailsRequest
	})
}

func (m *JacClient) SendPreBookRequest(preBookRequest *PreBookRequest, allowLogging bool) *PreBookResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.PreBookEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(preBookRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - serialization", fmt.Sprintf("%v", preBookRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if allowLogging && logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPreBookRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
			Timeout:              time.Duration(10 * time.Second),
		})

	if allowLogging && logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPreBookRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - request", httpResponse.Err.Error())
	}

	var preBookResponse *PreBookResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &preBookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - deserialization", fmt.Sprintf("%v", preBookRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreBookRequest - deserialization", errResponse.Error())
	}

	return preBookResponse
}

func (m *JacClient) NewPropertyDetailsRequest() *PropertyDetailsRequest {

	propertyDetailsRequest := &PropertyDetailsRequest{}

	propertyDetailsRequest.LoginDetails = m.NewLoginDetails()

	return propertyDetailsRequest
}

func (m *JacClient) SendPropertyDetailsRequest(propertyDetailsRequest *PropertyDetailsRequest) *PropertyDetailsResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.PropertyDetailsEndpoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(propertyDetailsRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - serialization", fmt.Sprintf("%v", propertyDetailsRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPropertyDetailsRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
			Timeout:              time.Duration(10 * time.Second),
		})

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPropertyDetailsRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - request", httpResponse.Err.Error())
	}

	var propertyDetailsResponse *PropertyDetailsResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &propertyDetailsResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - deserialization", fmt.Sprintf("%v", propertyDetailsRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPropertyDetailsRequest - deserialization", errResponse.Error())
	}

	return propertyDetailsResponse
}

func (m *JacClient) NewPreCancelRequest() *PreCancelRequest {

	preCancelRequest := &PreCancelRequest{}

	preCancelRequest.LoginDetails = m.NewLoginDetails()

	return preCancelRequest
}

func (m *JacClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	preCancelRequest := m.NewPreCancelRequest()

	mapping_hbe_to_precancelrequest(bookingCancelRequest, preCancelRequest)

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.BookingCancellationEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(preCancelRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - serialization", fmt.Sprintf("%v", preCancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPreCancellationRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
			Timeout:              time.Duration(180 * time.Second),
		})

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendPreCancellationRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - request", httpResponse.Err.Error())
	}

	var preCancelResponse *PreCancelResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &preCancelResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - deserialization", fmt.Sprintf("%v", preCancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendPreCancellationRequest - deserialization", errResponse.Error())
	}

	hbeBookingCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_precancelresponse_to_hbe(preCancelResponse, bookingCancelRequest, hbeBookingCancelResponse)

	return hbeBookingCancelResponse
}

func (m *JacClient) NewCancelRequest() *CancelRequest {

	cancelRequest := &CancelRequest{}

	cancelRequest.LoginDetails = m.NewLoginDetails()

	return cancelRequest
}

func (m *JacClient) CancelBookingConfirm(bookingCancelConfirmationRequest *hbecommon.BookingCancelConfirmationRequest) *hbecommon.BookingCancelConfirmationResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	cancelRequest := m.NewCancelRequest()

	mapping_hbe_to_cancelrequest(bookingCancelConfirmationRequest, cancelRequest)

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.BookingCancellationConfirmationEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(cancelRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - serialization", fmt.Sprintf("%v", cancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendCancellationConfirmationRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
			Timeout:              time.Duration(180 * time.Second),
		})

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC SendCancellationConfirmationRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - request", httpResponse.Err.Error())
	}

	var cancelResponse *CancelResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &cancelResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - deserialization", fmt.Sprintf("%v", cancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC SendCancellationConfirmationRequest - deserialization", errResponse.Error())
	}

	hbeBookingCancelConfirmationResponse := &hbecommon.BookingCancelConfirmationResponse{}

	mapping_cancelresponse_to_hbe(cancelResponse, hbeBookingCancelConfirmationResponse)

	return hbeBookingCancelConfirmationResponse
}

func (m *JacClient) SendSpecialRequest(bookingSpecialRequestRequest *hbecommon.BookingSpecialRequestRequest) *hbecommon.BookingSpecialRequestResponse {
	return nil
}

func (m *JacClient) NewBookRequest() *BookRequest {

	bookRequest := &BookRequest{}

	bookRequest.LoginDetails = m.NewLoginDetails()

	return bookRequest
}

func (m *JacClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {

	/* error similation  */
	/*
		hbeBookingResponseSim := &hbecommon.BookingResponse{}

		hbeBookingResponseSim.BookingStatus = hbecommon.BookingStatusFailedRestartEnum
		hbeBookingResponseSim.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "Error Similation"}}

		return hbeBookingResponseSim
	*/
	/* error similation  */

	preBookRequest := m.NewPreBookRequest()

	mapping_hbe_bookingrequest_to_prebookrequest(hbeBookingRequest, preBookRequest)
	preBookResponse := m.SendPreBookRequest(preBookRequest, true)

	if (preBookResponse.TotalPrice-hbeBookingRequest.Total)/hbeBookingRequest.Total > 0.05 {

		bookingResponse := &hbecommon.BookingResponse{}

		bookingResponse.BookingStatus = hbecommon.BookingStatusFailedRestartEnum
		bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "The hotel room rate has changed since your selection. Please confirm that you are willing to accept the new rate and re-submit the booking."}}

		return bookingResponse
	}

	bookRequest := m.NewBookRequest()

	mapping_hbe_to_bookingrequest(hbeBookingRequest, preBookResponse, bookRequest)

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.BookEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(bookRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - serialization", fmt.Sprintf("%v", bookRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - serialization", err.Error())
	}

	requestData = []byte(fmt.Sprintf("Data=%s", url.QueryEscape(string(requestData))))

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC BookRequest - request data", fmt.Sprintf("%s", body))
	}

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
			Timeout:              time.Duration(180 * time.Second),
		})

	if logScope.AllowLogging(roomresutils.EventTypeInfo1) {

		logScope.LogEvent(roomresutils.EventTypeInfo1, "JAC BookRequest - response data", fmt.Sprintf("%s", httpResponse.ResponseBody))
	}

	if httpResponse.Err != nil {
		body, _ := url.QueryUnescape(string(requestData))

		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - request", fmt.Sprintf("%s", body))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - request", httpResponse.Err.Error())
	}

	var bookResponse *BookResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &bookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - deserialization", fmt.Sprintf("%v", bookRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, "JAC BookRequest - deserialization", errResponse.Error())
	}

	hbeBookingResponse := &hbecommon.BookingResponse{}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingResponse)

	return hbeBookingResponse
}

func (m *JacClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {

	preBookRequest := m.NewPreBookRequest()

	mapping_hbe_prebookrequest(roomCxlRequest.RoomRef, preBookRequest)

	preBookResponse := m.SendPreBookRequest(preBookRequest, false)

	roomCxlResponse := &hbecommon.RoomCxlResponse{}

	mapping_prebookresponse_to_roomcxlresponse(preBookResponse, roomCxlResponse)

	return roomCxlResponse
}

func (m *JacClient) GetBookingInfo(request *hbecommon.BookingInfoRequest) *hbecommon.BookingInfoResponse {
	return &hbecommon.BookingInfoResponse{HotelPhoneNumber: ""}
}

func (m *JacClient) ListBookingReport(request *hbecommon.BookingReportRequest) *hbecommon.BookingReportResponse {
	return nil
}

func (m *JacClient) Init() {

}
