package gta

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"

	//repository "roomres/repository"

	. "github.com/ahmetb/go-linq"
)

func mapping_searchresponsedetail_to_hbe(
	searchResponse *SearchResponse,
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	hotelSearchResponse *hbecommon.HotelSearchResponse,
) {
	hotelSearchResponse.Hotels = []*hbecommon.Hotel{}

	for _, hotel := range searchResponse.Hotels {
		roomCategories := []*HotelRoomCategory{}
		for _, paxRoom := range hotel.HotelPaxRooms {
			for _, roomCategory := range paxRoom.HotelRoomCategories {
				if len(hotelSearchRequest.SpecificRoomRefs) > 0 {
					if !checkElementInArray(hotelSearchRequest.SpecificRoomRefs, roomCategory.Id) {
						continue
					}
				}
				roomCategories = append(roomCategories, roomCategory)
			}
		}

		if len(roomCategories) > 0 {
			hbeHotel := &hbecommon.Hotel{}

			mapping_hoteldetail_to_hbe(roomCategories, hotel, hbeHotel, hotelSearchRequest.GetLos())

			hotelSearchResponse.Hotels = append(hotelSearchResponse.Hotels, hbeHotel)
		}
	}
}

func GetRoomsDetail(hotelSearchRequest *hbecommon.HotelSearchRequest, hotel *Hotel) (bool, []*RoomCombination) {

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

func mapping_hoteldetail_to_hbe(
	roomCategories []*HotelRoomCategory,
	hotel *Hotel,
	hbeHotel *hbecommon.Hotel,
	los int,
) {
	hbeHotel.HotelId = fmt.Sprintf("%[1]s;%[2]s", hotel.HotelCity.Code, hotel.HotelItemCode.Code)

	hbeHotel.CheapestRoom = &hbecommon.RoomType{}
	roomTypes := []*hbecommon.RoomType{}
	mapping_roomtypedetail_to_hbe(
		roomCategories,
		&roomTypes, los)

	//fmt.Printf("Len of room types = %d\n", len(roomTypes))
	From(roomTypes).OrderBy(func(roomType interface{}) interface{} {
		return roomType.(*hbecommon.RoomType).Rate.PerNight
	}).ToSlice(&roomTypes)

	if len(roomTypes) > 0 {
		hbeHotel.CheapestRoom = roomTypes[0]
		hbeHotel.RoomTypes = roomTypes
	}
}

func mapping_roomtypedetail_to_hbe(
	roomCategories []*HotelRoomCategory,
	roomTypes *[]*hbecommon.RoomType,
	los int,
) {

	for idx, roomCategory := range roomCategories {
		//fmt.Printf("roomCategory = %+v\n", roomCategory)
		cancelPolicy := GenerateCancellationPolicyShort(roomCategory.ChargeConditions.ChargeConditions)
		roomType := hbecommon.RoomType{
			Taxes:       []*hbecommon.Tax{},
			Surcharges:  []*hbecommon.Surcharge{},
			ShortRef:    ConvertIntToString(idx + 1),
			Ref:         roomCategory.Id,
			Description: roomCategory.Description,
			Rate: &hbecommon.Rate{
				PerNight:     roomCategory.RoomCategoryPrice.Price / float32(los),
				PerNightBase: roomCategory.RoomCategoryPrice.Price / float32(los),
				Total:        roomCategory.RoomCategoryPrice.Price,
			},
			CancellationPolicy:     cancelPolicy,
			FreeCancellationPolicy: GetFreeCancelDate(roomCategory.ChargeConditions.ChargeConditions),
		}

		if cancelPolicy == "" {
			roomType.NonRefundable = true
		}

		*roomTypes = append(*roomTypes, &roomType)
	}

}

func makeBookingContent(hbeBookingRequest *hbecommon.BookingRequest, currencyCode string) *AddBookingRequest {
	compoundHotelId := hbeBookingRequest.Hotel.HotelId
	hotelCodes := strings.Split(compoundHotelId, ";")
	BookingItems := &BookingItems{
		BookingItems: []*BookingRequestItem{},
	}
	PaxNames := &PaxNames{}
	BookingRequestItem := &BookingRequestItem{
		ItemType:      "hotel",
		ExpectedPrice: hbeBookingRequest.Total,
		ItemReference: 1,
		ItemCity: &ItemCity{
			Code: hotelCodes[0],
		},
		ItemCode: &ItemCode2{
			Code: hotelCodes[1],
		},
		CheckInDate:  hbeBookingRequest.CheckIn,
		CheckOutDate: hbeBookingRequest.CheckOut,
	}

	hotelRooms := []*HotelRoom{}
	idx := 0
	for _, room := range hbeBookingRequest.Hotel.Rooms {
		PaxIds := []int{}

		for _, guest := range room.Guests {
			idx++
			PaxNames.PaxNames = append(PaxNames.PaxNames, &PaxName{
				PaxId:   idx,
				PaxName: guest.FirstName + " " + guest.LastName,
			})
			PaxIds = append(PaxIds, idx)
		}
		HotelRoom := &HotelRoom{
			Id:     room.Ref,
			Code:   getHotelRoomCode(len(PaxIds)),
			PaxIds: PaxIds,
		}
		hotelRooms = append(hotelRooms, HotelRoom)
	}
	BookingRequestItem.HotelRooms = &HotelRooms{
		HotelRooms: hotelRooms,
	}
	BookingItems.BookingItems = append(BookingItems.BookingItems, BookingRequestItem)
	bookingContent := &AddBookingRequest{
		Currency:             currencyCode,
		BookingReference:     makeBookingRefernce(),
		BookingDepartureDate: hbeBookingRequest.CheckIn,
		PaxNames:             PaxNames,
		BookingItems:         BookingItems,
	}

	return bookingContent
}

func mapping_bookresponse_to_hbe(
	bookResponse *BookingResponse,
	hbeBookingRequest *hbecommon.BookingRequest,
	bookingResponse *hbecommon.BookingResponse) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	bookingResult := bookResponse.BookingResponseData

	if bookingResult == nil {
		return
	}

	if strings.TrimSpace(bookingResult.BookingStatus.Code) == "C" {
		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		referenceId := ""

		for _, ref := range bookingResult.BookingReferences.BookingReferences {
			if ref.ReferenceSource == "client" {
				referenceId = ref.Value
			}
		}
		bookItem := bookingResult.BookingItems.BookingResponseItems[0]
		bookingResponse.Booking = &hbecommon.Booking{
			Ref:                    referenceId,
			ItineraryId:            bookItem.ItemConfirmationReference,
			Total:                  bookItem.ItemPrice.Nett,
			CancellationPolicy:     GenerateCancellationPolicyShort(bookItem.ChargeConditions),
			FreeCancellationPolicy: GetFreeCancelDate(bookItem.ChargeConditions),
		}
	} else {
		bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
		bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "Failed to booking reservation"}}
	}
}

func mapping_cancelresponse_to_hbe(
	cancelResponse *CancelResponse,
	hbeBookingCancelRequest *hbecommon.BookingCancelRequest,
	bookingCancelResponse *hbecommon.BookingCancelResponse) {

	if cancelResponse.BookingResponseData != nil && strings.TrimSpace(cancelResponse.BookingResponseData.BookingStatus.Code) == "X" {
		referenceId := ""
		cancelResult := cancelResponse.BookingResponseData
		for _, ref := range cancelResult.BookingReferences.BookingReferences {
			if ref.ReferenceSource == "client" {
				referenceId = ref.Value
			}
		}

		bookingCancelResponse.Status = "200"
		bookingCancelResponse.Ref = referenceId
	} else {
		bookingCancelResponse.Status = "Failed"
	}
}

func makeBookingRefernce() string {
	//Need to create unique ID...
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("B%s-%x", time.Now().Format("2006010215"), b[0:4])
}

func getHotelRoomCode(guests int) string {
	if guests == 1 {
		return "SB"
	} else if guests == 2 {
		return "DB"
	} else if guests == 3 {
		return "TR"
	} else if guests == 4 {
		return "Q"
	} else {
		return "TB"
	}
}
