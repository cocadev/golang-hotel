package gta

import (
	"fmt"

	hbecommon "../roomres/hbe/common"

	//repository "roomres/repository"
	"strings"
	"time"

	roomresutils "../roomres/utils"

	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_bookingreportrequest(reportRequestHBE *hbecommon.BookingReportRequest, reportRequest *SearchBookingRequest) {

	reportRequest.BookingDateRange = &BookingDateRange{
		DateType: "creation",
		FromDate: reportRequestHBE.BookingDateFrom.Format(roomresutils.LayoutYYYYMMDD),
		ToDate:   reportRequestHBE.BookingDateTo.Format(roomresutils.LayoutYYYYMMDD),
	}
}

func mapping_bookingreportresponse_to_hbe(searchBooking *SearchBooking, searchBookingItem *SearchBookingItemServiceResponse, reportResponseHBE *hbecommon.BookingReportResponse) {

	hbeBooking := &hbecommon.BookingReport{}

	hbeBooking.Ref = func() string {
		if ref := searchBooking.GetReference("api"); ref != nil {
			return ref.Value
		}
		return ""
	}()

	hbeBooking.InternalRef = func() string {
		if ref := searchBooking.GetReference("client"); ref != nil {
			return ref.Value
		}
		return ""
	}()

	if searchBookingItem.SearchBookingItemResponse == nil {
		return
	}

	if item := From(searchBookingItem.SearchBookingItemResponse.BookingItems).Where(func(item interface{}) bool {
		return strings.ToUpper(item.(*BookingItem).ItemType) == "HOTEL"
	}).First(); item != nil {

		hotelBookingItem := item.(*BookingItem)

		hbeBooking.CheckIn, _ = time.Parse(roomresutils.LayoutYYYYMMDD, hotelBookingItem.CheckIn)
		hbeBooking.CheckOut, _ = time.Parse(roomresutils.LayoutYYYYMMDD, hotelBookingItem.CheckOut)

		hbeBooking.HotelName = hotelBookingItem.HotelName
		hbeBooking.CityName = hotelBookingItem.CityName

		hbeBooking.Total = searchBookingItem.SearchBookingItemResponse.BookingPrice.Gross

		if searchBookingItem.SearchBookingItemResponse.BookingStatus.IsConfirmed() {
			hbeBooking.Status = hbecommon.BookingReportStatusConfirmed
		} else {
			hbeBooking.Status = hbecommon.BookingReportStatusNotConfirmed
		}

		reportResponseHBE.Bookings = append(reportResponseHBE.Bookings, hbeBooking)
	}
}

// func mapping_hbe_to_searchiteminformationrequest(request *hbecommon.BookingInfoRequest, hotelMappings []*repository.HotelMapping, searchItemInformationServiceRequest *SearchItemInformationServiceRequest) {

// 	hotelMapping := hotelMappings[0]

// 	hotelCodes := strings.Split(hotelMapping.SupplierId, ";")

// 	searchItemInformationServiceRequest.SearchItemInformationRequest = &SearchItemInformationRequest{}

// 	searchItemInformationServiceRequest.SearchItemInformationRequest.ItemType = "hotel"
// 	searchItemInformationServiceRequest.SearchItemInformationRequest.ItemCode = hotelCodes[1]

// 	searchItemInformationServiceRequest.SearchItemInformationRequest.ItemDestination = &ItemDestination{DestinationType: "city", DestinationCode: hotelCodes[0]}
// }

func mapping_searchiteminformationrequest_to_hbe(searchItemInformationServiceResponse *SearchItemInformationServiceResponse, response *hbecommon.BookingInfoResponse) {

	if searchItemInformationServiceResponse == nil ||
		len(searchItemInformationServiceResponse.SearchItems) == 0 ||
		searchItemInformationServiceResponse.SearchItems[0].HotelInformation == nil {
		return
	}

	response.HotelPhoneNumber = searchItemInformationServiceResponse.SearchItems[0].HotelInformation.Telephone
}

func mapping_searchresponse_to_hbe(
	searchResponse *SearchResponse,
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	hotelSearchResponse *hbecommon.HotelSearchResponse,
) {

	hotelSearchResponse.Hotels = []*hbecommon.Hotel{}

	for _, hotel := range searchResponse.Hotels {

		if ok, roomCombinations := GetRooms(hotelSearchRequest, hotel); ok {

			hbeHotel := &hbecommon.Hotel{}

			mapping_hotel_to_hbe(roomCombinations, hotel, hbeHotel, hotelSearchRequest.GetLos())

			hotelSearchResponse.Hotels = append(hotelSearchResponse.Hotels, hbeHotel)
		}

	}
}

type RoomCombination struct {
	HotelPaxRoom      *HotelPaxRoom
	HotelRoomCategory *HotelRoomCategory
}

func GetRooms(hotelSearchRequest *hbecommon.HotelSearchRequest, hotel *Hotel) (bool, []*RoomCombination) {

	if len(hotel.HotelPaxRooms) < len(hotelSearchRequest.RequestedRooms) {
		return false, nil
	}

	roomCombinations := []*RoomCombination{}

	for i := 1; i <= len(hotelSearchRequest.RequestedRooms); i++ {

		roomCheck := false
		for _, hotelPaxRoom := range hotel.HotelPaxRooms {
			if hotelPaxRoom.RoomIndex == i && len(hotelPaxRoom.HotelRoomCategories) > 0 {

				var hotelRoomCategories []*HotelRoomCategory

				From(hotelPaxRoom.HotelRoomCategories).OrderBy(func(category interface{}) interface{} {
					return category.(*HotelRoomCategory).RoomCategoryPrice.Price
				}).ToSlice(&hotelRoomCategories)

				roomCombinations = append(roomCombinations, &RoomCombination{
					HotelPaxRoom:      hotelPaxRoom,
					HotelRoomCategory: hotelRoomCategories[0]})

				roomCheck = true
				break
			}
		}

		if !roomCheck {
			return false, nil
		}
	}

	return true, roomCombinations
}

func mapping_hotel_to_hbe(
	roomCombinations []*RoomCombination,
	hotel *Hotel,
	hbeHotel *hbecommon.Hotel,
	los int,

) {

	hbeHotel.HotelId = fmt.Sprintf("%[1]s;%[2]s", hotel.HotelCity.Code, hotel.HotelItemCode.Code)

	hbeHotel.CheapestRoom = &hbecommon.RoomType{}

	mapping_roomtype_to_hbe(
		roomCombinations[0].HotelPaxRoom,
		hbeHotel.CheapestRoom, los)

	total := From(roomCombinations).Select(func(roomCombination interface{}) interface{} {
		return roomCombination.(*RoomCombination).HotelRoomCategory.RoomCategoryPrice.Price
	}).SumFloats()

	total = total / float64(los) / float64(len(roomCombinations))

	hbeHotel.CheapestRoom.Rate.PerNight = float32(total)
	hbeHotel.CheapestRoom.Rate.PerNightBase = float32(total)
}

func mapping_roomtype_to_hbe(
	paxRoom *HotelPaxRoom,
	roomType *hbecommon.RoomType,
	los int,
) {

	roomType.Taxes = []*hbecommon.Tax{}
	roomType.Surcharges = []*hbecommon.Surcharge{}

	roomCategory := paxRoom.HotelRoomCategories[0]

	roomType.ShortRef = roomCategory.Id
	roomType.Ref = roomCategory.Id

	roomType.Description = roomCategory.Description

	roomType.Rate = &hbecommon.Rate{}
	roomType.Rate.PerNight = roomCategory.RoomCategoryPrice.Price / float32(los)
	roomType.Rate.PerNightBase = roomCategory.RoomCategoryPrice.Price / float32(los)
}

func mapping_hbe_to_SearchHotelPricePaxRequests(hotelsPerBatch int, hotelSearchRequest *hbecommon.HotelSearchRequest, requests *[]*SearchHotelPricePaxRequest) {

	destinationGroups := SplitDestinationGroups(hotelsPerBatch, CreateDestinationGroups(hotelSearchRequest.HotelIds))

	for _, destinationGroup := range destinationGroups {

		request := &SearchHotelPricePaxRequest{}

		mapping_hbe_to_SearchHotelPricePaxRequest(hotelSearchRequest, destinationGroup, request)

		*requests = append(*requests, request)
	}
}

func mapping_hbe_to_SearchHotelPricePaxRequest(hotelSearchRequest *hbecommon.HotelSearchRequest, destinationGroup *DestinationGroup, request *SearchHotelPricePaxRequest) {

	request.CheckInDate = hotelSearchRequest.CheckIn
	request.CheckOutDate = hotelSearchRequest.CheckOut
	request.IncludePriceBreakdown = &IncludePriceBreakdown{}
	request.IncludeChargeConditions = &IncludeChargeConditions{DateFormatResponse: true}
	request.ImmediateConfirmationOnly = &ImmediateConfirmationOnly{}

	/*
		maxAdults := From(hotelSearchRequest.RequestedRooms).Select(func(requestedRoom interface{}) interface{} {
			return requestedRoom.(hbecommon.RoomRequest).Adults
		}).Max().(int)
	*/

	maxRequestedRoom := hotelSearchRequest.RequestedRooms[0]
	for i := 1; i < len(hotelSearchRequest.RequestedRooms); i++ {

		requestedRoom := hotelSearchRequest.RequestedRooms[i]

		if requestedRoom.Adults > maxRequestedRoom.Adults ||
			requestedRoom.Cots+len(requestedRoom.ChildAges) > maxRequestedRoom.Cots+len(maxRequestedRoom.ChildAges) {

			maxRequestedRoom = requestedRoom
		}
	}

	request.PaxRooms = []*PaxRoom{}
	for i, hbeRequestedRoom := range hotelSearchRequest.RequestedRooms {

		paxRoom := &PaxRoom{}

		mapping_hbe_RequestedRoom(hbeRequestedRoom, paxRoom, i+1, maxRequestedRoom)

		request.PaxRooms = append(request.PaxRooms, paxRoom)
	}

	request.ItemDestination = &ItemDestination{DestinationType: "city", DestinationCode: destinationGroup.DestinationCode}
	request.ItemCodes = []*ItemCode{}

	for _, hotelId := range destinationGroup.HotelIds {

		itemCode := &ItemCode{ItemCodeText: hotelId}

		request.ItemCodes = append(request.ItemCodes, itemCode)
	}
}

func mapping_hbe_RequestedRoom(hbeRequestedRoom *hbecommon.RoomRequest, requestedPaxRoom *PaxRoom, index int, maxRequestedRoom *hbecommon.RoomRequest) {

	requestedPaxRoom.RoomIndex = index
	requestedPaxRoom.Adults = maxRequestedRoom.Adults //hbeRequestedRoom.Adults
	requestedPaxRoom.Cots = hbeRequestedRoom.Cots

	for _, childAge := range maxRequestedRoom.ChildAges /*hbeRequestedRoom.ChildAges*/ {
		requestedPaxRoom.ChildAges = append(requestedPaxRoom.ChildAges, &ChildAge{Age: childAge.Age})
	}
}
