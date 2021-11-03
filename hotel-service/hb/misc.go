package hb

import (
	"encoding/json"
	"fmt"
	"strings"

	roomresutils "../roomres/utils"

	hbecommon "../roomres/hbe/common"
)

// func (m *HBClient) GetBookingInfo(request *hbecommon.BookingInfoRequest) *hbecommon.BookingInfoResponse {

// 	repository := m.RepositoryFactory.CreateHotelMappingRepository()

// 	hotelMappings := repository.ListHotelMappingByProviderAndPrimaryProviderId(HBProviderId, 1, []string{strconv.Itoa(request.HotelId)})

// 	if len(hotelMappings) == 0 {
// 		return nil
// 	}

// 	hotelContentRequest := &HotelContentRequest{}

// 	mapping_hbe_to_hotelcontentrequest(request, hotelMappings, hotelContentRequest)

// 	hotelContentResponse := m.GetHotelContent(hotelContentRequest)

// 	response := &hbecommon.BookingInfoResponse{}

// 	mapping_hotelcontentresponse_to_hbe(hotelContentResponse, response)

// 	return response
// }

func (m *HBClient) GetHotelContent(request *HotelContentRequest) *HotelContentResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.ContentEndPoint, m.CreateHttpHeaders())

	requestedUrlParams := request.GenerateUrlParams()

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendSearchRequest - request data", fmt.Sprintf("%s - %s", m.ContentEndPoint, requestedUrlParams))
	}

	httpResponse := httpClient.SendRequest(&roomresutils.IntegrationHttpRequest{Method: "GET", UrlParameters: requestedUrlParams, RequestBodySpecified: false})

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendSearchRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	var serviceResponse HotelContentResponse

	errResponse := json.Unmarshal(
		httpResponse.ResponseBody,
		&serviceResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", requestedUrlParams))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", errResponse.Error())
	}

	return &serviceResponse
}

func (m *HBClient) ListBookingReport(request *hbecommon.BookingReportRequest) *hbecommon.BookingReportResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	endPoint := strings.Replace(m.SearchEndPoint, "/hotels", "", 1)

	httpClient := roomresutils.NewIntegrationHttp(endPoint, m.CreateHttpHeaders())

	reportRequest := &BookingReportRequest{}

	mapping_hbe_to_bookingreportrequest(request, reportRequest)

	requestedUrlParams := reportRequest.GenerateUrlParams()

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" BookingReportRequest - request data", fmt.Sprintf("%s - %s", endPoint, requestedUrlParams))
	}

	httpResponse := httpClient.SendRequest(&roomresutils.IntegrationHttpRequest{Method: "GET", UrlParameters: requestedUrlParams, RequestBodySpecified: false})

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" BookingReportRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	var reportResponse BookingReportResponse

	errResponse := json.Unmarshal(
		httpResponse.ResponseBody,
		&reportResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" BookingReportRequest - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" BookingReportRequest - deserialization", fmt.Sprintf("%s", requestedUrlParams))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" BookingReportRequest - deserialization", fmt.Sprintf("%s", httpResponse.ResponseBody))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" BookingReportRequest - deserialization", errResponse.Error())
	}

	response := &hbecommon.BookingReportResponse{}

	mapping_bookingreportresponse_to_hbe(&reportResponse, response, m.Destinations)

	return response
}

func (m *HBClient) LoadDestinations(request *DestinationRequest) *DestinationResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.ContentEndPoint, m.CreateHttpHeaders())

	requestedUrlParams := request.GenerateUrlParams()

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" LoadDestinations - request data", fmt.Sprintf("%s - %s", m.ContentEndPoint, requestedUrlParams))
	}

	httpResponse := httpClient.SendRequest(&roomresutils.IntegrationHttpRequest{Method: "GET", UrlParameters: requestedUrlParams, RequestBodySpecified: false})

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" LoadDestinations - response data", fmt.Sprintf("%s", httpResponse))
	}

	var serviceResponse DestinationResponse

	errResponse := json.Unmarshal(
		httpResponse.ResponseBody,
		&serviceResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" LoadDestinations - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" LoadDestinations - deserialization", fmt.Sprintf("%s", requestedUrlParams))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" LoadDestinations - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" LoadDestinations - deserialization", errResponse.Error())
	}

	return &serviceResponse
}

func (m *HBClient) LoadAllDestinations() *DestinationResponse {

	response := &DestinationResponse{Destinations: []*Destination{}}
	request := &DestinationRequest{IndexFrom: 1}

	for {

		request.IndexTo = request.IndexFrom + 999

		localResponse := m.LoadDestinations(request)

		if len(localResponse.Destinations) > 0 {
			response.Destinations = append(response.Destinations, localResponse.Destinations...)

			request.IndexFrom += len(response.Destinations)
		} else {
			break
		}
	}

	return response
}

func (m *HBClient) Init() {

	response := m.LoadAllDestinations()

	m.Destinations = map[string]*Destination{}

	for _, destination := range response.Destinations {
		m.Destinations[strings.ToUpper(destination.Code)] = destination
	}

	//fmt.Printf("%v", m.Destinations)
}
