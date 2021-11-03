package derby

import (
	"encoding/json"

	hbecommon "../roomres/hbe/common"
	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_searchrequest(hotelSearchRequest *hbecommon.HotelSearchRequest) []*SearchRequest {
	requests := []*SearchRequest{}

	for _, hotelId := range hotelSearchRequest.HotelIds {
		searchRequest := &SearchRequest{}
		searchRequest.HotelId = hotelId
		searchRequest.StayRange = &StayRange{
			CheckIn:  hotelSearchRequest.CheckIn,
			CheckOut: hotelSearchRequest.CheckOut,
		}

		searchRequest.Iata = "string"

		ages := []int{}
		NumberOfRooms := len(hotelSearchRequest.RequestedRooms)
		NumberOfAdults := 0
		NumberOfChildren := 0
		for _, requestRoom := range hotelSearchRequest.RequestedRooms {
			NumberOfAdults += requestRoom.Adults
			NumberOfChildren += requestRoom.Children

			if len(requestRoom.ChildAges) > 0 {
				for _, childAge := range requestRoom.ChildAges {
					ages = append(ages, childAge.Age)
				}
			}
		}

		searchRequest.RoomCriteria = &RoomCriteria{
			RoomCount:  NumberOfRooms,
			AdultCount: NumberOfAdults,
			ChildCount: NumberOfChildren,
			ChildAges:  ages,
		}

		requests = append(requests, searchRequest)
	}

	return requests
}

func mapping_searchresponse_to_hbe(
	searchResponses []*SearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse) {

	los := hbeHotelSearchRequest.GetLos()

	//fmt.Println("searchResponse.Hotels = %d\n", len(searchResponse.Hotels))

	for _, hotelItem := range searchResponses {

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
	hotelItem *SearchResponse,
	hbeHotel *hbecommon.Hotel,
	los int) {

	hbeHotel.HotelId = hotelItem.HotelId
	hbeRoomTypes := []*hbecommon.RoomType{}

	//fmt.Printf("hotelItem.RoomTypes = %d\n", len(hotelItem.RoomTypes))
	for _, Room := range hotelItem.RoomRates {
		hbeRoomType := &hbecommon.RoomType{}

		var total float32 = 0.0
		if Room.AmountAfterTax != nil {
			for _, amt := range Room.AmountAfterTax {
				total += amt
			}
		} else {
			for _, amt := range Room.AmountBeforeTax {
				total += amt
			}
		}

		hbeRoomType.Rate = &hbecommon.Rate{
			PerNight:         total / float32(los),
			PerNightBase:     total / float32(los),
			CurrecyCode:      Room.Currency,
			NightlyRateTotal: total,
			Total:            total,
		}

		roomRef := &RoomRef{
			RoomId:          Room.RoomId,
			RateId:          Room.RateId,
			Currency:        Room.Currency,
			AmountBeforeTax: Room.AmountBeforeTax,
			AmountAfterTax:  Room.AmountAfterTax,
		}

		//check special refs
		if hbeHotelSearchRequest.Details == true && hbeHotelSearchRequest.SpecificRoomRefs != nil && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
			var isMatched = false

			for _, roomRefStr := range hbeHotelSearchRequest.SpecificRoomRefs {
				var specialRoomRef SpecialRoomRef
				err := json.Unmarshal([]byte(roomRefStr), &specialRoomRef)
				if err == nil {
					if specialRoomRef.RoomId != "" && specialRoomRef.RoomId != Room.RoomId {
						continue
					}
					if specialRoomRef.RateId != "" && specialRoomRef.RateId != Room.RateId {
						continue
					}
					if specialRoomRef.MealPlan != "" && specialRoomRef.MealPlan != Room.MealPlan {
						continue
					}
					if specialRoomRef.Currency != "" && specialRoomRef.Currency != Room.Currency {
						continue
					}

					isMatched = true
					break
				}
			}

			if !isMatched {
				continue
			}
		}

		str, _ := json.Marshal(roomRef)
		hbeRoomType.Ref = string(str)

		if hbeHotelSearchRequest.Details == true {
			hbeRoomType.Nights = los

			if Room.CancelPolicy != nil {
				hbeRoomType.CancellationPolicy = GenerateCancellationPolicyShort(Room.CancelPolicy, hbeHotelSearchRequest.CheckIn, total, los)
				hbeRoomType.FreeCancellationPolicy = GetFreeCancelDate(Room.CancelPolicy, hbeHotelSearchRequest.CheckIn)

				cancellable := false
				for _, cancelPenalty := range Room.CancelPolicy.CancelPenalties {
					if cancelPenalty.Cancellable {
						cancellable = true
						break
					}
				}
				if cancellable {
					hbeRoomType.NonRefundable = false
				} else {
					hbeRoomType.NonRefundable = true
				}
			}
		}

		hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)
	}

	From(hbeRoomTypes).OrderBy(func(roomType interface{}) interface{} {
		return roomType.(*hbecommon.RoomType).Rate.PerNight
	}).ToSlice(&hbeHotel.RoomTypes)

	if len(hbeHotel.RoomTypes) > 0 {
		hbeHotel.CheapestRoom = hbeHotel.RoomTypes[0]
	}

	if !hbeHotelSearchRequest.Details {
		hbeHotel.RoomTypes = []*hbecommon.RoomType{}

		if hbeHotel.CheapestRoom != nil {
			hbeHotel.CheapestRoom.Ref = ""
			hbeHotel.CheapestRoom.ShortRef = ""
		}
	}
}

func makePreBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, preBookingRequest *PreBookingRequest) {
	preBookingRequest.ReservationIds = &ReservationIds{
		DistributorResId: GenerateDistributorResID(),
		DerbyResId:       GenerateDerbyResID(),
		SupplierResId:    GenerateSupplierResID(),
	}

	hotel := hbeBookingRequest.Hotel
	preBookingRequest.HotelId = hotel.HotelId
	preBookingRequest.StayRange = &StayRange{
		CheckIn:  hbeBookingRequest.CheckIn,
		CheckOut: hbeBookingRequest.CheckOut,
	}
	preBookingRequest.ContactPerson = &Guest{
		FirstName: hbeBookingRequest.Customer.FirstName,
		LastName:  hbeBookingRequest.Customer.LastName,
		Email:     hbeBookingRequest.Customer.Email,
		Phone:     hbeBookingRequest.Customer.PhoneCountryCode + hbeBookingRequest.Customer.PhoneNumber,
	}

	preBookingRequest.Guests = []*Guest{}

	roomCount := 0
	adultCount := 0
	childCount := 0
	childAges := []int{}
	preBookingRequest.RoomRates = []*RoomRate{}

	for index, room := range hotel.Rooms {
		roomCount += room.Count
		adultCount += room.Adults
		childCount += room.Children
		for _, guest := range room.Guests {
			if !guest.IsAdult {
				childAges = append(childAges, guest.Age)
			}

			preBookingRequest.Guests = append(preBookingRequest.Guests, &Guest{
				FirstName: guest.FirstName,
				LastName:  guest.LastName,
				Index:     index,
			})
		}

		var roomRate RoomRate
		_ = json.Unmarshal([]byte(room.Ref), &roomRate)
		preBookingRequest.RoomRates = append(preBookingRequest.RoomRates, &roomRate)
	}

	preBookingRequest.RoomCriteria = &RoomCriteria{
		RoomCount:  roomCount,
		AdultCount: adultCount,
		ChildCount: childCount,
		ChildAges:  childAges,
	}

	preBookingRequest.Total = &Total{
		AmountBeforeTax: hbeBookingRequest.Total,
		AmountAfterTax:  hbeBookingRequest.Total,
	}
}

func makeBookingRequest(preBookRequest *PreBookingRequest, preBookResponse *PreBookResponse) *BookingRequest {
	bookingRequest := &BookingRequest{
		Header:         preBookRequest.Header,
		ReservationIds: preBookRequest.ReservationIds,
		Iata:           preBookRequest.Iata,
		HotelId:        preBookRequest.HotelId,
		StayRange:      preBookRequest.StayRange,
		ContactPerson:  preBookRequest.ContactPerson,
		RoomCriteria:   preBookRequest.RoomCriteria,
		Total:          preBookRequest.Total,
		Payment:        preBookRequest.Payment,
		LoyaltyAccount: preBookRequest.LoyaltyAccount,
		Guests:         preBookRequest.Guests,
		Comments:       preBookRequest.Comments,
		RoomRates:      preBookRequest.RoomRates,
		BookingToken:   preBookResponse.BookingToken,
	}

	return bookingRequest
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

	if bookingResult.ReservationIds != nil {
		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum
		bookingResponse.Booking = &hbecommon.Booking{
			Ref:               bookingResult.ReservationIds.DistributorResId,
			ItineraryId:       bookingResult.ReservationIds.DerbyResId,
			SupplierReference: bookingResult.ReservationIds.SupplierResId,
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

	if cancelResult.CancellationId != "" {
		bookingCancelResponse.Status = "200"
		bookingCancelResponse.InternalRef = cancelResult.ReservationIds.DistributorResId
		bookingCancelResponse.Ref = cancelResult.ReservationIds.DerbyResId
		bookingCancelResponse.CancelRef = cancelResult.CancellationId
	} else {
		bookingCancelResponse.Status = "Failed"
	}
}

func SplitBatchSearchRequest(searchRequest *SearchRequest, maxBatchSize int) []*SearchRequest {

	// batches := []*SearchRequest{}

	// var counter int = 0
	// for counter < len(searchRequest.HotelIDs) {

	// 	batch := searchRequest.Clone()

	// 	var batchSize int = maxBatchSize
	// 	if counter+batchSize > len(searchRequest.HotelIDs) {

	// 		batchSize = len(searchRequest.HotelIDs) - counter
	// 	}

	// 	batch.HotelIDs = searchRequest.HotelIDs[counter : counter+batchSize]

	// 	batches = append(batches, batch)

	// 	counter += batchSize
	// }

	// return batches

	return nil
}
