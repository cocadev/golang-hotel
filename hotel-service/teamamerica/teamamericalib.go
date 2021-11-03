package teamamerica

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"
	roomresutils "../roomres/utils"
	"golang.org/x/net/html/charset"
)

type TeamAmericaClient struct {
	Settings *TeamAmericaSettings

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

	CommissionProvider hbecommon.ICommissionProvider

	//RepositoryFactory  *repository.RepositoryFactory

	//Log roomresutils.ILog
}

type TeamAmericaSettings struct {
	BaseEndPoint   string `json:"BaseEndPoint"`
	UserName       string `json:"UserName"`
	Password       string `json:"Password"`
	CreditCardInfo *hbecommon.CreditCardInfo
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings) *TeamAmericaClient {

	teamAmericaClient := &TeamAmericaClient{

		//Log: logging,
		HotelsPerBatch:     hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads:    hotelProviderSettings.NumberOfThreads,
		SmartTimeout:       hotelProviderSettings.SmartTimeout,
		CommissionProvider: hotelProviderSettings.CommissionProvider,
		//RepositoryFactory:  hotelProviderSettings.RepositoryFactory,
	}

	var teamAmericaClientSettings TeamAmericaSettings
	err := json.Unmarshal([]byte(hotelProviderSettings.Metadata), &teamAmericaClientSettings)
	if err != nil {
		fmt.Printf("Error : TeamAmerica settings : %s", hotelProviderSettings.Metadata)
	}
	teamAmericaClient.Settings = &teamAmericaClientSettings
	teamAmericaClient.Settings.CreditCardInfo = hotelProviderSettings.CreditCardInfo

	return teamAmericaClient
}

func (m *TeamAmericaClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "text/xml",
	}
}

func (m *TeamAmericaClient) CreateHttpBodyParameters() map[string]string {
	return map[string]string{
		"UserName": m.Settings.UserName,
		"Password": m.Settings.Password,
	}
}

func (m *TeamAmericaClient) CreateHttpSearchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate",
	}
}

//Plugin APIs
func (m *TeamAmericaClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {
	//fmt.Printf("search request data = %+v\n", hotelSearchRequest)

	hotelSearchResponse := &hbecommon.HotelSearchResponse{}

	if CheckMaxNumberOfRooms(hotelSearchRequest) && CheckMaxNumberOfPax(hotelSearchRequest) && hotelSearchRequest.CheckTreshold(12) {
		requests := HbeToTeamAmericaHotelSearchRequests(hotelSearchRequest)
		//added by Li
		requests = GroupSearchRequestsByCity(requests)
		requests = SplitBatchSearchRequests(requests, m.HotelsPerBatch)
		response := m.SendSearchRequests(requests)

		roomCxlResponses := []*TeamAmericaCancelPolicyResponse{}

		if hotelSearchRequest.Details && len(hotelSearchRequest.SpecificRoomRefs) > 0 {
			//added li, 2018/10/20
			//filter rooms first
			var hotels []*HotelData

			for _, hotel := range response.HotelSearchResponse.HotelDatas {
				var total float32 = 0
				for _, nightInfo := range hotel.NightlyInfo {
					total += nightInfo.Prices.AdultPrice
				}

				var searchRoomRefs []*SearchRoomRef
				for _, roomRefStr := range hotelSearchRequest.SpecificRoomRefs {
					var searchRoomRef SearchRoomRef
					err := json.Unmarshal([]byte(roomRefStr), &searchRoomRef)
					if err != nil {
						fmt.Printf("Cannot parse specialRoomRef = %s\n", roomRefStr)
					} else {
						searchRoomRefs = append(searchRoomRefs, &searchRoomRef)
					}
				}
				//check special condition, let's consider specialRoomRef is one
				checkCondition := searchRoomRefs[0]
				if checkCondition.ProductCode != "" && checkCondition.ProductCode != hotel.ProductCode {
					continue
				}
				if checkCondition.ProductDate != "" && checkCondition.ProductDate != hotel.ProductDate {
					continue
				}
				if checkCondition.RoomType != "" && checkCondition.RoomType != hotel.RoomType {
					continue
				}
				if checkCondition.MealPlan != "" && checkCondition.MealPlan != hotel.MealPlan {
					continue
				}
				if checkCondition.ChildAge != 0 && checkCondition.ChildAge != hotel.ChildAge {
					continue
				}
				if checkCondition.FamilyPlan != "" && checkCondition.FamilyPlan != hotel.FamilyPlan {
					continue
				}
				if checkCondition.NonRefundable != 0 && checkCondition.NonRefundable != hotel.NonRefundable {
					continue
				}
				if checkCondition.MaxOccupancy != 0 && checkCondition.ProductCode != hotel.ProductCode {
					continue
				}
				if checkCondition.MinPrice != 0 && checkCondition.MinPrice > total {
					continue
				}
				if checkCondition.MaxPrice != 0 && checkCondition.MaxPrice < total {
					continue
				}

				hotels = append(hotels, hotel)
			}
			response.HotelSearchResponse.HotelDatas = hotels

			if response.HotelSearchResponse.HotelDatas != nil && len(response.HotelSearchResponse.HotelDatas) > 0 {
				//call room cxl policy request only when case 3
				var productCodes []string
				for _, hotel := range response.HotelSearchResponse.HotelDatas {
					productCodes = append(productCodes, hotel.ProductCode)
				}
				roomCxlResponses = m.SendCancelPolicyRequests(productCodes)
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

func (m *TeamAmericaClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {

	bookRequest := &TeamAmericaBookReserveRequest{}
	makeBookingRequest(hbeBookingRequest, bookRequest)

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	roomCxlResponses := []*TeamAmericaCancelPolicyResponse{}

	items := bookRequest.Items
	if items != nil && len(items.NewItems) > 0 {
		//call room cxl policy request only when case 3
		var productCodes []string
		for _, item := range items.NewItems {
			productCodes = append(productCodes, item.ProductCode)
		}
		roomCxlResponses = m.SendCancelPolicyRequests(productCodes)
		//fmt.Printf("roomCxlResponses=%+v\n", roomCxlResponses)
	}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse, roomCxlResponses)

	return hbeBookingResponse
}

func (m *TeamAmericaClient) CancelBooking(hbeBookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {

	bookCancelRequest := &TeamAmericaCancelReservationRequest{}
	bookCancelRequest.ReservationNumber = hbeBookingCancelRequest.Ref

	bookCancelResponse := m.SendCancelReservationRequest(bookCancelRequest)

	hbeBookingCancelResponse := &hbecommon.BookingCancelResponse{}

	mapping_cancelresponse_to_hbe(bookCancelResponse, hbeBookingCancelResponse)

	return hbeBookingCancelResponse
}

func (m *TeamAmericaClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {
	roomCancelRequest := &TeamAmericaCancelPolicyRequest{}
	roomRef := &RoomRef{}
	err := json.Unmarshal([]byte(roomCxlRequest.RoomRef), roomRef)
	if err != nil {
		return nil
	}
	roomCancelRequest.ProductCode = roomRef.ProductCode

	roomCancelResponse := m.SendCancelPolicyRequest(roomCancelRequest)

	hbeRoomCancelResponse := &hbecommon.RoomCxlResponse{}

	mapping_roomcancelresponse_to_hbe(roomCancelResponse, hbeRoomCancelResponse, roomRef.ProductDate)

	return hbeRoomCancelResponse
}

// func (m *TeamAmericaClient) GetBookingInfo(request *hbecommon.BookingInfoRequest) *hbecommon.BookingInfoResponse {

// 	bookingInfoRequest := &TeamAmericaRetrieveReservationRequest{}
// 	bookingInfoRequest.ReservationNumber = request.BookingId

// 	bookInfoResponse := m.SendBookingInfoRequest(bookingInfoRequest)

// 	hbeBookingInfoResponse := &hbecommon.BookingInfoResponse{}

// 	if bookInfoResponse.BookingInfoDetail != nil {
// 		hbeBookingInfoResponse.HotelSupplierId = "" //not sure what could mapping here
// 		hbeBookingInfoResponse.HotelPhoneNumber = bookInfoResponse.BookingInfoDetail.Phone
// 	}

// 	return hbeBookingInfoResponse
// }

// func (m *TeamAmericaClient) SendSpecialRequest(bookingSpecialRequestRequest *hbecommon.BookingSpecialRequestRequest) *hbecommon.BookingSpecialRequestResponse {
// 	return nil
// }

// func (m *TeamAmericaClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {
// 	return nil
// }

// func (m *TeamAmericaClient) ListBookingReport(request *hbecommon.BookingReportRequest) *hbecommon.BookingReportResponse {
// 	return nil
// }

func (m *TeamAmericaClient) Init() {

}

func (m *TeamAmericaClient) SendSearchRequests(searchRequests []*TeamAmericaHotelSearchRequest) *TeamAmericaHotelSearchResponse {

	//var requestChannel chan *AvailabilityResponse = make(chan *AvailabilityResponse)
	closableChannel := roomresutils.NewClosableChannel()

	var responses []*TeamAmericaHotelSearchResponse

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
			responses = append(responses, response.(*TeamAmericaHotelSearchResponse))
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

	combinedSearchResponse := &TeamAmericaHotelSearchResponse{HotelSearchResponse: &HotelSearchResponse{
		HotelDatas: []*HotelData{},
	}}

	for _, response := range responses {

		if response != nil && response.HotelSearchResponse.HotelDatas != nil {
			combinedSearchResponse.HotelSearchResponse.HotelDatas =
				append(
					combinedSearchResponse.HotelSearchResponse.HotelDatas,
					response.HotelSearchResponse.HotelDatas...,
				)
		}
	}

	return combinedSearchResponse
}

func (m *TeamAmericaClient) SendSearchRequestChannel(closableChannel *roomresutils.ClosableChannel, searchRequest *TeamAmericaHotelSearchRequest) {

	defer func() {
		if err := recover(); err != nil {

			fmt.Printf("SendSearchRequestChannel Error : %+v\n", err)

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

func (m *TeamAmericaClient) SendSearchRequest(searchRequest *TeamAmericaHotelSearchRequest) *TeamAmericaHotelSearchResponse {
	fmt.Println("SendSearchRequest ----------- Start ----------")

	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	searchRequest.UserName = m.Settings.UserName
	searchRequest.Password = m.Settings.Password

	fullRequestData := TeamAmericaRequest{
		SoapBody: &SoapBody{
			BodyData: searchRequest,
		},
		SoapEnv:    "http://schemas.xmlsoap.org/soap/envelope/",
		Xsd:        "http://www.wso2.org/php/xsd",
		SoapHeader: &SoapHeader{},
	}

	xmlStr, _ := xml.Marshal(fullRequestData)

	//fmt.Printf("Request Data = %s\n", xmlStr)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          []byte(xmlStr),
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var fullResponse *SearchResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&fullResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	searchResponse := fullResponse.SoapBody

	//setting compound hotelId
	cityCode := searchRequest.CityCode
	productCode := searchRequest.ProductCode

	for _, hotel := range searchResponse.HotelSearchResponse.HotelDatas {
		hotel.CompoundHotelId = fmt.Sprintf("%d|%s|%s", hotel.TeamVendorID, cityCode, productCode)
	}
	return searchResponse
}

func (m *TeamAmericaClient) SendCancelReservationRequest(cancelRequest *TeamAmericaCancelReservationRequest) *TeamAmericaCancelReservationResponse {
	fmt.Println("SendCancelReservationRequest ----------- Start ----------")

	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	cancelRequest.UserName = m.Settings.UserName
	cancelRequest.Password = m.Settings.Password

	fullRequestData := TeamAmericaRequest{
		SoapBody: &SoapBody{
			BodyData: cancelRequest,
		},
		SoapEnv:    "http://schemas.xmlsoap.org/soap/envelope/",
		Xsd:        "http://www.wso2.org/php/xsd",
		SoapHeader: &SoapHeader{},
	}

	xmlStr, _ := xml.Marshal(fullRequestData)

	//fmt.Printf("Request Data = %s\n", xmlStr)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          []byte(xmlStr),
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var fullResponse *CancelReservationResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&fullResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	responseBody := fullResponse.SoapBody
	return responseBody
}

func (m *TeamAmericaClient) SendBookRequest(bookRequest *TeamAmericaBookReserveRequest) *TeamAmericaBookReserveResponse {
	fmt.Println("SendBookRequest ----------- Start ----------")

	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	bookRequest.UserName = m.Settings.UserName
	bookRequest.Password = m.Settings.Password

	fullRequestData := TeamAmericaRequest{
		SoapBody: &SoapBody{
			BodyData: bookRequest,
		},
		SoapEnv:    "http://schemas.xmlsoap.org/soap/envelope/",
		Xsd:        "http://www.wso2.org/php/xsd",
		SoapHeader: &SoapHeader{},
	}

	xmlStr, _ := xml.Marshal(fullRequestData)

	//fmt.Printf("Request Data = %s\n", xmlStr)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          []byte(xmlStr),
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var fullResponse *BookReserveResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&fullResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	bookResponse := fullResponse.SoapBody
	fmt.Printf("Booking Reservation count = %+v\n", fullResponse.SoapBody)
	return bookResponse
}

func (m *TeamAmericaClient) SendCancelPolicyRequest(cancelPolicyRequest *TeamAmericaCancelPolicyRequest) *TeamAmericaCancelPolicyResponse {
	fmt.Println("SendBookRequest ----------- Start ----------")

	httpClient := NewIntegrationHttp(m.Settings.BaseEndPoint, m.CreateHttpHeaders())

	cancelPolicyRequest.UserName = m.Settings.UserName
	cancelPolicyRequest.Password = m.Settings.Password

	fullRequestData := TeamAmericaRequest{
		SoapBody: &SoapBody{
			BodyData: cancelPolicyRequest,
		},
		SoapEnv:    "http://schemas.xmlsoap.org/soap/envelope/",
		Xsd:        "http://www.wso2.org/php/xsd",
		SoapHeader: &SoapHeader{},
	}

	xmlStr, _ := xml.Marshal(fullRequestData)

	//fmt.Printf("Request Data = %s\n", xmlStr)

	httpResponse := httpClient.SendRequest(
		&IntegrationHttpRequest{
			Method:               "POST",
			UrlParameters:        "",
			RequestBodySpecified: true,
			RequestBody:          []byte(xmlStr),
		})

	if httpResponse.Err != nil {
		fmt.Printf("Error:%+v\n", httpResponse.Err)
		return nil
	}

	//fmt.Printf("Response = %s\n", httpResponse.ResponseBody)
	var fullResponse *CancelPolicyResponse

	decoder := xml.NewDecoder(strings.NewReader(string(httpResponse.ResponseBody)))
	decoder.CharsetReader = charset.NewReaderLabel
	errResponse := decoder.Decode(&fullResponse)

	if errResponse != nil {
		fmt.Printf("error = %+v, %s\n", errResponse, errResponse)
		return nil
	}

	response := fullResponse.SoapBody

	//allocate productCode
	response.CancellationPolicyResponse.ProductCode = cancelPolicyRequest.ProductCode

	return response
}

func (m *TeamAmericaClient) SendCancelPolicyRequests(searchRequests []string) []*TeamAmericaCancelPolicyResponse {

	//var requestChannel chan *AvailabilityResponse = make(chan *AvailabilityResponse)
	closableChannel := roomresutils.NewClosableChannel()

	var responses []*TeamAmericaCancelPolicyResponse

	var i int = 0
	var maxBatches int = m.NumberOfThreads

	if len(searchRequests) < maxBatches {
		maxBatches = len(searchRequests)
	}

	for i < maxBatches {
		go m.SendCancelPolicyRequestChannel(closableChannel, searchRequests[i])
		i++
	}

	for len(responses) < len(searchRequests) {

		exit := false

		select {
		case response := <-closableChannel.Channel:
			responses = append(responses, response.(*TeamAmericaCancelPolicyResponse))
		case <-time.After(time.Second * 20):
			exit = true
			closableChannel.Close()
			break
		}

		if exit {
			break
		}

		if i < len(searchRequests) {
			go m.SendCancelPolicyRequestChannel(closableChannel, searchRequests[i])
			i++
		}
	}

	return responses
}

func (m *TeamAmericaClient) SendCancelPolicyRequestChannel(closableChannel *roomresutils.ClosableChannel, searchRequest string) {

	defer func() {
		if err := recover(); err != nil {

			fmt.Printf("SendSearchRequestChannel Error : %+v\n", err)

			closableChannel.Execute(func(channel chan interface{}) {
				channel <- nil
			})
		}
	}()

	searchResponse := m.SendCancelPolicyRequest(&TeamAmericaCancelPolicyRequest{
		ProductCode: searchRequest,
	})

	closableChannel.Execute(func(channel chan interface{}) {
		channel <- searchResponse
	})
}
