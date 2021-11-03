package main

import (
	"encoding/json"
	"fmt"
	"testing"

	hbecommon "./roomres/hbe/common"
	"./teamamerica"
)

var taHotelProviderSetting hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	Metadata: `{
		"BaseEndPoint":"http://javatest.teamamericany.com:8080/TADoclit/services/TADoclit",
		"UserName":"XMLROOMRESAU",
		"Password":"9VA2fH8m"
	}`,
	HotelsPerBatch:  3,
	NumberOfThreads: 3,
	SmartTimeout:    5000,
	CreditCardInfo: &hbecommon.CreditCardInfo{
		CardType:   "VISA",
		Number:     "4242424242424242",
		ExpiryDate: "20201020",
		Cvc:        "123",
		HolderName: "LiXing",
	},
}

//two hotels
func TestMultiHotelForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-21",
		"Rooms":[{"adults":1,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|"],
		"ExternalRefs":null,
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdultsForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdultsForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildrenForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRoomsForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":1,"children":1},{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//many hotels, 1 city and multi hotels
//Added by Li, 2018/10/18
func TestMultiHotels1ForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-22",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|", "ta_16964|NYC|", "ta_158|NYC|", "ta_15963|NYC|", "ta_15960|NYC|", "ta_15968|NYC|", "ta_1441|NYC|", "ta_8258|NYC|", "ta_15975|NYC|", "ta_15962|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//many hotels, 3 city and multihotels
//Added by Li, 2018/10/18
func TestMultiHotels2ForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-22",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|", "ta_2|NYC|", "ta_16964|NYC|", "ta_158|NYC|", "ta_8113|OAK|", "ta_9032|OAK|", "ta_16701|PBI|", "ta_16019|PBI|", "ta_8027|PBI|", "ta_7978|PBI|", "ta_8139|PBI|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypesForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefsForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-22",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_131|NYC|"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"ProductCode\":\"NYCHPCGGB1\"}"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs:MaxPrice
func Test2WithSpecialRoomRefsForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":1,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["ta_2|NYC|"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"MaxPrice\":800.0}"],
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

	searchResponse := teamamericaClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBookingForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	jsonBookingRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"InternalRef":"",
		"Hotel":{
			"HotelId":"ta_131|NYC|",
			"Rooms":[{
				"Ref":"{\"ProductCode\":\"NYCHPCGGB1\",\"ProductDate\":\"2018-09-20\",\"MealPlan\":\"GRAB \\u0026 GO BREAKFAST\",\"RoomType\":\"CLASSIC 1 BED WITH FACILITY FEE INCLUDED!\",\"ChildAge\":18,\"FamilyPlan\":\"Y\",\"NonRefundable\":0,\"MaxOccupancy\":2,\"AverageNightlyRate\":410.45}",
				"Count":1,
				"Adults":1,
				"Children":1,
				"SpecialRequest":"",
				"Guests":[{
					"Primary":true,
					"Title":"Mr",
					"FirstName":"Li",
					"LastName":"Xing",
					"IsAdult":true,
					"Age":30
				},{
					"Primary":false,
					"Title":"",
					"FirstName":"Li",
					"LastName":"Tian",
					"IsAdult":false,
					"Age":10
				}]
			}]
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := teamamericaClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingReservationForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"1043929"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := teamamericaClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Cancel Bookinig Result = %s\n", resultStr)
}

//Get Room Cxl
func TestGetRoomCxlForTA(t *testing.T) {
	teamamericaClient := teamamerica.NewHotelBookingProvider(&taHotelProviderSetting)
	teamamericaClient.Init()

	//only need booking number
	jsonRoomCxlRequest := `{
		"RoomRef":"{\"ProductCode\":\"NYCHGCPREM\",\"ProductDate\":\"2018-09-20\",\"MealPlan\":\"EP-NO MEALS\",\"RoomType\":\"PREMIUM CORNER KING\",\"ChildAge\":0,\"FamilyPlan\":\"N\",\"NonRefundable\":0,\"MaxOccupancy\":2,\"AverageNightlyRate\":389.18}"
	}`

	var roomCxlRequest hbecommon.RoomCxlRequest
	err := json.Unmarshal([]byte(jsonRoomCxlRequest), &roomCxlRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	roomCxlResponse := teamamericaClient.GetRoomCxl(&roomCxlRequest)

	resultStr, err := json.Marshal(roomCxlResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHBE Get Room Cxl Result = %s\n", resultStr)
}
