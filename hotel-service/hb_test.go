package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"./hb"

	hbecommon "./roomres/hbe/common"
	utils "./roomres/utils"
)

var hotelProviderSettingForHB hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	SearchEndPoint:                "https://api.test.hotelbeds.com/hotel-api/1.0/hotels",
	BookingSpecialRequestEndPoint: "https://api.test.hotelbeds.com/hotel-api/1.0/checkrates",
	BookingConfirmationEndPoint:   "https://api.test.hotelbeds.com/hotel-api/1.0/bookings",
	BookingDetailsEndpoint:        "https://api.test.hotelbeds.com/hotel-api/1.0/bookings",
	PropertyDetailsEndpoint:       "https://api.test.hotelbeds.com/hotel-content-api/1.0",
	ProviderId:                    0,
	ApiKey:                        "h55jva9hu38jkpgc3xgk8ges",
	SecretKey:                     "vxMbEBDBAx",
	HotelsPerBatch:                3,
	NumberOfThreads:               3,
	SmartTimeout:                  5000,
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VISA",
		Number:     "4242424242424242",
		ExpiryDate: "20201020",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

//two hotels
func TestMultiHotelForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-21",
		"CheckOut":"2018-09-23",
		"Rooms":[{"adults":1,"children":0}],
		"HotelIds":["330352","232197"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdultsForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["160022"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdultsForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["160022"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildrenForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["160022"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]},{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["160022"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRooms2ForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-01",
		"CheckOut":"2018-09-03",
		"Rooms":[{"adults":2,"children":0},{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["160022"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-25",
		"CheckOut":"2018-09-28",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["330352"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHBE Search Result = %s\n", resultStr)
}

//Scenario 3, Full specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefsForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-25",
		"CheckOut":"2018-09-28",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["330352"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["20180925|20180928|W|255|330352|DBL.KG|FIT|DB||1~2~0||N@353B644997404F59AAE281B11D50C0971258"],
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

	searchResponse := hbClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBookingForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	jsonBookingRequest := `{
		"Hotel":{
			"Rooms":[{
				"Ref":"20180925|20180926|W|254|323357|QUA.ST|NRF-ALL|RO||1~2~0||N@3D46A77064054E45814F49F562FC04E11737",
				"Adults":1,
				"Children":0,
				"Guests":[{
					"FirstName":"Li",
					"LastName":"Xing",
					"Gender":1,
					"IsAdult":true
				}]
			}]
		},
		"Customer":{
			"FirstName":"IntegrationTestFirstName",
			"LastName":"IntegrationTestLastName"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := hbClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingForHB(t *testing.T) {
	hbClient := hb.NewHotelBookingProvider(&hotelProviderSettingForHB, utils.CreateRootLogScope())
	hbClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"254-1888845"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := hbClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Cancel Bookinig Result = %s\n", resultStr)
}
