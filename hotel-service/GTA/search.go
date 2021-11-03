package gta

import (
	"encoding/xml"
	"fmt"
	"time"

	roomresutils "../roomres/utils"

	"github.com/go-errors/errors"
)

func (m *GtaClient) SendSearchRequests(searchRequests []*SearchRequest) *SearchResponse {

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

	for len(responses) < len(searchRequests) {

		remaining := float64(m.SmartTimeout) - totalTime.Seconds()
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

		if totalTime.Seconds() >= float64(m.SmartTimeout) {
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
			"GTA SendSearchRequest ",
			fmt.Sprintf("%d - %d", len(searchRequests), len(responses)),
		)
	}

	combinedSearchResponse := &SearchResponse{}

	for _, response := range responses {

		if response != nil {
			combinedSearchResponse.Hotels = append(combinedSearchResponse.Hotels, response.Hotels...)
		}
	}

	return combinedSearchResponse
}

func (m *GtaClient) SendSearchRequestChannel(closableChannel *roomresutils.ClosableChannel, searchRequest *SearchRequest) {

	defer func() {
		if err := recover(); err != nil {

			m.Log.LogEvent(roomresutils.EventTypeError,
				"GTA SendSearchRequest ",
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

func (m *GtaClient) SendSearchRequest(searchRequest *SearchRequest) *SearchResponse {

	logScope := m.Log.StartLogScope(roomresutils.LogScopeRef{})

	serializer := roomresutils.NewSerializer(true)
	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, err := serializer.Serialize(searchRequest)

	fmt.Printf("Request = %s\n", requestData)

	if err != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - serialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - serialization", err.Error())
	}

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SendSearchRequest - request data", fmt.Sprintf("%s", requestData))
	}

	responseData, errResponse := httpClient.Send(requestData)

	if logScope.AllowLogging(roomresutils.EventTypeInfo3) {

		logScope.LogEvent(roomresutils.EventTypeInfo3, "GTA SendSearchRequest - request data", fmt.Sprintf("%s", responseData))
	}

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - request", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - request", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - request", errResponse.Error())
	}

	fmt.Printf("responseData = %s\n", responseData)

	var searchResponse *SearchResponse

	errResponse = xml.Unmarshal(responseData, &searchResponse)

	if errResponse != nil {
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - deserialization", fmt.Sprintf("%v", searchRequest))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - deserialization", fmt.Sprintf("%s", requestData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - deserialization", fmt.Sprintf("%s", responseData))
		logScope.LogEvent(roomresutils.EventTypeError, "GTA SendSearchRequest - deserialization", errResponse.Error())
	}

	return searchResponse
}
