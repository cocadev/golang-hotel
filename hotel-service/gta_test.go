package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"./gta"

	hbecommon "./roomres/hbe/common"
	utils "./roomres/utils"
)

var hotelProviderSettingForGTA hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	SearchEndPoint:              "https://interface.demo.gta-travel.com/wbsapi/RequestListenerServlet",
	BookingConfirmationEndPoint: "https://interface.demo.gta-travel.com/wbsapi/RequestListenerServlet",
	BookingDetailsEndpoint:      "https://interface.demo.gta-travel.com/wbsapi/RequestListenerServlet",
	PropertyDetailsEndpoint:     "https://interface.demo.gta-travel.com/wbsapi/RequestListenerServlet",
	UserName:                    "XML.DPLOY@ROOMRES.COM",
	SiteId:                      "2027",
	ApiKey:                      "PASS",
	ProfileCurrency:             "USD",
	ProfileCountry:              "US",
	HotelsPerBatch:              30,
	NumberOfThreads:             5,
	SmartTimeout:                5000,
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VISA",
		Number:     "4242424242424242",
		ExpiryDate: "20201020",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

//two hotels
func TestMultiHotelForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-20",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":1,"children":0}],
		"HotelIds":["AMS;ACO","AMS;BEL","AMS;COM"],
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)
	//fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdultsForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-23",
		"Rooms":[{"adults":2,"children":0}],
		"HotelIds":["AMS;ACO","AMS;BEL","AMS;COM"],
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdultsForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["AMS;ACO","AMS;BEL","AMS;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildrenForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"HotelIds":["AMS;ACO","AMS;BEL","AMS;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]},{"adults":1,"children":1,"childages":[{"age":10}]}],
		"HotelIds":["AMS;ACO","AMS;BEL","AMS;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRooms2ForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-10",
		"CheckOut":"2018-09-11",
		"Rooms":[{"adults":2,"children":0},{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["NCE;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":null,
		"Details":false,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":2,"children":0}],
		"HotelIds":["NCE;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":null,
		"Details":true,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHBE Search Result = %s\n", resultStr)
}

//Scenario 3
func TestOneHotelWithSpecialRoomRefsForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["NCE;COM"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["001:COM:9848:S9660:11135:45854"],
		"Details":true,
		"Packaging":false,
		"SortType":"",
		"AutocompleteId":"",
		"HotelFilterId":"",
		"RecommendationOnly":false,
		"MaxPrice":0,
		"MinPrice":0,
		"StarRatings":null,
		"PageIndex":1,
		"PageSize":10}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := gtaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBookingForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	jsonBookingRequest := `{
		"Total":497,
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Hotel":{
			"HotelId":"NCE;COM",
			"Rooms":[{
				"Ref":"001:COM:9848:S9660:11135:45835",
				"Adults":2,
				"Children":0,
				"Guests":[{
					"FirstName":"Li",
					"LastName":"Xing"
				},{
					"FirstName":"Zhang",
					"LastName":"Hong"
				}]
			}]
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := gtaClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingForGTA(t *testing.T) {
	gtaClient := gta.NewHotelBookingProvider(&hotelProviderSettingForGTA, utils.CreateRootLogScope())
	gtaClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"B2018090509-e4c944bc"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := gtaClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Cancel Bookinig Result = %s\n", resultStr)
}
