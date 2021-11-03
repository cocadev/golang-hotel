package main

import (
	"encoding/json"
	"fmt"
	"testing"

	hbecommon "./roomres/hbe/common"
	"./webbed"
)

var hotelProviderSettingForWD hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	SearchEndPoint:              "http://search.fitruums.com/1/PostGet/StaticXMLAPI.asmx/Search",
	PreBookEndPoint:             "http://book.fitruums.com/1/PostGet/Booking.asmx/PreBook",
	BookingConfirmationEndPoint: "http://book.fitruums.com/1/PostGet/Booking.asmx/Book",
	BookingCancellationEndPoint: "http://book.fitruums.com/1/PostGet/Booking.asmx/CancelBooking",
	UserName:                    "TestRoomRes",
	SecretKey:                   "Test1234",
	HotelsPerBatch:              3,
	NumberOfThreads:             3,
	SmartTimeout:                10000,
	ProfileCountry:              "gb",
	ProfileCurrency:             "USD",
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VISA",
		Number:     "4242424242424242",
		ExpiryDate: "20201020",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

//two hotels
func TestHotelSearchForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":2,"children":0}],
		"HotelIds":["wbb_28038", "wbb_226618"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)
	fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

func TestMultiHotel1ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":1,"children":0}],
		"HotelIds":["wbb_28038", "wbb_226618", "wbb_257196", "wbb_6184", "wbb_46083", "wbb_132500", "wbb_24591", "wbb_9156", "wbb_17630", "wbb_1036", "wbb_28672", "wbb_45830", "wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)
	//fmt.Printf("searchResponse = %d\n", len(searchResponse.Hotels))
	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Two adults
func Test2_0_0ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-22",
		"CheckOut":"2018-11-24",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_28038"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//3 adults
func Test_3_0_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 1children
func Test1_1_0ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 2children
func Test1_2_0ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":8}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 adults 1children
func Test_2_1_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":2,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 adults 2children
func Test_2_2_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":2,"children":2,"childages":[{"age":10},{"age":8}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 3children
func Test_1_3_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":3,"childages":[{"age":10},{"age":8},{"age":4}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 1children 1 baby
func Test1_1_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":1}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 2children 1baby
func Test1_2_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":3,"childages":[{"age":10},{"age":8},{"age":1}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 adults 1children 1bady
func Test_2_1_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":2,"children":2,"childages":[{"age":10},{"age":1}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 adults 2children 1baby
func Test_2_2_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":2,"children":3,"childages":[{"age":10},{"age":8},{"age":1}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//1 adults 3children 1 baby
func Test_1_3_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Rooms":[{"adults":1,"children":4,"childages":[{"age":10},{"age":8},{"age":4},{"age":1}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":1,"children":1,"childages":[{"age":10}]},{"adults":1,"children":1,"childages":[{"age":10}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_24591"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":1,"children":2,"childages":[{"age":10},{"age":6}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_2376"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHDE Search Result = %s\n", resultStr)
}

//Scenario 3, Full specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefsForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["wbb_30518"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["2931663"],
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

	searchResponse := webbedClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBooking_2_0_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-21",
		"CheckOut":"2018-11-23",
		"Hotel":{
			"Rooms":[{
				"Ref":"2931663",
				"Count":1,
				"Adults":2,
				"Children":0,
				"SpecialRequest":"{\"mealId\":\"3\"}",
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
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

func TestMakeBooking_3_0_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Hotel":{
			"Rooms":[{
				"Ref":"12912493",
				"Count":1,
				"Adults":3,
				"Children":0,
				"SpecialRequest":"{\"mealId\":\"3\"}",
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
				},{
					"FirstName":"Wang",
					"LastName":"Hong",
					"IsAdult":true,
					"Age":20
				}]
			}]
		},
		"Customer":{
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

func TestMakeBooking_2_1_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Hotel":{
			"Rooms":[{
				"Ref":"12912493",
				"Count":1,
				"Adults":2,
				"Children":1,
				"SpecialRequest":"{\"mealId\":\"3\"}",
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
				},{
					"FirstName":"Wang",
					"LastName":"Hong",
					"IsAdult":false,
					"Age":14
				}]
			}]
		},
		"Customer":{
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

func TestMakeBooking_2_2_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Hotel":{
			"Rooms":[{
				"Ref":"6241976",
				"Count":1,
				"Adults":2,
				"Children":2,
				"SpecialRequest":"{\"mealId\":\"1\"}",
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
				},{
					"FirstName":"Li",
					"LastName":"Z",
					"IsAdult":false,
					"Age":3
				},{
					"FirstName":"Li",
					"LastName":"J",
					"IsAdult":false,
					"Age":2
				}]
			}]
		},
		"Customer":{
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

func TestMakeBooking_1_2_0_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Hotel":{
			"Rooms":[{
				"Ref":"12912493",
				"Count":1,
				"Adults":1,
				"Children":2,
				"SpecialRequest":"{\"mealId\":\"1\"}",
				"Guests":[{
					"FirstName":"Li",
					"LastName":"Xing",
					"IsAdult":true,
					"Age":30
				},{
					"FirstName":"Li",
					"LastName":"Tian",
					"IsAdult":false,
					"Age":7
				},{
					"FirstName":"Li",
					"LastName":"Y",
					"IsAdult":false,
					"Age":2
				}]
			}]
		},
		"Customer":{
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

func TestMakeBooking_2_1_1_ForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-11-24",
		"CheckOut":"2018-11-25",
		"Hotel":{
			"Rooms":[{
				"Ref":"9972092",
				"Count":1,
				"Adults":2,
				"Children":2,
				"SpecialRequest":"{\"mealId\":\"1\"}",
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
				},{
					"FirstName":"Li",
					"LastName":"E",
					"IsAdult":false,
					"Age":7
				},{
					"FirstName":"Li",
					"LastName":"Y",
					"IsAdult":false,
					"Age":1
				}]
			}]
		},
		"Customer":{
			"Email":"polarislee1984@outlook.com"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := webbedClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingForWD(t *testing.T) {
	webbedClient := webbed.NewHotelBookingProvider(&hotelProviderSettingForWD)
	webbedClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"SH6637436"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := webbedClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HDE Cancel Bookinig Result = %s\n", resultStr)
}
