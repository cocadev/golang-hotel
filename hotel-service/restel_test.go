package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"./restel"

	hbecommon "./roomres/hbe/common"
)

var hotelProviderSetting hbecommon.HotelProviderSettings = hbecommon.HotelProviderSettings{
	Metadata: `{
		"BaseEndPoint":"http://xml.hotelresb2b.com/xml/listen_xml.jsp",
		"UserCode":"TJYU",
		"UserPassword":"xml490952",
		"Client":"RS",
		"AccessCode":"123477",
		"AgencyAffiliation":"HA",
		"AgencyUserCode":"LOVI3"
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
func TestMultiHotel(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":1,"cots":0,"childages":[{"age":8}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//added by Li, 2018/10/19
//10 hotels
func TestMulti1Hotel(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":1,"cots":0,"childages":[{"age":8}]}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN", "rtl_001578|ES|ESBCN", "rtl_450704|ES|ESBCN", "rtl_973655|ES|ESBCN", "rtl_623391|ES|ESBCN", "rtl_667479|ES|ESBCN", "rtl_667485|ES|ESBCN", "rtl_796035|ES|ESBCN", "rtl_115046|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Two adults
func TestTwoAdults(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//3 adults
func TestThreeAdults(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":3,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//1 adults 1children
func TestOneAdultsOneChildren(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-21",
		"Rooms":[{"adults":1,"children":1}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//2 rooms
func TestTwoRooms(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-08-20",
		"CheckOut":"2018-08-21",
		"Rooms":[{"adults":2,"children":0},{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN", "rtl_600515|ES|ESBCN"],
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

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 2, Added details parameter
func TestOneHotelWithDetailedRoomTypes(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, Full specialRoomRefs parameters
func TestOneHotelWithSpecialRoomRefs(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"RoomType\":\"9+\",\"MealPlanType\":\"BB\",\"MinPrice\":\"0.0\",\"MaxPrice\":\"200.0\",\"RefundableType\":\"0\"}"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs:MaxPrice
func Test2WithSpecialRoomRefs(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"MaxPrice\":\"130.0\"}"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs:Room Type
func Test3WithSpecialRoomRefs(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-11-20",
		"CheckOut":"2018-11-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"RoomType\":\"DB\"}"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs:RefundableType
func Test4WithSpecialRoomRefs(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-21",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_626886|ES|ESBCN"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"RefundableType\":\"0\"}"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Scenario 3, specialRoomRefs:MealPlanType
func Test5WithSpecialRoomRefs(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonSearchRequest := `{
		"CheckIn":"2018-09-20",
		"CheckOut":"2018-09-22",
		"Rooms":[{"adults":2,"children":0}],
		"Lat":0,
		"Lon":0,
		"HotelIds":["rtl_450441#|ES|ESEAS"],
		"CurrencyCode":"",
		"SpecificRoomRefs":["{\"MealPlanType\":\"BB\"}"],
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
		"PageSize":10,
		"Details":true}`

	var searchRequest hbecommon.HotelSearchRequest
	err := json.Unmarshal([]byte(jsonSearchRequest), &searchRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	searchResponse := restelClient.SearchRequest(&searchRequest)

	resultStr, err := json.Marshal(searchResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Search Result = %s\n", resultStr)
}

//Booking
func TestMakeBooking(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	jsonBookingRequest := `{
		"Total":450.0,
		"InternalRef":"{\"Remark\":\"Sea views\",\"InternalUse\":\"\",\"PaymentForm\":12}",
		"Hotel":{
			"HotelId":"rtl_450441|ES|ESEAS",
			"Rooms":[{
				"Ref":"{\"RoomType\":\"D8\",\"RoomDesc\":\"DOBLE STANDARD\",\"RoomLine\":{\"MealPlanType\":\"BB\",\"PlanPrice\":\"224.66\",\"Currency\":\"DO\",\"Status\":\"OK\",\"MinPrice\":\"224.66\",\"FeeRate\":\"0\",\"CompressedAvailability\":[\"D8#1#VR#80.25#0#OB#OK#20180920#20180921#DO#2-0#75735#10#201809060403#450441#\",\"D8#1#VR#80.25#0#OB#OK#20180921#20180922#DO#2-0#75735#10#201809060403#450441#\",\"D8#1#VR#112.33#0#BB#OK#20180920#20180921#DO#2-0#0#0#201809060403#450441#\",\"D8#1#VR#112.33#0#BB#OK#20180921#20180922#DO#2-0#0#0#201809060403#450441#\"]}}",
				"Count":1,
				"Adults":2,
				"Children":0,
				"SpecialRequest":"",
				"Guests":[{
					"Primary":true,
					"Title":"Mr",
					"FirstName":"Li",
					"LastName":"Xing",
					"CountryOfPassport":"China",
					"CountryOfNationality":"China",
					"Gender":1,
					"IsAdult":true,
					"Age":30
				},{
					"Primary":false,
					"Title":"Mrss",
					"FirstName":"Wang",
					"LastName":"Xu",
					"CountryOfPassport":"China",
					"CountryOfNationality":"China",
					"Gender":2,
					"IsAdult":true,
					"Age":26
				}]
			}]
		},
		"Customer":{
			"Title":"Mr.",
			"FirstName":"Li",
			"LastName":"Xing",
			"Email":"polarislee1984@outlook.com",
			"PhoneCountryCode":"86",
			"PhoneAreaCode":"130",
			"PhoneNumber":"1233432535"
		}}`

	var bookingRequest hbecommon.BookingRequest
	err := json.Unmarshal([]byte(jsonBookingRequest), &bookingRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingResponse := restelClient.MakeBooking(&bookingRequest)

	resultStr, err := json.Marshal(bookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Reservation
func TestCancelBookingReservation(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	//only need booking number
	jsonBookingCancelRequest := `{
		"Ref":"01999731"
	}`

	var bookingCancelRequest hbecommon.BookingCancelRequest
	err := json.Unmarshal([]byte(jsonBookingCancelRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := restelClient.CancelBooking(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Cancel Bookinig Result = %s\n", resultStr)
}

//Cancel Booking Confirm
func TestCancelBookingConfirm(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	//only need booking number
	jsonBookingCancelConfirmRequest := `{
		"Ref":"01999731",
		"InternalRef":"10327"
	}`

	var bookingCancelRequest hbecommon.BookingCancelConfirmationRequest
	err := json.Unmarshal([]byte(jsonBookingCancelConfirmRequest), &bookingCancelRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	cancelBookingResponse := restelClient.CancelBookingConfirm(&bookingCancelRequest)

	resultStr, err := json.Marshal(cancelBookingResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Cancel Bookinig Confirm Result = %s\n", resultStr)
}

//Get Booking Info
func TestGetBookingInfo(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	//only need booking number
	jsonBookingInfoRequest := `{
		"HotelId":123456
	}`

	var bookingInfoRequest hbecommon.BookingInfoRequest
	err := json.Unmarshal([]byte(jsonBookingInfoRequest), &bookingInfoRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	bookingInfoResponse := restelClient.GetBookingInfo(&bookingInfoRequest)

	resultStr, err := json.Marshal(bookingInfoResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("HBE Get Bookinig Info Result = %s\n", resultStr)
}

//Get Room Cxl
func TestGetRoomCxl(t *testing.T) {
	restelClient := restel.NewHotelBookingProvider(&hotelProviderSetting)
	restelClient.Init()

	//only need booking number
	jsonRoomCxlRequest := `{
		"HotelId":"626886",
		"RoomRef":"{\"RoomType\":\"9X\",\"RoomDesc\":\"Twin con desayuno\",\"RoomLine\":{\"MealPlanType\":\"BB\",\"PlanPrice\":\"287.97\",\"Currency\":\"DO\",\"Status\":\"OK\",\"MinPrice\":\"\",\"FeeRate\":\"0\",\"CompressedAvailability\":[\"9X#1#VR#151.56#0#BB#OK#20180920#20180921#DO#2-0#0#0#201808301113#626886#\",\"9X#1#VR#136.41#0#BB#OK#20180921#20180922#DO#2-0#0#0#201808301113#626886#\"]}}"
	}`

	var roomCxlRequest hbecommon.RoomCxlRequest
	err := json.Unmarshal([]byte(jsonRoomCxlRequest), &roomCxlRequest)
	if err != nil {
		t.Errorf("Cannot parse request json : %s", err)
		return
	}

	roomCxlResponse := restelClient.GetRoomCxl(&roomCxlRequest)

	resultStr, err := json.Marshal(roomCxlResponse)
	if err != nil {
		t.Errorf("Wrong response json : %s", err)
		return
	}
	fmt.Printf("\nHBE Get Room Cxl Result = %s\n", resultStr)
}
