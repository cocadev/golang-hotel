package gta

import (
	"encoding/xml"
	"fmt"

	roomresutils "../roomres/utils"

	hbecommon "../roomres/hbe/common"
)

func (m *GtaClient) NewSearchItemInformationServiceRequest(request *hbecommon.BookingInfoRequest) *SearchItemInformationServiceRequest {

	serviceRequest := &SearchItemInformationServiceRequest{
		//Source: m.NewSource(request.CurrencyCode),
		Source: m.NewSource(m.ProfileCurrency),
	}

	return serviceRequest
}

func (m *GtaClient) GetBookingInfo(request *hbecommon.BookingInfoRequest) *hbecommon.BookingInfoResponse {

	// repository := m.RepositoryFactory.CreateHotelMappingRepository()

	// hotelMappings := repository.ListHotelMappingByProviderAndPrimaryProviderId(GtaProviderId, 1, []string{strconv.Itoa(request.HotelId)})

	// if len(hotelMappings) == 0 {
	// 	return nil
	// }

	// searchItemInformationRequest := m.NewSearchItemInformationServiceRequest(request)

	// mapping_hbe_to_searchiteminformationrequest(request, hotelMappings, searchItemInformationRequest)

	// searchItemInformationResponse := m.GetSearchInformationItem(searchItemInformationRequest)

	// response := &hbecommon.BookingInfoResponse{}

	// mapping_searchiteminformationrequest_to_hbe(searchItemInformationResponse, response)

	// return response
	return nil
}

func (m *GtaClient) GetSearchInformationItem(request *SearchItemInformationServiceRequest) *SearchItemInformationServiceResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(request)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - serialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA GetSearchInformationItem - request data", fmt.Sprintf("%s", requestData))
	}

	responseData, errResponse := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA GetSearchInformationItem - request data", fmt.Sprintf("%s", responseData))
	}

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - request", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - request", errResponse.Error())
	}

	var response *SearchItemInformationServiceResponse

	errResponse = xml.Unmarshal(responseData, &response)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - deserialization", fmt.Sprintf("%s", responseData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA GetSearchInformationItem - deserialization", errResponse.Error())
	}

	return response
}

func (m *GtaClient) NewSearchBookingRequest() *SearchBookingServiceRequest {

	searchRequest := &SearchBookingServiceRequest{
		//Source: m.NewSource(hotelSearchRequest.CurrencyCode),
		Source:               m.NewSource(m.ProfileCurrency),
		SearchBookingRequest: &SearchBookingRequest{},
	}

	return searchRequest
}

func (m *GtaClient) SearchBookings(request *SearchBookingServiceRequest) *SearchBookingServiceReponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(request)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - serialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SearchBookings - request data", fmt.Sprintf("%s", requestData))
	}

	responseData, errResponse := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SearchBookings - request data", fmt.Sprintf("%s", responseData))
	}

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - request", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - request", errResponse.Error())
	}

	var response *SearchBookingServiceReponse

	errResponse = xml.Unmarshal(responseData, &response)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%s", responseData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", errResponse.Error())
	}

	return response
}

func (m *GtaClient) NewSearchBookinItemRequest() *SearchBookingItemServiceRequest {

	searchRequest := &SearchBookingItemServiceRequest{
		//Source: m.NewSource(hotelSearchRequest.CurrencyCode),
		Source: m.NewSource(m.ProfileCurrency),
	}

	return searchRequest
}

func (m *GtaClient) SearchBookingItem(request *SearchBookingItemServiceRequest) *SearchBookingItemServiceResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(request)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookingItem - serialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookingItem - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SearchBookingItem - request data", fmt.Sprintf("%s", requestData))
	}

	responseData, errResponse := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SearchBookingItem - request data", fmt.Sprintf("%s", responseData))
	}

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookingItem - request", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookingItem - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookingItem - request", errResponse.Error())
	}

	var response *SearchBookingItemServiceResponse

	errResponse = xml.Unmarshal(responseData, &response)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%v", request))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", fmt.Sprintf("%s", responseData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SearchBookings - deserialization", errResponse.Error())
	}

	return response
}

func (m *GtaClient) ListBookingReport(request *hbecommon.BookingReportRequest) *hbecommon.BookingReportResponse {

	requestReport := m.NewSearchBookingRequest()
	mapping_hbe_to_bookingreportrequest(request, requestReport.SearchBookingRequest)

	//fmt.Printf("1. %v\n", requestReport)

	response := m.SearchBookings(requestReport)

	//fmt.Printf("2. %v\n", response)

	responseHbe := &hbecommon.BookingReportResponse{}

	for _, booking := range response.SearchBookings {

		//fmt.Printf("3. %v\n", booking)

		bookingReference := booking.GetReference("client")

		if bookingReference == nil {
			continue
		}

		requestBookingItem := m.NewSearchBookinItemRequest()

		requestBookingItem.SearchBookingItemRequest = &SearchBookingItemRequest{BookingReference: bookingReference}

		responseBookingItem := m.SearchBookingItem(requestBookingItem)

		mapping_bookingreportresponse_to_hbe(booking, responseBookingItem, responseHbe)
	}

	return responseHbe
}

func (m *GtaClient) Init() {

}
