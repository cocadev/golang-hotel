package hb

import (
	"fmt"
	"strings"

	hbecommon "../roomres/hbe/common"
	. "github.com/ahmetb/go-linq"
	//repository "../roomres/repository"
)

func mapping_searchresponsedetail_to_hbe(
	availabilityResponse *AvailabilityResponse,
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse) {

	if availabilityResponse.Hotels == nil {
		return
	}

	hbeHotelSearchResponse.Hotels = []*hbecommon.Hotel{}
	for _, hotel := range availabilityResponse.Hotels.Hotels {

		hbeHotel := &hbecommon.Hotel{HotelId: fmt.Sprintf("%s-%d", hotel.DestinationCode, hotel.HotelId)}

		mapping_hoteldetail_to_hbe(hotelSearchRequest, hotel, hbeHotel)

		if hbeHotel.CheapestRoom != nil {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
		}
	}
}

func mapping_hoteldetail_to_hbe(hotelSearchRequest *hbecommon.HotelSearchRequest, hotel *Hotel, hbeHotel *hbecommon.Hotel) {

	roomRates := FindBestRoomDetailRateCombinations(hotelSearchRequest.RequestedRooms, hotel.Rooms)
	hbeRoomTypes := []*hbecommon.RoomType{}
	specialRoomRefs := hotelSearchRequest.SpecificRoomRefs

	for _, roomRate := range roomRates {
		hbeRoomType := &hbecommon.RoomType{}
		//check max/min range
		if hotelSearchRequest.MaxPrice > 0 && (hotelSearchRequest.MaxPrice < roomRate.Rate.NetValue) {
			continue
		}
		if hotelSearchRequest.MinPrice > 0 && (hotelSearchRequest.MinPrice > roomRate.Rate.NetValue) {
			continue
		}
		if specialRoomRefs != nil && len(specialRoomRefs) > 0 {
			if !checkSpecialRoomRefs(specialRoomRefs, roomRate) {
				continue
			}
		}

		mapping_roomtypedetail_to_hbe(
			hotelSearchRequest.GetLos(),
			hotel,
			roomRate.Room,
			roomRate.Rate,
			hbeRoomType,
		)

		hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)
	}

	if len(hbeRoomTypes) > 0 {
		//sort room by price

		From(hbeRoomTypes).OrderBy(func(roomType interface{}) interface{} {
			return roomType.(*hbecommon.RoomType).Rate.PerNight
		}).ToSlice(&hbeHotel.RoomTypes)

		cheapestRoomType := hbeHotel.RoomTypes[0]
		hbeHotel.CheapestRoom = cheapestRoomType
	}
}

func mapping_roomtypedetail_to_hbe(los int, hotel *Hotel, room *Room, rate *Rate, hbeRoomType *hbecommon.RoomType) {

	hbeRoomType.ShortRef = room.Code
	hbeRoomType.Ref = rate.RateKey

	hbeRoomType.Description = room.Name

	hbeRoomType.Taxes = []*hbecommon.Tax{}
	hbeRoomType.Surcharges = []*hbecommon.Surcharge{}

	hbeRoomType.Rate = &hbecommon.Rate{
		PerNight:     rate.NetValue / float32(los),
		PerNightBase: rate.NetValue / float32(los),
		Total:        rate.NetValue,
		CurrecyCode:  hotel.CurrencyCode,
	}

	if rate.CancellationPolicies != nil && len(rate.CancellationPolicies) > 0 {
		hbeRoomType.CancellationPolicy = GenerateCancellationPolicyShort(rate.CancellationPolicies)
		hbeRoomType.FreeCancellationPolicy = GetFreeCancelDate(rate.CancellationPolicies)
	} else {
		hbeRoomType.NonRefundable = true
	}

	//added by Li, 20180917
	//Set taxes
	if rate.Taxes != nil && len(rate.Taxes.Taxes) > 0 {
		for _, tax := range rate.Taxes.Taxes {
			hbeRoomType.Taxes = append(hbeRoomType.Taxes, &hbecommon.Tax{
				Type:   tax.Type,
				Amount: tax.Amount,
			})
		}
	}
}

func FindBestRoomDetailRateCombinations(roomRequests []*hbecommon.RoomRequest, rooms []*Room) []*RoomRate {

	roomRates := map[string]*RoomRate{}
	rateCount := 0

	for _, room := range rooms {

		for _, rate := range room.Rates {
			rateCount++
			if strings.ToUpper(rate.RateType) != "BOOKABLE" && strings.ToUpper(rate.RateType) != "RECHECK" {
				continue
			}

			if strings.ToUpper(rate.RateType) == "RECHECK" && len(roomRequests) > 1 {
				continue
			}

			rate.NetValue = rate.GetNet()

			if roomRate, ok := roomRates[rate.RateKey]; !ok {
				roomRates[rate.RateKey] = &RoomRate{Room: room, Rate: rate, OccupanciesFits: 0}
			} else {

				if roomRate.Rate.NetValue > rate.NetValue {
					roomRate.Room = room
					roomRate.Rate = rate
				}
			}
		}
	}

	for _, roomRequest := range roomRequests {
		for _, roomRate := range roomRates {
			if roomRate.Rate.Adults >= roomRequest.Adults && roomRate.Rate.Children >= roomRequest.Children {
				roomRate.OccupanciesFits++
			}
		}
	}

	validCombinations := []*RoomRate{}

	for _, roomRate := range roomRates {
		if roomRate.OccupanciesFits == len(roomRequests) {
			validCombinations = append(validCombinations, roomRate)
		}
	}

	return validCombinations
}

func checkSpecialRoomRefs(specialRoomRefs []string, roomRate *RoomRate) bool {
	for _, specialRoomRef := range specialRoomRefs {
		params := strings.Split(specialRoomRef, "|")
		roomRateParams := strings.Split(roomRate.Rate.RateKey, "|")

		for i := 0; i < 11; i++ {
			if len(params) > i && params[i] != "" && params[i] != roomRateParams[i] {
				return false
			}
		}
	}
	return true
}

func makeBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, bookingRequest *BookingRequest) {
	bookingRequest.ClientReference = "IntegrationAgency"
	bookingRequest.Holder = &Holder{
		FirstName: hbeBookingRequest.Customer.FirstName,
		LastName:  hbeBookingRequest.Customer.LastName,
	}
	bookingRooms := []*BookingRoom{}
	//hotelId := hbeBookingRequest.Hotel.HotelId
	for _, room := range hbeBookingRequest.Hotel.Rooms {
		bookingPaxes := []*BookingPax{}
		for _, guest := range room.Guests {
			gender := "AD"
			if !guest.IsAdult {
				gender = "CH"
			}
			bookingPaxes = append(bookingPaxes, &BookingPax{
				RoomID:    "1",
				Type:      gender,
				FirstName: guest.FirstName,
				LastName:  guest.LastName,
			})
		}

		bookingRooms = append(bookingRooms, &BookingRoom{
			RateKey: room.Ref,
			Paxes:   bookingPaxes,
		})
	}

	bookingRequest.BookingRooms = bookingRooms
}

func mapping_bookresponse_to_hbe(
	bookResponse *BookingResponse,
	hbeBookingRequest *hbecommon.BookingRequest,
	bookingResponse *hbecommon.BookingResponse) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	bookingResult := bookResponse.BookingResult

	if bookingResult == nil {
		return
	}

	if bookingResult.Status == "CONFIRMED" {
		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		bookingResponse.Booking = &hbecommon.Booking{
			Ref:               bookingResult.Reference,
			ItineraryId:       bookingResult.ClientReference,
			SupplierReference: bookingResult.Hotel.Supplier.VatNumber,
			SupplierName:      bookingResult.Hotel.Supplier.Name,
			Total:             bookingResult.TotalNet,
		}

		if bookingResult.Hotel.Rooms != nil {
			room := bookingResult.Hotel.Rooms[0]
			if room.Rates != nil && len(room.Rates) > 0 {
				bookingResponse.Booking.CancellationPolicy = GenerateCancellationPolicyShort(room.Rates[0].CancellationPolicies)
				bookingResponse.Booking.FreeCancellationPolicy = GetFreeCancelDate(room.Rates[0].CancellationPolicies)
			}
		}

	} else {
		bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
		bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "Failed to booking reservation"}}
	}
}

func mapping_cancelresponse_to_hbe(
	cancelResponse *BookingResponse,
	hbeBookingCancelRequest *hbecommon.BookingCancelRequest,
	bookingCancelResponse *hbecommon.BookingCancelResponse) {

	cancelResult := cancelResponse.BookingResult

	if cancelResult == nil {
		bookingCancelResponse.Status = "Failed"
		return
	}

	if cancelResult.Status == "CANCELLED" {
		bookingCancelResponse.Status = "200"
		bookingCancelResponse.Ref = hbeBookingCancelRequest.Ref
	} else {
		bookingCancelResponse.Status = "Failed"
	}
}

//added by Li, 20180917
func mapping_hbe_to_checkraterequest(
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	checkRateRequest *CheckRateRequest) {

	checkRateRequest.Language = "ENG"
	checkRateRequest.Upselling = "False"
	checkRateRequest.Rooms = []*SimpleRate{}

	for _, ratekey := range hotelSearchRequest.SpecificRoomRefs {
		checkRateRequest.Rooms = append(checkRateRequest.Rooms, &SimpleRate{
			RateKey: ratekey,
		})
	}
}
