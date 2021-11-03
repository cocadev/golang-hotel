package hb

import (
	"encoding/json"
	"fmt"
	"time"

	roomresutils "../roomres/utils"

	"github.com/go-errors/errors"
)

func (m *HBClient) SplitSearchRequests(combinedSearchRequest *AvailabilityRequest) (searchRequests []*AvailabilityRequest) {

	if len(combinedSearchRequest.Hotels.HotelIds) > m.HotelsPerBatch {

		searchRequest := combinedSearchRequest.Clone()
		searchRequest.Hotels.HotelIds = []int{}

		for _, hotelId := range combinedSearchRequest.Hotels.HotelIds {

			if len(searchRequest.Hotels.HotelIds) >= m.HotelsPerBatch {
				searchRequests = append(searchRequests, searchRequest)

				searchRequest = combinedSearchRequest.Clone()
				searchRequest.Hotels.HotelIds = []int{}
			}

			searchRequest.Hotels.HotelIds = append(searchRequest.Hotels.HotelIds, hotelId)
		}

		searchRequests = append(searchRequests, searchRequest)

	} else {
		searchRequests = append(searchRequests, combinedSearchRequest)
	}

	return searchRequests
}

func (m *HBClient) SendSearchRequests(searchRequests []*AvailabilityRequest) *AvailabilityResponse {

	//var requestChannel chan *AvailabilityResponse = make(chan *AvailabilityResponse)
	closableChannel := roomresutils.NewClosableChannel()

	var responses []*AvailabilityResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(searchRequests) < maxBatches {
		maxBatches = len(searchRequests)
	}

	for i < maxBatches {
		go m.SendSearchRequestChannel(closableChannel, searchRequests[i])
		i++
	}

	for len(responses) < len(searchRequests) {

		exit := false

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*AvailabilityResponse))
		case <-time.After(time.Second * 20):
			exit = true
			closableChannel.Close()
			break
		}

		if exit {
			break
		}

		if i < len(searchRequests) {
			go m.SendSearchRequestChannel(closableChannel, searchRequests[i])
			i++
		}
	}

	combinedSearchResponse := &AvailabilityResponse{Hotels: &ResponseHotels{}}

	for _, response := range responses {

		if response != nil && response.Hotels != nil {
			combinedSearchResponse.Hotels.Hotels =
				append(
					combinedSearchResponse.Hotels.Hotels,
					response.Hotels.Hotels...,
				)
		}
	}

	return combinedSearchResponse
}

func (m *HBClient) SendSearchRequestChannel(closableChannel *roomresutils.ClosableChannel, searchRequest *AvailabilityRequest) {

	defer func() {
		if err := recover(); err != nil {

			m.Log.LogEvent(roomresutils.EventTypeError,
				m.GetProviderName()+" SendSearchRequest ",
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

func (m *HBClient) SendSearchRequest(searchRequest *AvailabilityRequest) *AvailabilityResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := json.Marshal(searchRequest)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendSearchRequest - serialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendSearchRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendSearchRequest - request data", fmt.Sprintf("%s", requestData))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendSearchRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - request", errResponseA.Error())
	}

	var serviceResponse AvailabilityResponse

	//fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := json.Unmarshal(
		httpResponse,
		&serviceResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", errResponse.Error())
	}

	return &serviceResponse
}

//Added by Li, 20180917
func (m *HBClient) SendCheckRateRequest(searchRequest *CheckRateRequest) *AvailabilityResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	httpClient := roomresutils.NewIntegrationHttp(m.CheckRateEndPoint, m.CreateHttpHeaders())

	requestData, err := json.Marshal(searchRequest)
	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendCheckRateRequest - serialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "HB SendCheckRateRequest - serialization", err.Error())
	}
	fmt.Printf("httpRequest = %s\n\n", requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendCheckRateRequest - request data", fmt.Sprintf("%s", requestData))
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, m.GetProviderName()+" SendCheckRateRequest - response data", fmt.Sprintf("%s", httpResponse))
	}

	if errResponseA != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCheckRateRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendCheckRateRequest - request", errResponseA.Error())
	}

	var originResponse CheckRateResponse

	fmt.Printf("httpResponse = %s\n\n", httpResponse)

	errResponse := json.Unmarshal(
		httpResponse,
		&originResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", fmt.Sprintf("%s", httpResponse))
		logScope.LogEvent(roomresutils.EventTypeError, m.GetProviderName()+" SendSearchRequest - deserialization", errResponse.Error())
	}
	hotels := []*Hotel{}
	if originResponse.Hotel != nil {
		hotels = append(hotels, originResponse.Hotel)
	}
	serviceResponse := AvailabilityResponse{
		Hotels: &ResponseHotels{
			Hotels: hotels,
		},
	}

	return &serviceResponse
}
