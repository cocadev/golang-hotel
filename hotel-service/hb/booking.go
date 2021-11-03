package hb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	roomresutils "../roomres/utils"
)

func (m *HBClient) SendBookRequest(bookingRequest *BookingRequest) *BookingResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.BookingEndPoint, m.CreateHttpHeaders())

	requestData, err := json.Marshal(bookingRequest)

	fmt.Printf("requestData=%s\n", requestData)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendBookRequest - request data", fmt.Sprintf("%s", requestData))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendBookRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - request", errResponseA.Error())
	}

	var bookResponse BookingResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := json.Unmarshal(
		httpResponse,
		&bookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", errResponse.Error())
	}

	return &bookResponse
}

func (m *HBClient) SendBookDetailRequest(bookingRequest *BookingRequest) *BookingResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.BookingEndPoint, m.CreateHttpHeaders())

	requestData, err := json.Marshal(bookingRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendBookRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendBookRequest - request data", fmt.Sprintf("%s", requestData))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendBookRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - request", errResponseA.Error())
	}

	var bookResponse BookingResponse

	//fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := json.Unmarshal(
		httpResponse,
		&bookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%v", bookingRequest))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendBookRequest - deserialization", errResponse.Error())
	}

	return &bookResponse
}

func (m *HBClient) SendCancelRequest(bookingRef string) *BookingResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	requestUrl := m.BookingEndPoint + "/" + bookingRef
	httpResponse, errResponseA := SendDelete(m.CreateHttpHeaders(), requestUrl)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendCancelRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCancelRequest - request", fmt.Sprintf("%s", requestUrl))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCancelRequest - request", errResponseA.Error())
	}

	var bookResponse BookingResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := json.Unmarshal(
		httpResponse,
		&bookResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCancelRequest - deserialization", fmt.Sprintf("%s", requestUrl))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCancelRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCancelRequest - deserialization", errResponse.Error())
	}

	return &bookResponse
}

func SendDelete(headers map[string]string, endPoint string) ([]byte, error) {

	client := http.Client{}

	req, reqErr := http.NewRequest("DELETE", endPoint, bytes.NewBuffer([]byte{}))

	if reqErr != nil {
		return nil, reqErr
	}

	for headerName, headerValue := range headers {
		req.Header.Set(headerName, headerValue)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	respData, errResp := ioutil.ReadAll(resp.Body)

	return respData, errResp
}
