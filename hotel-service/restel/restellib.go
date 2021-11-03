package restel

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"
	"golang.org/x/net/html/charset"
	//. "github.com/ahmetb/go-linq"
)

type RestelClient struct {
	Settings *RestelSettings

	HotelsPerBatch  int
	NumberOfThreads int
	SmartTimeout    int

	ProfileCurrency string //TODO:remove?

	SearchEndPoint                          string
	PreBookEndPoint                         string
	BookEndPoint                            string
	BookingCancellationEndPoint             string
	BookingCancellationConfirmationEndPoint string
	BookingInfoEndpoint                     string
	CxlFeeDetailEndpoint                    string

	CommissionProvider hbecommon.ICommissionProvider
	CreditCardInfo     *hbecommon.CreditCardInfo

	//RepositoryFactory  *repository.RepositoryFactory

	//Log roomresutils.ILog
}

type RestelSettings struct {
	BaseEndPoint      string `json:"BaseEndPoint"`
	UserCode          string `json:"UserCode"`
	UserPassword      string `json:"UserPassword"`
	Client            string `json:"Client"`
	AccessCode        string `json:"AccessCode"`
	AgencyAffiliation string `json:"AgencyAffiliation"`
	AgencyUserCode    string `json:"AgencyUserCode"`
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings) *RestelClient {

	restelClient := &RestelClient{

		//Log: logging,
		SearchEndPoint:                          "110",
		PreBookEndPoint:                         "202",
		BookEndPoint:                            "3",
		BookingCancellationEndPoint:             "3",
		BookingCancellationConfirmationEndPoint: "401",
		BookingInfoEndpoint:                     "15",
		CxlFeeDetailEndpoint:                    "144",

		HotelsPerBatch:     hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads:    hotelProviderSettings.NumberOfThreads,
		SmartTimeout:       hotelProviderSettings.SmartTimeout,
		CommissionProvider: hotelProviderSettings.CommissionProvider,
		CreditCardInfo:     hotelProviderSettings.CreditCardInfo,
		//RepositoryFactory:  hotelProviderSettings.RepositoryFactory,
	}

	var restelClientSettings RestelSettings
	err := json.Unmarshal([]byte(hotelProviderSettings.Metadata), &restelClientSettings)
	if err != nil {
		fmt.Printf("Error : Restel settings : %s", hotelProviderSettings.Metadata)
	}
	restelClient.Settings = &restelClientSettings

	if strings.ToUpper(hotelProviderSettings.ProfileCurrency) == "AUD" {
		restelClient.ProfileCurrency = "5"
	} else if strings.ToUpper(hotelProviderSettings.ProfileCurrency) == "USD" {
		restelClient.ProfileCurrency = "2"
	}

	return restelClient
}

func (m *RestelClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func (m *RestelClient) CreateHttpBodyParameters() map[string]string {
	return map[string]string{
		"clausu":    m.Settings.UserPassword,
		"codigousu": m.Settings.UserCode,
		"afiliacio": m.Settings.Client,
		"secacc":    m.Settings.AccessCode,
	}
}

func (m *RestelClient) CreateHttpSearchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate",
	}
}

//Plugin APIs
func (m *RestelClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {
	//fmt.Printf("search request data = %+v\n", hotelSearchRequest)

	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	if CheckMaxNumberOfRooms(hotelSearchRequest) && CheckMaxNumberOfPax(hotelSearchRequest) && hotelSearchRequest.CheckTreshold(12) {
		requests := mapping_hbe_to_searchrequest(hotelSearchRequest)
		for _, req := range requests {
			req.ServiceNumber = m.SearchEndPoint
			req.RestelSearchParam.Affiliation = m.Settings.AgencyAffiliation
			req.RestelSearchParam.UserCode = m.Settings.AgencyUserCode
		}
		requests = GroupSearchRequestsByCity(requests)
		requests = SplitBatchSearchRequests(requests, m.HotelsPerBatch)
		//fmt.Printf("multiple search requests count : %d\n", len(requests))

		response := m.SendSearchRequests(requests)

		//fmt.Printf("Result hotel counts = %d\n", len(response.ParamResult.HotelResults.Hotels))
		roomCxlResponses := []*CancelFeeDetailResponse{}

		if hotelSearchRequest.Details && len(hotelSearchRequest.SpecificRoomRefs) > 0 {
			var searchRoomRefs []*ShortRoomRef

			for _, roomRefStr := range hotelSearchRequest.SpecificRoomRefs {
				var shortRoomRef ShortRoomRef
				err := json.Unmarshal([]byte(roomRefStr), &shortRoomRef)
				if err != nil {
					fmt.Printf("Cannot parse specialRoomRef = %s\n", roomRefStr)
				} else {
					searchRoomRefs = append(searchRoomRefs, &shortRoomRef)
				}
			}

			if response.ParamResult.HotelResults != nil && len(response.ParamResult.HotelResults.Hotels) > 0 {
				//call room cxl policy request only when case 3
				hotel := response.ParamResult.HotelResults.Hotels[0]

				var roomPlans []*RoomPlan
				for _, roomPlan := range hotel.HotelRestrict.RoomPlans {
					if len(searchRoomRefs) > 0 {
						//should check specifiedRoomRefs
						isMatch := false
						for _, searchRoomRef := range searchRoomRefs {
							if searchRoomRef.AdultChildren != "" && hotel.HotelRestrict.AdultChldren != searchRoomRef.AdultChildren {
								continue
							}
							if searchRoomRef.RoomType != "" && roomPlan.RoomType != searchRoomRef.RoomType {
								continue
							}
							if searchRoomRef.MealPlanType != "" && roomPlan.RoomLine.MealPlanType != searchRoomRef.MealPlanType {
								continue
							}
							if searchRoomRef.RefundableType != "" && roomPlan.RoomLine.FeeRate != searchRoomRef.RefundableType {
								continue
							}

							if searchRoomRef.MinPrice != "" && ConvertStringToFloat32(roomPlan.RoomLine.PlanPrice) < ConvertStringToFloat32(searchRoomRef.MinPrice) {
								continue
							}

							if searchRoomRef.MaxPrice != "" && ConvertStringToFloat32(roomPlan.RoomLine.PlanPrice) > ConvertStringToFloat32(searchRoomRef.MaxPrice) {
								continue
							}
							isMatch = true
							break
						}

						if !isMatch {
							continue
						}

						roomPlans = append(roomPlans, roomPlan)
					}
				}

				hotel.HotelRestrict.RoomPlans = roomPlans

				var cancelFeeRequests []*CancelFeeDetailRequest
				for _, roomPlan := range hotel.HotelRestrict.RoomPlans {
					fmt.Printf("Hotel=%s, plan=%s\n", hotel.HotelCode, roomPlan.RoomLine.CompressedAvailability)
					cancelFeeDetailRequest := &CancelFeeDetailRequest{
						CancelFeeDetailRequestParam: &CancelFeeDetailRequestParam{
							Hotel:       hotel.HotelCode,
							BookingLine: roomPlan.RoomLine.CompressedAvailability,
							Lang:        2,
						},
					}
					cancelFeeDetailRequest.ServiceNumber = m.CxlFeeDetailEndpoint
					cancelFeeRequests = append(cancelFeeRequests, cancelFeeDetailRequest)
				}
				roomCxlResponses = m.SendCancelFeeDetailRequests(cancelFeeRequests)
			}
		}

		mapping_searchresponse_to_hbe(
			response,
			hotelSearchRequest,
			hotelSearchResponse,
			roomCxlResponses,
			//m.Log.StartLogScope(roomresutils.LogScopeRef{}),
		)
	}

	return hotelSearchResponse
}

func (m *RestelClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {

	//set credit card information from restel client information
	hbeBookingRequest.CreditCardInfo = m.CreditCardInfo

	preBookRequests := mapping_hbe_bookingrequest_to_prebookrequest(hbeBookingRequest)
	fmt.Printf("Request counts = %d\n", len(preBookRequests))

	if len(preBookRequests) > 1 {
		fmt.Printf("Cannot booking multi rooms. Please try with one room.")
		return &hbecommon.BookingResponse{
			BookingStatus: hbecommon.BookingStatusFailedEnum,
			ErrorMessages: []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{Message: "Cannot booking multi rooms. Please try with one room."}},
		}
	}

	for _, req := range preBookRequests {
		req.ServiceNumber = m.PreBookEndPoint
	}

	preBookResponses := m.SendPreBookRequests(preBookRequests)

	//fmt.Printf("pre booking responses : %d\n", len(preBookResponses))
	preBookResponse := preBookResponses[0]
	//fmt.Printf("pre booking response : %+v\n", preBookResponse)
	if preBookResponse != nil && preBookResponse.ParamPreBookResult.ReservationStatus == "00" {
		//Success Booking reservation
		if (preBookResponse.ParamPreBookResult.TotalReservationAmount-hbeBookingRequest.Total)/hbeBookingRequest.Total > 0.05 {
			return &hbecommon.BookingResponse{
				BookingStatus: hbecommon.BookingStatusFailedRestartEnum,
				ErrorMessages: []*hbecommon.ErrorMessage{
					&hbecommon.ErrorMessage{
						Message: "The hotel room rate has changed since your selection. Please confirm that you are willing to accept the new rate and re-submit the booking.",
					},
				},
			}
		}

		bookRequest := &ConfirmBookRequest{}
		bookRequest.ServiceNumber = m.BookEndPoint
		makeBookingConfirmRequest(hbeBookingRequest, preBookResponse, bookRequest)

		bookResponse := m.SendBookRequest(bookRequest)

		hbeBookingResponse := &hbecommon.BookingResponse{}

		var cancelFeeRequests []*CancelFeeDetailRequest
		for _, req := range preBookRequests {
			bookParam := req.RestelBookingParam

			cancelFeeDetailRequest := &CancelFeeDetailRequest{
				CancelFeeDetailRequestParam: &CancelFeeDetailRequestParam{
					Hotel:       bookParam.HotelCode,
					BookingLine: bookParam.SimpleLines.CompressedAvailability,
					Lang:        2,
				},
			}
			cancelFeeDetailRequest.ServiceNumber = m.CxlFeeDetailEndpoint
			cancelFeeRequests = append(cancelFeeRequests, cancelFeeDetailRequest)
		}
		roomCxlResponses := m.SendCancelFeeDetailRequests(cancelFeeRequests)

		mapping_bookresponse_to_hbe(bookResponse, preBookResponse, hbeBookingResponse, roomCxlResponses, hbeBookingRequest)

		return hbeBookingResponse
	} else {
		//Faile to booking reservation
		return &hbecommon.BookingResponse{
			BookingStatus: hbecommon.BookingStatusFailedEnum,
			ErrorMessages: []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{Message: preBookResponse.ParamPreBookResult.Reservado2}},
		}
	}
}

func (m *RestelClient) CancelBooking(hbeBookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {

	bookCancelRequest := &ConfirmBookRequest{}
	bookCancelRequest.ServiceNumber = m.BookEndPoint
	bookCancelRequest.RestelConfirmBookingParam = &RestelConfirmBookingParam{
		BookingNumber: hbeBookingCancelRequest.Ref,
		Action:        "AI",
	}

	bookCancelResponse := m.SendBookRequest(bookCancelRequest)

	hbeBookingCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_precancelresponse_to_hbe(bookCancelResponse, hbeBookingCancelResponse)

	return hbeBookingCancelResponse
}

func (m *RestelClient) CancelBookingConfirm(hbeBookingCancelConfirmationRequest *hbecommon.BookingCancelConfirmationRequest) *hbecommon.BookingCancelConfirmationResponse {

	bookCancelRequest := &CancelConfirmedBookRequest{}
	bookCancelRequest.ServiceNumber = m.BookingCancellationConfirmationEndPoint
	bookCancelRequest.CancelConfirmedBookRequestParam = &CancelConfirmedBookRequestParam{
		BookingNumber:      hbeBookingCancelConfirmationRequest.Ref,
		ShortBookingNumber: hbeBookingCancelConfirmationRequest.InternalRef,
	}

	bookCancelResponse := m.SendCancelBookingConfirmRequest(bookCancelRequest)

	hbeBookingCancelConfirmationResponse := &hbecommon.BookingCancelConfirmationResponse{}

	mapping_cancelresponse_to_hbe(bookCancelResponse, hbeBookingCancelConfirmationResponse)

	return hbeBookingCancelConfirmationResponse
}

func (m *RestelClient) GetBookingInfo(request *hbecommon.BookingInfoRequest) *hbecommon.BookingInfoResponse {

	bookingInfoRequest := &BookingInfoRequest{}
	bookingInfoRequest.ServiceNumber = m.BookingInfoEndpoint
	bookingInfoRequest.BookingInfoRequestParam = &BookingInfoRequestParam{
		HotelCode: "626886", //TODO:hard code temporary
		Language:  "2",      //english
	}

	bookInfoResponse := m.SendBookingInfoRequest(bookingInfoRequest)

	hbeBookingInfoResponse := &hbecommon.BookingInfoResponse{}

	if bookInfoResponse.BookingInfoDetail != nil {
		hbeBookingInfoResponse.HotelSupplierId = "" //not sure what could mapping here
		hbeBookingInfoResponse.HotelPhoneNumber = bookInfoResponse.BookingInfoDetail.Phone
	}

	return hbeBookingInfoResponse
}

func (m *RestelClient) SendSpecialRequest(bookingSpecialRequestRequest *hbecommon.BookingSpecialRequestRequest) *hbecommon.BookingSpecialRequestResponse {
	return nil
}

func (m *RestelClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {
	request := &CancelFeeDetailRequest{}
	request.ServiceNumber = m.CxlFeeDetailEndpoint

	//roomCxlRequest to request
	var ref RoomPlan
	json.Unmarshal([]byte(roomCxlRequest.RoomRef), &ref)

	//fmt.Printf("BookingLines = %s\n", ref.RoomLine.CompressedAvailability)

	request.CancelFeeDetailRequestParam = &CancelFeeDetailRequestParam{
		Hotel:       roomCxlRequest.HotelId,
		BookingLine: ref.RoomLine.CompressedAvailability,
		Lang:        2,
	}

	response := m.SendCancelFeeDetailRequest(request)

	hbeResponse := &hbecommon.RoomCxlResponse{}

	compressData := ref.RoomLine.CompressedAvailability[0]
	params := strings.Split(compressData, "#")

	checkInDate, _ := time.Parse("20060102", params[7])

	mapping_cancelfeedetail_to_hbe(response, hbeResponse, checkInDate.Format("2006-01-02"))

	return hbeResponse

	return nil
}

func (m *RestelClient) ListBookingReport(request *hbecommon.BookingReportRequest) *hbecommon.BookingReportResponse {
	return nil
}

func (m *RestelClient) Init() {

}

//local functions
func (m *RestelClient) SendSearchRequests(searchRequests []*SearchRequest) *SearchResponse {

	//fmt.Printf("SendSearchRequests - Start %d\n", len(searchRequests))
	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := NewClosableChannel()

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

	//fmt.Printf("RESTEL SendSearchRequest %d - %d\n", len(searchRequests), len(responses))

	combinedSearchResponse := &SearchResponse{
		ParamResult: &ParamResult{
			HotelResults: &HotelResults{
				Hotels: []*HotelResult{},
			},
		},
	}

	for _, response := range responses {

		if response != nil && response.ParamResult.HotelResults != nil {
			combinedSearchResponse.ParamResult.HotelResults.Hotels =
				append(
					combinedSearchResponse.ParamResult.HotelResults.Hotels,
					response.ParamResult.HotelResults.Hotels...,
				)
		}
	}

	return combinedSearchResponse
}

func (m *RestelClient) SendSearchRequestChannel(closableChannel *ClosableChannel, searchRequest *SearchRequest) {

	//("SendSearchRequestChannel : %+v, %+v\n", closableChannel, searchRequest)

	defer func() {
		if err := recover(); err != nil {

			// m.Log.LogEvent(roomresutils.EventTypeError,
			// 	"RESTEL SendSearchRequest ",
			// 	fmt.Sprintf("%s", errors.Wrap(err, 2).ErrorStack()),
			// )

			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	searchResponse := m.SendSearchRequest(searchRequest)

	//HotelID update with request one
	for _, hotel := range searchResponse.ParamResult.HotelResults.Hotels {
		hotel.HotelCode = "rtl_" + searchRequest.RestelSearchParam.CompoundHotelId
	}

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}

func (m *RestelClient) SendSearchRequest(searchRequest *SearchRequest) *SearchResponse {
	fmt.Println("SendSearchRequest ----------- Start ----------")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(searchRequest)

	//fmt.Printf("%+v\n", searchRequest)
	//fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var searchResponse *SearchResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&searchResponse)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	//fmt.Printf("Response = %+v\n", searchResponse)

	return searchResponse
}

func (m *RestelClient) SendPreBookRequest(preBookRequest *PreBookRequest, allowLogging bool) *PreBookResponse {
	fmt.Println("SendPreBookRequest ----------- Start ----------")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(preBookRequest)

	//fmt.Printf("%+v\n", searchRequest)
	fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var preBookResponse *PreBookResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&preBookResponse)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	//fmt.Printf("Response = %+v\n", searchResponse)

	return preBookResponse
}

func (m *RestelClient) SendBookRequest(bookRequest *ConfirmBookRequest) *ConfirmBookResponse {
	fmt.Println("Call SendBookRequest")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(bookRequest)

	//fmt.Printf("%+v\n", searchRequest)
	fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var bookResponse *ConfirmBookResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&bookResponse)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	return bookResponse
}

func (m *RestelClient) SendCancelBookingConfirmRequest(bookCancelRequest *CancelConfirmedBookRequest) *CancelConfirmedBookResponse {
	fmt.Println("Call SendCancelBookingConfirmRequest")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(bookCancelRequest)

	//fmt.Printf("%+v\n", searchRequest)
	fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var bookResponse *CancelConfirmedBookResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&bookResponse)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	return bookResponse
}

func (m *RestelClient) SendBookingInfoRequest(bookInfoRequest *BookingInfoRequest) *BookingInfoResponse {
	fmt.Println("Call SendBookingInfoRequest")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(bookInfoRequest)

	//fmt.Printf("%+v\n", searchRequest)
	fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var bookResponse *BookingInfoResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&bookResponse)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	return bookResponse
}

func (m *RestelClient) SendPreBookRequests(preBookRequests []*PreBookRequest) []*PreBookResponse {

	//fmt.Printf("SendSearchRequests - Start %d\n", len(searchRequests))
	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := NewClosableChannel()

	var responses []*PreBookResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(preBookRequests) < maxBatches {
		maxBatches = len(preBookRequests)
	}

	for i < maxBatches {
		go m.SendPreBookRequestChannel(closableChannel, preBookRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.SmartTimeout

	if len(preBookRequests) > 0 {
		smartTimeout = 30
	}

	for len(responses) < len(preBookRequests) {

		remaining := float64(smartTimeout) - totalTime.Seconds()
		time1 := time.Now()

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*PreBookResponse))
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

		if i < len(preBookRequests) {
			go m.SendPreBookRequestChannel(closableChannel, preBookRequests[i])
			i++
		}
	}

	return responses
}

func (m *RestelClient) SendPreBookRequestChannel(closableChannel *ClosableChannel, preBookRequest *PreBookRequest) {

	//("SendSearchRequestChannel : %+v, %+v\n", closableChannel, searchRequest)

	defer func() {
		if err := recover(); err != nil {
			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	preBookResponse := m.SendPreBookRequest(preBookRequest, true)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- preBookResponse
	})
}

func (m *RestelClient) SendCancelFeeDetailRequests(requests []*CancelFeeDetailRequest) []*CancelFeeDetailResponse {

	//fmt.Printf("SendSearchRequests - Start %d\n", len(searchRequests))
	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := NewClosableChannel()

	var responses []*CancelFeeDetailResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(requests) < maxBatches {
		maxBatches = len(requests)
	}

	for i < maxBatches {
		go m.SendCancelFeeDetailRequestChannel(closableChannel, requests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.SmartTimeout

	for len(responses) < len(requests) {

		remaining := float64(smartTimeout) - totalTime.Seconds()
		time1 := time.Now()

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*CancelFeeDetailResponse))
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

		if i < len(requests) {
			go m.SendCancelFeeDetailRequestChannel(closableChannel, requests[i])
			i++
		}
	}

	return responses
}

func (m *RestelClient) SendCancelFeeDetailRequestChannel(closableChannel *ClosableChannel, request *CancelFeeDetailRequest) {

	//("SendSearchRequestChannel : %+v, %+v\n", closableChannel, searchRequest)

	defer func() {
		if err := recover(); err != nil {

			// m.Log.LogEvent(roomresutils.EventTypeError,
			// 	"RESTEL SendSearchRequest ",
			// 	fmt.Sprintf("%s", errors.Wrap(err, 2).ErrorStack()),
			// )

			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	searchResponse := m.SendCancelFeeDetailRequest(request)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}

func (m *RestelClient) SendCancelFeeDetailRequest(request *CancelFeeDetailRequest) *CancelFeeDetailResponse {
	fmt.Println("SendCancelFeeDetailRequest ----------- Start ----------")

	serializer := NewSerializer(true)
	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	requestData, _ := serializer.Serialize(request)

	//fmt.Printf("%+v\n", searchRequest)
	//fmt.Printf("requestData = %s\n", requestData)

	//requestData = []byte(fmt.Sprintf("XML=%s", url.QueryEscape(string(requestData))))

	bodyParameters := m.CreateHttpBodyParameters()
	bodyParameters["xml"] = string(requestData)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			//RequestBody:          requestData,
			BodyParameters: bodyParameters,
		})

	requestData = nil

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var response *CancelFeeDetailResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&response)

	//errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	response.CancelFeeDetailResponseParam.RequestKey = GetStringArrayHash(request.CancelFeeDetailRequestParam.BookingLine)

	//fmt.Printf("Response = %+v\n", searchResponse)

	return response
}
