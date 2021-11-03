package main

import (
	"encoding/json"
	"fmt"
	"testing"

	derby "./derby"
	hbecommon "./roomres/hbe/common"
)

var hotelProviderSettingForDerby hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	Metadata: `{
		"DistributorId":"ROOMRES",
		"Version":"v4",
		"SupplierId" : "GOHOTEL",
		"Token":"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJST09NUkVTIiwibmFtZSI6IlJvb20tUmVzLmNvbSIsInRpbWUiOiIxNTQ3NTIxOTg3In0.lgvotCu8pxGj934aaNoAYwpKkqH0OWPefhlwDFnn1Yc"
	}`,
	SearchEndPoint:              "https://solo.derbysoft-test.com/bookingusbv4/availability",
	PreBookEndPoint:             "https://solo.derbysoft-test.com/bookingusbv4/reservation/prebook",
	BookingConfirmationEndPoint: "https://solo.derbysoft-test.com/bookingusbv4/reservation/book",
	BookingCancellationEndPoint: "https://solo.derbysoft-test.com/bookingusbv4/reservation/cancel",
	BookingDetailsEndpoint:      "https://solo.derbysoft-test.com/bookingusbv4/reservation/detail",
	HotelsPerBatch:              3,
	NumberOfThreads:             3,
	SmartTimeout:                10000,
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VI",
		Number:     "4242424242424242",
		ExpiryDate: "0123",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

func TestHotelSearchForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()
	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":2,"children":0}],
		"HotelIds":["GOH201", "GOH202"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

func TestMultiHotel1ForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()
	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":2,"children":0}],
		"HotelIds":["GOH201", "GOH202", "GOH203", "GOH101", "GOH102", "GOH103"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdultsForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-22",
		"CheckOut":"2019-01-24",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdultsForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildrenForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]},{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-23",
		"CheckOut":"2019-01-25",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":6}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHDE Search Result = %s\n", resultStr)
}

//Scenario 3, Full specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefsForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["GOH201"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"roomId\":\"T2\",\"rateId\":\"PROMO\",\"currency\":\"EUR\"}"],
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

	searchResponse := derbyClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBookingForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2019-01-21",
		"CheckOut":"2019-01-23",
		"Total":216,
		"Hotel":{
			"HotelId":"GOH201",
			"Rooms":[{
				"Ref":"{\"roomId\":\"T2\",\"rateId\":\"PROMO\",\"currency\":\"EUR\",\"amountBeforeTax\":null,\"amountAfterTax\":[108,108]}",
				"Count":1,
				"Adults":2,
				"Children":0,
				"Guests":[{
					"FirstName":"Li",
					"LastName":"Xing",
					"IsAdult":true,
					"Age":30
				},{
					"FirstName":"Li",
					"LastName":"Tian",
					"IsAdult":true,
					"Age":27
				}]
			}]
		},
		"Customer":{
			"FirstName":"Li",
			"LastName":"Xing",
			"Email":"polarislee1984@outlook.com",
			"PhoneCountryCode":"86",
			"PhoneNumber":"135546547"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := derbyClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingForDerby(t *testing.T) {
	derbyClient := derby.NewHotelBookingProvider(&hotelProviderSettingForDerby)
	derbyClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"C2084DFL0",
		"InternalRef":"D15F893D34DF"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := derbyClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Cancel Bookinig Result = %s\n", resultStr)
}
