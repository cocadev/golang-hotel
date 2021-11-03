package webbed

import (
	"encoding/json"
	"fmt"
	"strings"

	hbecommon "../roomres/hbe/common"

	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_searchrequest(hotelSearchRequest *hbecommon.HotelSearchRequest, searchRequest *SearchRequest) {

	//	searchRequest.SearchDetails = &SearchDetails{}
	searchRequest.CheckInDate = hotelSearchRequest.CheckIn
	searchRequest.CheckOutDate = hotelSearchRequest.CheckOut
	searchRequest.Language = "en"
	searchRequest.NumberOfRooms = ConvertIntToString(len(hotelSearchRequest.RequestedRooms))

	searchRequest.HotelIDs = hotelSearchRequest.HotelIds
	for idx, id := range searchRequest.HotelIDs {
		searchRequest.HotelIDs[idx] = GetRawHotelId(id)
	}

	searchRequest.B2c = "0"

	if hotelSearchRequest.MinPrice > 0 {
		searchRequest.MinPrice = ConvertFloat32ToString(hotelSearchRequest.MinPrice)
	}
	if hotelSearchRequest.MaxPrice > 0 {
		searchRequest.MaxPrice = ConvertFloat32ToString(hotelSearchRequest.MaxPrice)
	}

	ages := []string{}
	NumberOfAdults := 0
	infants := 0
	for _, requestRoom := range hotelSearchRequest.RequestedRooms {
		NumberOfAdults += requestRoom.Adults

		if len(requestRoom.ChildAges) > 0 {
			for _, childAge := range requestRoom.ChildAges {
				if childAge.Age < 2 {
					infants++
				} else {
					ages = append(ages, ConvertIntToString(childAge.Age))
				}
			}
		}
	}
	if len(ages) > 0 {
		searchRequest.ChildrenAges = strings.Join(ages, ",")
	}
	searchRequest.NumberOfAdults = ConvertIntToString(NumberOfAdults)
	searchRequest.NumberOfChildren = ConvertIntToString(len(ages))
	searchRequest.Infant = ConvertIntToString(infants)
}

func mapping_searchresponse_to_hbe(
	searchResponse *SearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse) {

	los := hbeHotelSearchRequest.GetLos()

	//fmt.Println("searchResponse.Hotels = %d\n", len(searchResponse.Hotels))

	for _, hotelItem := range searchResponse.Hotels {

		hbeHotel := &hbecommon.Hotel{}

		mapping_hotel_to_hbe(
			hbeHotelSearchRequest,
			hotelItem,
			hbeHotel,
			los,
		)

		if hbeHotel.CheapestRoom != nil {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
		}
	}
}

func mapping_hotel_to_hbe(
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hotelItem *Hotel,
	hbeHotel *hbecommon.Hotel,
	los int) {

	//fmt.Printf("hotelItem=%+v\n", hotelItem)

	hbeHotel.HotelId = "wbb_" + hotelItem.Id
	hbeHotel.CustomTag = hotelItem.Type
	hbeHotel.Notes = hotelItem.Description

	hbeRoomTypes := []*hbecommon.RoomType{}

	//fmt.Printf("hotelItem.RoomTypes = %d\n", len(hotelItem.RoomTypes))
	for _, RoomType := range hotelItem.RoomTypes {
		//fmt.Printf("Rooms = %d\n", len(RoomType.Rooms))
		for _, Room := range RoomType.Rooms {
			//create room type per room
			hbeRoomType := &hbecommon.RoomType{}
			hbeRoomType.Description = RoomType.RoomType

			var Total float32 = 0
			CurrencyCode := "USD"
			for _, Meal := range Room.Meals {
				if len(Meal.Prices) > 0 {
					Total += Meal.Prices[0].Price
					CurrencyCode = Meal.Prices[0].Currency
				}
				if len(Meal.Discounts) > 0 {
					Total -= Meal.Discounts[0].Amount
				}
			}
			hbeRoomType.Rate = &hbecommon.Rate{}
			hbeRoomType.Rate.PerNight = Total / float32(los)
			hbeRoomType.Rate.PerNightBase = Total / float32(los)
			hbeRoomType.Rate.CurrecyCode = CurrencyCode
			hbeRoomType.Rate.NightlyRateTotal = Total
			hbeRoomType.Rate.Total = Total

			hbeRoomType.Ref = Room.Id

			//check special refs
			if hbeHotelSearchRequest.SpecificRoomRefs != nil && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
				var isMatched = false
				for _, roomRefStr := range hbeHotelSearchRequest.SpecificRoomRefs {
					if roomRefStr == Room.Id {
						isMatched = true
						break
					}
				}

				if !isMatched {
					continue
				}
			}

			if hbeHotelSearchRequest.Details == true {
				hbeRoomType.Nights = los

				if Room.CancellationPolicies != nil && len(Room.CancellationPolicies) > 0 {
					hbeRoomType.CancellationPolicy = GenerateCancellationPolicyShort(Room.CancellationPolicies, hbeHotelSearchRequest.CheckIn, Total)
					hbeRoomType.FreeCancellationPolicy = GetFreeCancelDate(Room.CancellationPolicies, hbeHotelSearchRequest.CheckIn)

					if Room.CancellationPolicies[0].Deadline == "" && Room.CancellationPolicies[0].Percentage == "100" {
						hbeRoomType.NonRefundable = true
					} else {
						hbeRoomType.NonRefundable = false
					}
				}
			}
			if hbeHotelSearchRequest.Details == true {
				hbeRoomType.NormalBeddingOccupancy = Room.Beds
				hbeRoomType.ExtraBeddingOccupancy = Room.Extrabeds
			}

			hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)
		}
	}

	//fmt.Printf("hbeRoomTypes = %+v\n", hbeRoomTypes)

	From(hbeRoomTypes).OrderBy(func(roomType interface{}) interface{} {
		return roomType.(*hbecommon.RoomType).Rate.PerNight
	}).ToSlice(&hbeHotel.RoomTypes)

	if len(hbeHotel.RoomTypes) > 0 {
		hbeHotel.CheapestRoom = hbeHotel.RoomTypes[0]
	}

	if !hbeHotelSearchRequest.Details {
		hbeHotel.RoomTypes = []*hbecommon.RoomType{}

		// if hbeHotel.CheapestRoom != nil {
		// 	hbeHotel.CheapestRoom.Ref = ""
		// 	hbeHotel.CheapestRoom.ShortRef = ""
		// }
	}
}

func makePreBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, preBookingRequest *PreBookingRequest) {

	preBookingRequest.CheckInDate = hbeBookingRequest.CheckIn
	preBookingRequest.CheckOutDate = hbeBookingRequest.CheckOut

	room := hbeBookingRequest.Hotel.Rooms[0]
	preBookingRequest.RoomId = room.Ref
	preBookingRequest.Rooms = room.Count
	preBookingRequest.Adults = room.Adults
	preBookingRequest.Children = room.Children
	preBookingRequest.B2c = 0
	preBookingRequest.Infant = 0
	preBookingRequest.Guests = room.Guests

	var specialRequest SpecialRequest
	if room.SpecialRequest != "" {
		err := json.Unmarshal([]byte(room.SpecialRequest), &specialRequest)
		if err != nil {
			fmt.Printf("Request parsing error")
			return
		}
		preBookingRequest.MealId = specialRequest.MealId
	}

	if room.Children > 0 {
		ages := []string{}
		for _, guest := range room.Guests {
			if !guest.IsAdult {
				if guest.Age < 2 {
					preBookingRequest.Infant++
				} else {
					ages = append(ages, ConvertIntToString(guest.Age))
				}
			}
		}

		preBookingRequest.ChildrenAges = strings.Join(ages, ",")
		preBookingRequest.Children = len(ages)
	}
}

func makeBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, preBookResponse *PreBookResponse) *BookingRequest {
	if preBookResponse.PreBookCode != "" {
		bookingRequest := &BookingRequest{
			CheckInDate:     hbeBookingRequest.CheckIn,
			CheckOutDate:    hbeBookingRequest.CheckOut,
			RoomId:          preBookResponse.RoomId,
			Rooms:           preBookResponse.Rooms,
			Adults:          preBookResponse.Adults,
			Children:        preBookResponse.Children,
			B2c:             0,
			Infant:          0,
			Email:           hbeBookingRequest.Customer.Email,
			PreBookCode:     preBookResponse.PreBookCode,
			PaymentMethodId: "1",
			YourRef:         "booking",
			MealId:          preBookResponse.MealId,
		}

		adultIdx := 0
		childIdx := 0
		for _, guest := range preBookResponse.Guests {
			if guest.IsAdult {
				adultIdx++
				SetField(bookingRequest, "AdultGuest"+ConvertIntToString(adultIdx)+"FirstName", guest.FirstName)
				SetField(bookingRequest, "AdultGuest"+ConvertIntToString(adultIdx)+"LastName", guest.LastName)
			} else {
				if guest.Age > 1 {
					childIdx++
					SetField(bookingRequest, "ChildrenGuest"+ConvertIntToString(childIdx)+"FirstName", guest.FirstName)
					SetField(bookingRequest, "ChildrenGuest"+ConvertIntToString(childIdx)+"LastName", guest.LastName)
					SetField(bookingRequest, "ChildrenGuestAge"+ConvertIntToString(childIdx), ConvertIntToString(guest.Age))
				} else {
					bookingRequest.Infant++
				}
			}
		}

		return bookingRequest
	}

	return nil
}

func mapping_bookresponse_to_hbe(
	bookResponse *BookResponse,
	hbeBookingRequest *hbecommon.BookingRequest,
	bookingResponse *hbecommon.BookingResponse) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	bookingResult := bookResponse

	if bookingResult == nil {
		return
	}

	if bookingResult.Booking != nil && bookingResult.Booking.BookingNumber != "" {
		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		var Total float32 = 0
		for _, Price := range bookingResult.Booking.Prices {
			if Price.Currency == bookingResult.Booking.Currency {
				Total = Price.Price
			}
		}

		bookingResponse.Booking = &hbecommon.Booking{
			Ref:   bookingResult.Booking.BookingNumber,
			Total: Total,
		}

		if bookingResult.Booking.CancellationPolicies != nil && len(bookingResult.Booking.CancellationPolicies) > 0 {
			bookingResponse.Booking.CancellationPolicy = GenerateCancellationPolicyShort(bookingResult.Booking.CancellationPolicies, hbeBookingRequest.CheckIn, Total)
			bookingResponse.Booking.FreeCancellationPolicy = GetFreeCancelDate(bookingResult.Booking.CancellationPolicies, hbeBookingRequest.CheckIn)
		}

	} else {
		bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
		bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "Failed to booking reservation"}}
	}
}

func mapping_cancelresponse_to_hbe(
	cancelResult *BookingCancelResponse,
	hbeBookingCancelRequest *hbecommon.BookingCancelRequest,
	bookingCancelResponse *hbecommon.BookingCancelResponse) {

	if cancelResult == nil {
		bookingCancelResponse.Status = "Failed"
		return
	}

	if cancelResult.Code == 1 {
		bookingCancelResponse.Status = "200"
		if len(cancelResult.Cancellationfees) > 0 {
			bookingCancelResponse.Payment = &hbecommon.Payment{
				Currency:        cancelResult.Cancellationfees[0].Currency,
				AmountInclusive: cancelResult.Cancellationfees[0].Price,
			}
		}
	} else {
		bookingCancelResponse.Status = "Failed"
	}
}

func SplitBatchSearchRequests(searchRequest *SearchRequest, maxBatchSize int) []*SearchRequest {

	batches := []*SearchRequest{}

	batches = append(batches, SplitBatchSearchRequest(searchRequest, maxBatchSize)...)

	return batches
}

func SplitBatchSearchRequest(searchRequest *SearchRequest, maxBatchSize int) []*SearchRequest {

	batches := []*SearchRequest{}

	var counter int = 0
	for counter < len(searchRequest.HotelIDs) {

		batch := searchRequest.Clone()

		var batchSize int = maxBatchSize
		if counter+batchSize > len(searchRequest.HotelIDs) {

			batchSize = len(searchRequest.HotelIDs) - counter
		}

		batch.HotelIDs = searchRequest.HotelIDs[counter : counter+batchSize]

		batches = append(batches, batch)

		counter += batchSize
	}

	return batches
}
