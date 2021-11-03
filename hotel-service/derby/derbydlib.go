package derby

import (
	"encoding/json"
	"fmt"
	"time"

	hbecommon "../roomres/hbe/common"
	roomresutils "../roomres/utils"
)

type DerbyClient struct {
	Settings        *DerbySettings
	HotelsPerBatch  int `json:"HotelsPerBatch"`
	NumberOfThreads int `json:"NumberOfThreads"`
	SmartTimeout    int `json:"SmartTimeout"`

	SearchEndPoint         string `json:"SearchEndPoint"`
	PreBookingEndPoint     string `json:"PreBookingEndPoint"`
	BookingEndPoint        string `json:"BookingEndPoint"`
	BookingCancelEndPoint  string `json:"BookingCancelEndPoint"`
	BookingDetailsEndpoint string `json:"BookingDetailsEndpoint"`
	CreditCardInfo         *hbecommon.CreditCardInfo
}

type DerbySettings struct {
	DistributorId string `json:"DistributorId"`
	SupplierId    string `json:"SupplierId"`
	Version       string `json:"Version"`
	Token         string `json:"Token"`
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings) *DerbyClient {

	var derbySettings *DerbySettings = &DerbySettings{}

	err := json.Unmarshal([]byte(hotelProviderSettings.Metadata), &derbySettings)
	if err != nil {
		fmt.Printf("Error : derbySettings : %s", hotelProviderSettings.Metadata)
	}

	derbyClient := &DerbyClient{
		Settings:               derbySettings,
		SearchEndPoint:         hotelProviderSettings.SearchEndPoint,
		PreBookingEndPoint:     hotelProviderSettings.PreBookEndPoint,
		BookingEndPoint:        hotelProviderSettings.BookingConfirmationEndPoint,
		BookingCancelEndPoint:  hotelProviderSettings.BookingCancellationEndPoint,
		BookingDetailsEndpoint: hotelProviderSettings.BookingDetailsEndpoint,

		HotelsPerBatch:  hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads: hotelProviderSettings.NumberOfThreads,
		SmartTimeout:    hotelProviderSettings.SmartTimeout,

		CreditCardInfo: hotelProviderSettings.CreditCardInfo,
	}

	return derbyClient
}

func (m *DerbyClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + m.Settings.Token,
	}
}

func (m *DerbyClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {
	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	requests := mapping_hbe_to_searchrequest(hotelSearchRequest)
	for _, request := range requests {
		request.Header = &Header{
			DistributorId: m.Settings.DistributorId,
			SupplierId:    m.Settings.SupplierId,
			Version:       m.Settings.Version,
			Token:         m.Settings.Token,
		}

		fmt.Printf("request = %+v\n", request)
	}

	responses := m.SendSearchRequests(requests)
	mapping_searchresponse_to_hbe(
		responses,
		hotelSearchRequest,
		hotelSearchResponse,
	)

	return hotelSearchResponse
}

func (m *DerbyClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {
	//prebooking first
	preBookRequest := &PreBookingRequest{}
	preBookRequest.Header = &Header{
		DistributorId: m.Settings.DistributorId,
		SupplierId:    m.Settings.SupplierId,
		Version:       m.Settings.Version,
		Token:         m.Settings.Token,
	}

	makePreBookingRequest(hbeBookingRequest, preBookRequest)
	preBookRequest.Payment = &Payment{
		CardCode:       m.CreditCardInfo.CardType,
		CardNumber:     m.CreditCardInfo.Number,
		CardHolderName: m.CreditCardInfo.HolderName,
		ExpireDate:     m.CreditCardInfo.ExpiryDate,
	}

	preBookResponse := m.SendPreBookRequest(preBookRequest)

	bookRequest := makeBookingRequest(preBookRequest, preBookResponse)

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse)

	return hbeBookingResponse
}

func (m *DerbyClient) SendSearchRequest(searchRequest *SearchRequest) (*SearchResponse, error) {
	fmt.Printf("Called SendSearchRequest = %+v\n", searchRequest)

	httpClient := roomresutils.NewIntegrationHttp(m.SearchEndPoint, m.CreateHttpHeaders())

	requestData, _ := json.Marshal(searchRequest)
	fmt.Printf("requestData = %s\n", requestData)
	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil, nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	searchResponse := &SearchResponse{}

	errResponse := json.Unmarshal(httpResponse.ResponseBody, searchResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v\n", errResponse)
		return nil, nil
	}

	return searchResponse, nil
}

func (m *DerbyClient) SendPreBookRequest(preBookingRequest *PreBookingRequest) *PreBookResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.PreBookingEndPoint, m.CreateHttpHeaders())

	requestData, _ := json.Marshal(preBookingRequest)
	//fmt.Printf("requestData = %s\n", requestData)
	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	response := &PreBookResponse{}

	errResponse := json.Unmarshal(httpResponse.ResponseBody, response)

	if errResponse != nil {
		fmt.Printf("error = %+v\n", errResponse)
		return nil
	}

	return response
}

func (m *DerbyClient) SendBookRequest(bookingRequest *BookingRequest) *BookResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.BookingEndPoint, m.CreateHttpHeaders())

	requestData, _ := json.Marshal(bookingRequest)
	//fmt.Printf("requestData = %s\n", requestData)
	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	response := &BookResponse{}

	errResponse := json.Unmarshal(httpResponse.ResponseBody, response)

	if errResponse != nil {
		fmt.Printf("error = %+v\n", errResponse)
		return nil
	}

	return response
}

func (m *DerbyClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {
	cancelRequest := &CancelBookRequest{}
	cancelRequest.Header = &Header{
		DistributorId: m.Settings.DistributorId,
		SupplierId:    m.Settings.SupplierId,
		Version:       m.Settings.Version,
		Token:         m.Settings.Token,
	}

	cancelRequest.ReservationIds = &ReservationIds{
		DistributorResId: bookingCancelRequest.Ref,
		DerbyResId:       bookingCancelRequest.InternalRef,
	}
	cancelResponse := m.SendCancelRequest(cancelRequest)

	hbeCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_cancelresponse_to_hbe(cancelResponse, bookingCancelRequest, hbeCancelResponse)

	return hbeCancelResponse
}

func (m *DerbyClient) Init() {

}

func (m *DerbyClient) SendCancelRequest(bookingCancelRequest *CancelBookRequest) *BookingCancelResponse {
	httpClient := roomresutils.NewIntegrationHttp(m.BookingCancelEndPoint, m.CreateHttpHeaders())

	requestData, _ := json.Marshal(bookingCancelRequest)
	//fmt.Printf("requestData = %s\n", requestData)
	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          requestData,
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	response := &BookingCancelResponse{}

	errResponse := json.Unmarshal(httpResponse.ResponseBody, response)

	if errResponse != nil {
		fmt.Printf("error = %+v\n", errResponse)
		return nil
	}

	return response
}

//local functions
func (m *DerbyClient) SendSearchRequests(searchRequests []*SearchRequest) []*SearchResponse {
	fmt.Printf("SendSearchRequests = %d\n", len(searchRequests))
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

	return responses
}

func (m *DerbyClient) SendSearchRequestChannel(closableChannel *ClosableChannel, searchRequest *SearchRequest) {

	defer func() {
		if err := recover(); err != nil {
			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	searchResponse, _ := m.SendSearchRequest(searchRequest)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}

func (m *DerbyClient) SendBookRequests(bookingRequests []*BookingRequest) []*BookResponse {

	closableChannel := NewClosableChannel()

	var responses []*BookResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(bookingRequests) < maxBatches {
		maxBatches = len(bookingRequests)
	}

	for i < maxBatches {
		go m.SendBookingRequestChannel(closableChannel, bookingRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.SmartTimeout

	for len(responses) < len(bookingRequests) {

		remaining := float64(smartTimeout) - totalTime.Seconds()
		time1 := time.Now()

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*BookResponse))
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

		if i < len(bookingRequests) {
			go m.SendBookingRequestChannel(closableChannel, bookingRequests[i])
			i++
		}
	}

	return responses
}

func (m *DerbyClient) SendBookingRequestChannel(closableChannel *ClosableChannel, bookingRequest *BookingRequest) {

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

	response := m.SendBookRequest(bookingRequest)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- response
	})
}
