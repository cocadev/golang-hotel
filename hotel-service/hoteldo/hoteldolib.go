package hoteldo

import (
	"encoding/xml"
	"fmt"
	"time"

	roomresutils "../roomres/utils"

	hbecommon "../roomres/hbe/common"
)

type HotelDoClient struct {
	Settings *HotelDoSettings
}

type HotelDoSettings struct {
	HotelsPerBatch  int `json:"HotelsPerBatch"`
	NumberOfThreads int `json:"NumberOfThreads"`
	SmartTimeout    int `json:"SmartTimeout"`

	SearchEndPoint        string `json:"SearchEndPoint"`
	BookingEndPoint       string `json:"BookingEndPoint"`
	BookingCancelEndPoint string `json:"BookingCancelEndPoint"`

	AffiliateID string `json:"AffiliateID"`
	SecretKey   string `json:"SecretKey"`

	ProfileCurrency string `json:"ProfileCurrency"`
	ProfileCountry  string `json:"ProfileCountry"`
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings) *HotelDoClient {

	var hotelDoSettings *HotelDoSettings = &HotelDoSettings{
		SearchEndPoint:        hotelProviderSettings.SearchEndPoint,
		BookingEndPoint:       hotelProviderSettings.BookingConfirmationEndPoint,
		BookingCancelEndPoint: hotelProviderSettings.BookingCancellationEndPoint,

		AffiliateID:     hotelProviderSettings.UserName,
		SecretKey:       hotelProviderSettings.SecretKey,
		HotelsPerBatch:  hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads: hotelProviderSettings.NumberOfThreads,
		SmartTimeout:    hotelProviderSettings.SmartTimeout,
		ProfileCurrency: hotelProviderSettings.ProfileCurrency,
		ProfileCountry:  hotelProviderSettings.ProfileCountry,
	}

	// err := json.Unmarshal([]byte(hotelProviderSettings.Metadata), &hotelDoSettings)

	// if err != nil {
	// 	fmt.Printf("Metadata Parsing Error=%s\n", err)
	// }

	hotelDoClient := &HotelDoClient{
		Settings: hotelDoSettings,
	}

	return hotelDoClient
}

func (m *HotelDoClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func (m *HotelDoClient) CreateHttpBookingHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/xml",
	}
}

func (m *HotelDoClient) CreateHttpSearchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate",
	}
}

func (m *HotelDoClient) NewSearchRequest() *SearchRequest {

	searchRequest := &SearchRequest{}
	searchRequest.SearchDetails = &SearchDetails{}
	searchRequest.SearchDetails.AffiliateID = m.Settings.AffiliateID
	searchRequest.SearchDetails.CountryCode = m.Settings.ProfileCountry
	searchRequest.SearchDetails.CurrencyCode = m.Settings.ProfileCurrency

	return searchRequest
}

/*
	Main Entry Point to Search Request
	Parameter - hbecommon.HotelSearchRequest
	Return Value - hbecommon.HotelSearchRequest
*/
func (m *HotelDoClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {

	hotelIdTransformer := NewHotelIdTransformer()

	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	request := m.NewSearchRequest()
	mapping_hbe_to_searchrequest(hotelSearchRequest, request, hotelIdTransformer)

	requests := SplitBatchSearchRequests(request, m.Settings.HotelsPerBatch)
	response := m.SendSearchRequests(requests)

	mapping_searchresponse_to_hbe(
		response,
		hotelSearchRequest,
		hotelSearchResponse,
		hotelIdTransformer,
	)

	return hotelSearchResponse
}

func (m *HotelDoClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {
	bookRequest := &BookingRequest{
		Type:        "Reservation",
		Version:     "1.0",
		Affiliateid: m.Settings.AffiliateID,
		Language:    "ing",
		Currency:    m.Settings.ProfileCurrency,
		//Uid:           fmt.Sprintf("%s", uuid.NewV4()), //TODO : how to get this uid?
		ClientCountry: m.Settings.ProfileCountry,
	}
	makeBookingRequest(hbeBookingRequest, bookRequest)

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse)

	return hbeBookingResponse
}

func (m *HotelDoClient) SendSearchRequest(searchRequest *SearchRequest) (*SearchResponse, string) {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.SearchEndPoint, m.CreateHttpHeaders())

	searchRequest.SearchDetails.StartDate = convertTime(searchRequest.SearchDetails.StartDate)
	searchRequest.SearchDetails.EndDate = convertTime(searchRequest.SearchDetails.EndDate)
	var urlParameters string

	urlParameters = fmt.Sprintf("?a=%s&co=%s&c=%s&sd=%s&ed=%s&mp=%s&r=%d&d=%d&l=%s&hash=ha:true",
		searchRequest.SearchDetails.AffiliateID, searchRequest.SearchDetails.CountryCode,
		searchRequest.SearchDetails.CurrencyCode, searchRequest.SearchDetails.StartDate,
		searchRequest.SearchDetails.EndDate, searchRequest.SearchDetails.MealPlan,
		searchRequest.SearchDetails.NumberOfRooms, searchRequest.SearchDetails.DestinationCode,
		searchRequest.SearchDetails.LanguageId)
	if searchRequest.SearchDetails.Order != "" {
		urlParameters += fmt.Sprintf("&order=%s", searchRequest.SearchDetails.Order)
	}

	urlParameters += fmt.Sprintf("&h=")
	for _, hotelId := range searchRequest.SearchDetails.HotelIds {
		urlParameters += fmt.Sprintf("%s,", hotelId)
	}
	if len(searchRequest.SearchDetails.HotelIds) > 0 {
		urlParameters = TrimSuffix(urlParameters, ",")
	}

	var i int = 1
	for _, room := range searchRequest.SearchDetails.RequestedRooms {
		urlParameters += fmt.Sprintf("&r%da=%d&r%dk=%d", i, room.Adults, i, room.Children)

		var j int = 1
		for _, childAge := range room.ChildAges {
			urlParameters += fmt.Sprintf("&r%dk%da=%d", i, j, childAge.Age)
			j++
		}
		i++
	}
	fmt.Printf("SendSearchRequest - %s\n", urlParameters)

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:        "GET",
			UrlParameters: urlParameters,
		})

	var searchResponse *SearchResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Println(fmt.Sprintf("%v", searchRequest))
		fmt.Println(errResponse.Error())
	}

	return searchResponse, urlParameters
}

func (m *HotelDoClient) SendBookRequest(bookingRequest *BookingRequest) *BookingResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.BookingEndPoint, m.CreateHttpBookingHeaders())

	requestData, err := xml.Marshal(bookingRequest)

	fmt.Printf("requestData=%s\n", requestData)

	if err != nil {
		fmt.Printf("HB SendBookRequest - serialization : %s\n", fmt.Sprintf("%v", bookingRequest))
		fmt.Printf("HB SendBookRequest - serialization : %s\n", err.Error())
	}

	httpResponse, errResponseA := httpClient.Send(requestData)

	if errResponseA != nil {
		fmt.Printf(" SendBookRequest - request : %s\n", fmt.Sprintf("%s", requestData))
		fmt.Printf(" SendBookRequest - request : %s\n", errResponseA.Error())
	}

	var bookResponse BookingResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := xml.Unmarshal(
		httpResponse,
		&bookResponse)

	if errResponse != nil {
		fmt.Printf(" SendBookRequest - deserialization : %s\n", fmt.Sprintf("%v", bookingRequest))
		fmt.Printf(" SendBookRequest - deserialization : %s\n", fmt.Sprintf("%s", requestData))
		fmt.Printf(" SendBookRequest - deserialization : %s\n", fmt.Sprintf("%s", httpResponse))
		fmt.Printf(" SendBookRequest - deserialization : %s\n", errResponse.Error())
	}

	return &bookResponse
}

func (m *HotelDoClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {
	// Hu B code here
	cancelResponse := m.SendCancelRequest(bookingCancelRequest)

	hbeCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_cancelresponse_to_hbe(cancelResponse, bookingCancelRequest, hbeCancelResponse)

	return hbeCancelResponse
}

func (m *HotelDoClient) Init() {

}

func (m *HotelDoClient) SendCancelRequest(bookingCancelRequest *hbecommon.BookingCancelRequest) *BookingCancelResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.BookingCancelEndPoint, m.CreateHttpHeaders())

	var urlParameters string

	urlParameters = fmt.Sprintf("?a=%s&l=%s&c=%s&bn=%s",
		m.Settings.AffiliateID, "ING",
		m.Settings.ProfileCountry,
		bookingCancelRequest.Ref)

	fmt.Printf("Request URL = %s, param=%s\n", m.Settings.BookingCancelEndPoint, urlParameters)

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:        "GET",
			UrlParameters: urlParameters,
		})

	fmt.Printf("response = %s\n", httpResponse.ResponseBody)
	var cancelResponse *BookingCancelResponse

	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &cancelResponse)

	if errResponse != nil {
		fmt.Println(fmt.Sprintf("%v", cancelResponse))
		fmt.Println(errResponse.Error())
	}

	return cancelResponse
}

//local functions
func (m *HotelDoClient) SendSearchRequests(searchRequests []*SearchRequest) *SearchResponse {

	//fmt.Printf("SendSearchRequests - Start %d\n", len(searchRequests))
	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := NewClosableChannel()

	var responses []*SearchResponse

	var i int = 0
	var maxBatches int = m.Settings.NumberOfThreads

	if len(searchRequests) < maxBatches {
		maxBatches = len(searchRequests)
	}

	for i < maxBatches {
		go m.SendSearchRequestChannel(closableChannel, searchRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.Settings.SmartTimeout

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
		Hotels: []*Hotel{},
	}

	for _, response := range responses {

		if response != nil && response.Hotels != nil {
			combinedSearchResponse.Hotels =
				append(
					combinedSearchResponse.Hotels,
					response.Hotels...,
				)
		}
	}

	return combinedSearchResponse
}

func (m *HotelDoClient) SendSearchRequestChannel(closableChannel *ClosableChannel, searchRequest *SearchRequest) {

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

	searchResponse, _ := m.SendSearchRequest(searchRequest)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}
