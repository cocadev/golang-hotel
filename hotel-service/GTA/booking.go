package gta

import (
	"encoding/xml"
	"fmt"

	roomresutils "../roomres/utils"
)

func (m *GtaClient) SendBookRequest(bookingRequest *BookingRequest) *BookingResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(bookingRequest)
	fmt.Printf("Request = %s\n", requestData)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "SendBookRequest - request data", fmt.Sprintf("%s", requestData))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "SendBookRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - request", errResponseA.Error())
	}

	var bookResponse BookingResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := xml.Unmarshal(
		httpResponse,
		&bookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - deserialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRequest - deserialization", errResponse.Error())
	}

	return &bookResponse
}

func (m *GtaClient) SendCancelRequest(cancelRequest *CancelRequest) *CancelResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(cancelRequest)
	fmt.Printf("Request = %s\n", requestData)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendCancelRequest - serialization", fmt.Sprintf("%v", cancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendCancelRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "SendCancelRequest - request data", fmt.Sprintf("%s", cancelRequest))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "SendCancelRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "SendCancelRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "SendBookRSendCancelRequestequest - request", errResponseA.Error())
	}

	var cancelResponse CancelResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := xml.Unmarshal(
		httpResponse,
		&cancelResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "SendCancelRequest - deserialization", fmt.Sprintf("%v", cancelRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "SendCancelRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "SendCancelRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, "SendCancelRequest - deserialization", errResponse.Error())
	}

	return &cancelResponse
}
