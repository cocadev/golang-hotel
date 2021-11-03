package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"./hoteldo"
	hbecommon "./roomres/hbe/common"
)

var hotelProviderSettingForHD hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	SearchEndPoint:              "http://testxml.e-tsw.com/AffiliateService/AffiliateService.svc/restful/GetQuoteHotels",
	BookingConfirmationEndPoint: "http://testxml.e-tsw.com/AffiliateService/AffiliateService.svc/restful/Book",
	BookingCancellationEndPoint: "http://testxml.e-tsw.com/AffiliateService/AffiliateService.svc/restful/CancelItineraryHotel",
	PropertyDetailsEndpoint:     "http://testxml.e-tsw.com/AffiliateService/AffiliateService.svc/restful/GetQuoteHotels",
	UserName:                    "9010477",
	SecretKey:                   "Roomres*18",
	HotelsPerBatch:              3,
	NumberOfThreads:             3,
	SmartTimeout:                5000,
	ProfileCountry:              "US",
	ProfileCurrency:             "US",
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VISA",
		Number:     "4242424242424242",
		ExpiryDate: "20201020",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

//two hotels
func TestMultiHotelForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":1,"children":0}],
		"HotelIds":["4656", "10"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//multi hotels
//Added by Li, 2018/10/19
func TestMultiHotel1ForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":1,"children":0}],
		"HotelIds":["10", "100", "104", "1053", "1073", "1074", "1077", "109", "1123", "118", "1209", "128", "1400"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdultsForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["4656"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdultsForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["3938"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildrenForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["4656"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]},{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["4656"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":6}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["10"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHDE Search Result = %s\n", resultStr)
}

//Scenario 3, Full specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefsForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":6}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["10"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"roomType\":\"STU\",\"mealplan\":\"EP\",\"marketid\":\"MOBILE\",\"contract\":1,\"currency\":\"US\",\"amount\":117.12285,\"status\":\"AV\",\"ratekey\":\"STUEP\",\"adults\":\"1\",\"kids\":\"2\",\"k1a\":\"10\",\"k2a\":\"6\"}"],
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

	searchResponse := hdClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBookingForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	jsonBookingRequest := `{
		"InternalRef":"{\"type\":\"CREDMX\",\"currency\":\"US\",\"amount\":117.12285}",
		"CheckIn":"2018-10-21",
		"CheckOut":"2018-10-23",
		"Total":117.12285,
		"Hotel":{
			"HotelId":"10",
			"Rooms":[{
				"Ref":"{\"roomType\":\"STU\",\"mealplan\":\"EP\",\"marketid\":\"MOBILE\",\"contract\":1,\"currency\":\"US\",\"amount\":117.12285,\"status\":\"AV\",\"ratekey\":\"STUEP\",\"adults\":\"1\",\"kids\":\"2\",\"k1a\":\"10\",\"k2a\":\"6\"}",
				"Adults":1,
				"Children":2
			}]
		},
		"Customer":{
			"FirstName":"Li",
			"LastName":"Xing"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := hdClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingForHD(t *testing.T) {
	hdClient := hoteldo.NewHotelBookingProvider(&hotelProviderSettingForHD)
	hdClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"47121828"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := hdClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Cancel Bookinig Result = %s\n", resultStr)
}
