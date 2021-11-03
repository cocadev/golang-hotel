package webbed

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	roomresutils "../roomres/utils"

	hbecommon "../roomres/hbe/common"
)

type WebBedClient struct {
	Settings *WebBedSettings
}

type WebBedSettings struct {
	HotelsPerBatch  int `json:"HotelsPerBatch"`
	NumberOfThreads int `json:"NumberOfThreads"`
	SmartTimeout    int `json:"SmartTimeout"`

	SearchEndPoint        string `json:"SearchEndPoint"`
	PreBookingEndPoint    string `json:"PreBookingEndPoint"`
	BookingEndPoint       string `json:"BookingEndPoint"`
	BookingCancelEndPoint string `json:"BookingCancelEndPoint"`

	UserName string `json:"UserName"`
	Password string `json:"Password"`

	ProfileCurrency string `json:"ProfileCurrency"`
	ProfileCountry  string `json:"ProfileCountry"`
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings) *WebBedClient {

	var webbedSettings *WebBedSettings = &WebBedSettings{
		SearchEndPoint:        hotelProviderSettings.SearchEndPoint,
		PreBookingEndPoint:    hotelProviderSettings.PreBookEndPoint,
		BookingEndPoint:       hotelProviderSettings.BookingConfirmationEndPoint,
		BookingCancelEndPoint: hotelProviderSettings.BookingCancellationEndPoint,

		UserName:        hotelProviderSettings.UserName,
		Password:        hotelProviderSettings.SecretKey,
		HotelsPerBatch:  hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads: hotelProviderSettings.NumberOfThreads,
		SmartTimeout:    hotelProviderSettings.SmartTimeout,
		ProfileCurrency: hotelProviderSettings.ProfileCurrency,
		ProfileCountry:  hotelProviderSettings.ProfileCountry,
	}

	webbedClient := &WebBedClient{
		Settings: webbedSettings,
	}

	return webbedClient
}

func (m *WebBedClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

func (m *WebBedClient) CreateHttpBookingHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/xml",
	}
}

func (m *WebBedClient) CreateHttpSearchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate",
	}
}

func (m *WebBedClient) NewSearchRequest() *SearchRequest {

	searchRequest := &SearchRequest{
		UserName:        m.Settings.UserName,
		Password:        m.Settings.Password,
		Currencies:      m.Settings.ProfileCurrency,
		CustomerCountry: m.Settings.ProfileCountry,
		Language:        "en",
	}

	return searchRequest
}

func (m *WebBedClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {

	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	request := m.NewSearchRequest()
	mapping_hbe_to_searchrequest(hotelSearchRequest, request)

	requests := SplitBatchSearchRequests(request, m.Settings.HotelsPerBatch)
	response := m.SendSearchRequests(requests)

	mapping_searchresponse_to_hbe(
		response,
		hotelSearchRequest,
		hotelSearchResponse,
	)

	return hotelSearchResponse
}

func (m *WebBedClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {
	//prebooking first
	preBookRequest := &PreBookingRequest{}
	preBookRequest.UserName = m.Settings.UserName
	preBookRequest.Password = m.Settings.Password
	preBookRequest.Currency = m.Settings.ProfileCurrency
	preBookRequest.CustomerCountry = m.Settings.ProfileCountry
	preBookRequest.Language = "en"

	makePreBookingRequest(hbeBookingRequest, preBookRequest)

	preBookResponse := m.SendPreBookRequest(preBookRequest)

	fmt.Printf("preBookResponse = %+v\n", preBookResponse)
	bookRequest := makeBookingRequest(hbeBookingRequest, preBookResponse)

	if bookRequest == nil {
		return nil
	}

	bookRequest.UserName = m.Settings.UserName
	bookRequest.Password = m.Settings.Password
	bookRequest.Currency = m.Settings.ProfileCurrency
	bookRequest.CustomerCountry = m.Settings.ProfileCountry
	bookRequest.Language = "en"

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse)

	return hbeBookingResponse
	//return nil
}

func (m *WebBedClient) SendSearchRequest(searchRequest *SearchRequest) (*SearchResponse, string) {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.SearchEndPoint, m.CreateHttpHeaders())

	var urlParameters string

	urlParameters = fmt.Sprintf("?userName=%s&password=%s&language=%s&currencies=%s&checkInDate=%s&checkOutDate=%s&numberOfRooms=%s&destination=%s&destinationID=%s&hotelIDs=%s&resortIDs=%s&accommodationTypes=%s&numberOfAdults=%s&numberOfChildren=%s&childrenAges=%s&infant=%s&sortBy=%s&sortOrder=%s&exactDestinationMatch=%s&blockSuperdeal=%s&showTransfer=%s&mealIds=%s&showCoordinates=%s&showReviews=%s&referencePointLatitude=%s&referencePointLongitude=%s&maxDistanceFromReferencePoint=%s&minStarRating=%s&maxStarRating=%s&featureIds=%s&minPrice=%s&maxPrice=%s&themeIds=%s&excludeSharedRooms=%s&excludeSharedFacilities=%s&prioritizedHotelIds=%s&totalRoomsInBatch=%s&paymentMethodId=%s&CustomerCountry=%s&b2c=%s",
		searchRequest.UserName,
		searchRequest.Password,
		searchRequest.Language,
		searchRequest.Currencies,
		searchRequest.CheckInDate,
		searchRequest.CheckOutDate,
		searchRequest.NumberOfRooms,
		searchRequest.Destination,
		searchRequest.DestinationID,
		strings.Join(searchRequest.HotelIDs, ","),
		searchRequest.ResortIDs,
		searchRequest.AccommodationTypes,
		searchRequest.NumberOfAdults,
		searchRequest.NumberOfChildren,
		searchRequest.ChildrenAges,
		searchRequest.Infant,
		searchRequest.SortBy,
		searchRequest.SortOrder,
		searchRequest.ExactDestinationMatch,
		searchRequest.BlockSuperdeal,
		searchRequest.ShowTransfer,
		searchRequest.MealIds,
		searchRequest.ShowCoordinates,
		searchRequest.ShowReviews,
		searchRequest.ReferencePointLatitude,
		searchRequest.ReferencePointLongitude,
		searchRequest.MaxDistanceFromReferencePoint,
		searchRequest.MinStarRating,
		searchRequest.MaxStarRating,
		searchRequest.FeatureIds,
		searchRequest.MinPrice,
		searchRequest.MaxPrice,
		searchRequest.ThemeIds,
		searchRequest.ExcludeSharedRooms,
		searchRequest.ExcludeSharedFacilities,
		searchRequest.PrioritizedHotelIds,
		searchRequest.TotalRoomsInBatch,
		searchRequest.PaymentMethodId,
		searchRequest.CustomerCountry,
		searchRequest.B2c,
	)

	//fmt.Printf("SendSearchRequest - %s\n", urlParameters)

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:        "GET",
			UrlParameters: urlParameters,
		})

	var searchResponse *SearchResponse

	//fmt.Printf("httpResponse.ResponseBody = %s\n", httpResponse.ResponseBody)
	errResponse := xml.Unmarshal(httpResponse.ResponseBody, &searchResponse)

	if errResponse != nil {
		fmt.Println(fmt.Sprintf("%v", searchRequest))
		fmt.Println(errResponse.Error())
	}

	return searchResponse, urlParameters
}

func (m *WebBedClient) SendPreBookRequest(preBookingRequest *PreBookingRequest) *PreBookResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.PreBookingEndPoint, m.CreateHttpHeaders())

	var urlParameters string

	urlParameters = fmt.Sprintf("?userName=%s&password=%s&currency=%s&language=%s&checkInDate=%s&checkOutDate=%s&roomId=%s&rooms=%d&adults=%d&children=%d&childrenAges=%s&infant=%d&mealId=%s&CustomerCountry=%s&b2c=%d&searchPrice=%s",
		preBookingRequest.UserName,
		preBookingRequest.Password,
		preBookingRequest.Currency,
		preBookingRequest.Language,
		preBookingRequest.CheckInDate,
		preBookingRequest.CheckOutDate,
		preBookingRequest.RoomId,
		preBookingRequest.Rooms,
		preBookingRequest.Adults,
		preBookingRequest.Children,
		preBookingRequest.ChildrenAges,
		preBookingRequest.Infant,
		preBookingRequest.MealId,
		preBookingRequest.CustomerCountry,
		preBookingRequest.B2c,
		preBookingRequest.SearchPrice,
	)

	fmt.Printf("SendPreBookingRequest - %s\n", urlParameters)

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:        "GET",
			UrlParameters: urlParameters,
		})

	var bookResponse PreBookResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := xml.Unmarshal(
		httpResponse.ResponseBody,
		&bookResponse)

	if errResponse != nil {
		fmt.Printf(" SendPreBookRequest - deserialization : %s\n", fmt.Sprintf("%v", preBookingRequest))
		fmt.Printf(" SendPreBookRequest - deserialization : %s\n", fmt.Sprintf("%s", httpResponse))
		fmt.Printf(" SendPreBookRequest - deserialization : %s\n", errResponse.Error())
	}

	bookResponse.RoomId = preBookingRequest.RoomId
	bookResponse.Rooms = preBookingRequest.Rooms
	bookResponse.Adults = preBookingRequest.Adults
	bookResponse.Children = preBookingRequest.Children
	bookResponse.ChildrenAges = preBookingRequest.ChildrenAges
	bookResponse.Infant = preBookingRequest.Infant
	bookResponse.MealId = preBookingRequest.MealId
	bookResponse.Guests = preBookingRequest.Guests

	return &bookResponse
}

func (m *WebBedClient) SendBookRequest(bookingRequest *BookingRequest) *BookResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.BookingEndPoint, m.CreateHttpHeaders())

	var urlParameters string

	urlParameters = fmt.Sprintf("?userName=%s&password=%s&currency=%s&language=%s&email=%s&checkInDate=%s&checkOutDate=%s&roomId=%s&rooms=%d&adults=%d&children=%d&infant=%d&yourRef=%s&specialrequest=%s&mealId=%s&adultGuest1FirstName=%s&adultGuest1LastName=%s&adultGuest2FirstName=%s&adultGuest2LastName=%s&adultGuest3FirstName=%s&adultGuest3LastName=%s&adultGuest4FirstName=%s&adultGuest4LastName=%s&adultGuest5FirstName=%s&adultGuest5LastName=%s&adultGuest6FirstName=%s&adultGuest6LastName=%s&adultGuest7FirstName=%s&adultGuest7LastName=%s&adultGuest8FirstName=%s&adultGuest8LastName=%s&adultGuest9FirstName=%s&adultGuest9LastName=%s&childrenGuest1FirstName=%s&childrenGuest1LastName=%s&childrenGuestAge1=%s&childrenGuest2FirstName=%s&childrenGuest2LastName=%s&childrenGuestAge2=%s&childrenGuest3FirstName=%s&childrenGuest3LastName=%s&childrenGuestAge3=%s&childrenGuest4FirstName=%s&childrenGuest4LastName=%s&childrenGuestAge4=%s&childrenGuest5FirstName=%s&childrenGuest5LastName=%s&childrenGuestAge5=%s&childrenGuest6FirstName=%s&childrenGuest6LastName=%s&childrenGuestAge6=%s&childrenGuest7FirstName=%s&childrenGuest7LastName=%s&childrenGuestAge7=%s&childrenGuest8FirstName=%s&childrenGuest8LastName=%s&childrenGuestAge8=%s&childrenGuest9FirstName=%s&childrenGuest9LastName=%s&childrenGuestAge9=%s&paymentMethodId=%s&creditCardType=%s&creditCardNumber=%s&creditCardHolder=%s&creditCardCVV2=%s&creditCardExpYear=%s&creditCardExpMonth=%s&customerEmail=%s&invoiceRef=%s&commissionAmountInHotelCurrency=%s&CustomerCountry=%s&b2c=%d&preBookCode=%s",
		bookingRequest.UserName,
		bookingRequest.Password,
		bookingRequest.Currency,
		bookingRequest.Language,
		bookingRequest.Email,
		bookingRequest.CheckInDate,
		bookingRequest.CheckOutDate,
		bookingRequest.RoomId,
		bookingRequest.Rooms,
		bookingRequest.Adults,
		bookingRequest.Children,
		bookingRequest.Infant,
		bookingRequest.YourRef,
		bookingRequest.Specialrequest,
		bookingRequest.MealId,
		bookingRequest.AdultGuest1FirstName,
		bookingRequest.AdultGuest1LastName,
		bookingRequest.AdultGuest2FirstName,
		bookingRequest.AdultGuest2LastName,
		bookingRequest.AdultGuest3FirstName,
		bookingRequest.AdultGuest3LastName,
		bookingRequest.AdultGuest4FirstName,
		bookingRequest.AdultGuest4LastName,
		bookingRequest.AdultGuest5FirstName,
		bookingRequest.AdultGuest5LastName,
		bookingRequest.AdultGuest6FirstName,
		bookingRequest.AdultGuest6LastName,
		bookingRequest.AdultGuest7FirstName,
		bookingRequest.AdultGuest7LastName,
		bookingRequest.AdultGuest8FirstName,
		bookingRequest.AdultGuest8LastName,
		bookingRequest.AdultGuest9FirstName,
		bookingRequest.AdultGuest9LastName,
		bookingRequest.ChildrenGuest1FirstName,
		bookingRequest.ChildrenGuest1LastName,
		bookingRequest.ChildrenGuestAge1,
		bookingRequest.ChildrenGuest2FirstName,
		bookingRequest.ChildrenGuest2LastName,
		bookingRequest.ChildrenGuestAge2,
		bookingRequest.ChildrenGuest3FirstName,
		bookingRequest.ChildrenGuest3LastName,
		bookingRequest.ChildrenGuestAge3,
		bookingRequest.ChildrenGuest4FirstName,
		bookingRequest.ChildrenGuest4LastName,
		bookingRequest.ChildrenGuestAge4,
		bookingRequest.ChildrenGuest5FirstName,
		bookingRequest.ChildrenGuest5LastName,
		bookingRequest.ChildrenGuestAge5,
		bookingRequest.ChildrenGuest6FirstName,
		bookingRequest.ChildrenGuest6LastName,
		bookingRequest.ChildrenGuestAge6,
		bookingRequest.ChildrenGuest7FirstName,
		bookingRequest.ChildrenGuest7LastName,
		bookingRequest.ChildrenGuestAge7,
		bookingRequest.ChildrenGuest8FirstName,
		bookingRequest.ChildrenGuest8LastName,
		bookingRequest.ChildrenGuestAge8,
		bookingRequest.ChildrenGuest9FirstName,
		bookingRequest.ChildrenGuest9LastName,
		bookingRequest.ChildrenGuestAge9,
		bookingRequest.PaymentMethodId,
		bookingRequest.CreditCardType,
		bookingRequest.CreditCardNumber,
		bookingRequest.CreditCardHolder,
		bookingRequest.CreditCardCVV2,
		bookingRequest.CreditCardExpYear,
		bookingRequest.CreditCardExpMonth,
		bookingRequest.CustomerEmail,
		bookingRequest.InvoiceRef,
		bookingRequest.CommissionAmountInHotelCurrency,
		bookingRequest.CustomerCountry,
		bookingRequest.B2c,
		bookingRequest.PreBookCode,
	)

	fmt.Printf("SendBookingRequest - %s\n", urlParameters)

	httpResponse := httpClient.SendRequest(
		&roomresutils.IntegrationHttpRequest{
			Method:        "GET",
			UrlParameters: urlParameters,
		})

	var bookResponse BookResponse

	fmt.Printf("httpResponse = %s\n", httpResponse)

	errResponse := xml.Unmarshal(
		httpResponse.ResponseBody,
		&bookResponse)

	if errResponse != nil {
		fmt.Printf(" SendBookRequest - deserialization : %s\n", fmt.Sprintf("%v", bookingRequest))
		fmt.Printf(" SendBookRequest - deserialization : %s\n", fmt.Sprintf("%s", httpResponse))
		fmt.Printf(" SendBookRequest - deserialization : %s\n", errResponse.Error())
	}

	return &bookResponse
}

func (m *WebBedClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {
	// Hu B code here
	cancelResponse := m.SendCancelRequest(bookingCancelRequest)

	hbeCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_cancelresponse_to_hbe(cancelResponse, bookingCancelRequest, hbeCancelResponse)

	return hbeCancelResponse
}

func (m *WebBedClient) Init() {

}

func (m *WebBedClient) SendCancelRequest(bookingCancelRequest *hbecommon.BookingCancelRequest) *BookingCancelResponse {

	httpClient := roomresutils.NewIntegrationHttp(m.Settings.BookingCancelEndPoint, m.CreateHttpHeaders())

	var urlParameters string

	urlParameters = fmt.Sprintf("?userName=%s&password=%s&bookingID=%s&language=%s",
		m.Settings.UserName,
		m.Settings.Password,
		bookingCancelRequest.Ref,
		"en",
	)

	fmt.Printf("Request URL = %s\n", urlParameters)

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
func (m *WebBedClient) SendSearchRequests(searchRequests []*SearchRequest) *SearchResponse {

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

func (m *WebBedClient) SendSearchRequestChannel(closableChannel *ClosableChannel, searchRequest *SearchRequest) {

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

func (m *WebBedClient) SendPreBookRequests(preBookingRequests []*PreBookingRequest) []*PreBookResponse {

	//fmt.Printf("SendSearchRequests - Start %d\n", len(searchRequests))
	//var requestChannel chan *SearchResponse = make(chan *SearchResponse)
	closableChannel := NewClosableChannel()

	var responses []*PreBookResponse

	var i int = 0
	var maxBatches int = m.Settings.NumberOfThreads

	if len(preBookingRequests) < maxBatches {
		maxBatches = len(preBookingRequests)
	}

	for i < maxBatches {
		go m.SendPreBookingRequestChannel(closableChannel, preBookingRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.Settings.SmartTimeout

	for len(responses) < len(preBookingRequests) {

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

		if i < len(preBookingRequests) {
			go m.SendPreBookingRequestChannel(closableChannel, preBookingRequests[i])
			i++
		}
	}

	return responses
}

func (m *WebBedClient) SendPreBookingRequestChannel(closableChannel *ClosableChannel, preBookingRequest *PreBookingRequest) {

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

	response := m.SendPreBookRequest(preBookingRequest)

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- response
	})
}

func (m *WebBedClient) SendBookRequests(bookingRequests []*BookingRequest) []*BookResponse {

	closableChannel := NewClosableChannel()

	var responses []*BookResponse

	var i int = 0
	var maxBatches int = m.Settings.NumberOfThreads

	if len(bookingRequests) < maxBatches {
		maxBatches = len(bookingRequests)
	}

	for i < maxBatches {
		go m.SendBookingRequestChannel(closableChannel, bookingRequests[i])
		i++
	}

	var totalTime time.Duration

	var smartTimeout int = m.Settings.SmartTimeout

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

func (m *WebBedClient) SendBookingRequestChannel(closableChannel *ClosableChannel, bookingRequest *BookingRequest) {

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
