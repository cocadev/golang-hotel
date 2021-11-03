package hoteldo

import (
	"encoding/json"
	"fmt"
	"time"

	hbecommon "../roomres/hbe/common"

	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_searchrequest(hotelSearchRequest *hbecommon.HotelSearchRequest, searchRequest *SearchRequest, hotelIdTransformer *HotelIdTransformer) {

	//	searchRequest.SearchDetails = &SearchDetails{}
	searchRequest.SearchDetails.StartDate = hotelSearchRequest.CheckIn
	searchRequest.SearchDetails.EndDate = hotelSearchRequest.CheckOut
	searchRequest.SearchDetails.MealPlan = ""
	searchRequest.SearchDetails.NumberOfRooms = len(hotelSearchRequest.RequestedRooms)
	searchRequest.SearchDetails.DestinationCode = 2
	searchRequest.SearchDetails.LanguageId = "ING"
	searchRequest.SearchDetails.Order = hotelSearchRequest.SortType

	searchRequest.SearchDetails.Details = hotelSearchRequest.Details

	for _, hotelId := range hotelSearchRequest.HotelIds {
		searchRequest.SearchDetails.HotelIds = append(searchRequest.SearchDetails.HotelIds, hotelId)
	}

	for _, specificRoomRef := range hotelSearchRequest.SpecificRoomRefs {
		searchRequest.SearchDetails.SpecificRoomRefs = append(searchRequest.SearchDetails.SpecificRoomRefs, specificRoomRef)
	}

	for _, hbeRequestedRoom := range hotelSearchRequest.RequestedRooms {

		roomRequest := &RoomRequest{}

		mapping_hbe_roomrequest(hbeRequestedRoom, roomRequest)

		searchRequest.SearchDetails.RequestedRooms = append(searchRequest.SearchDetails.RequestedRooms, roomRequest)
	}
}

func mapping_hbe_roomrequest(hbeRequestedRoom *hbecommon.RoomRequest, roomRequest *RoomRequest) {

	roomRequest.Adults = hbeRequestedRoom.Adults
	roomRequest.Children = 0

	for _, childAge := range hbeRequestedRoom.ChildAges {
		roomRequest.Children++
		roomRequest.ChildAges = append(roomRequest.ChildAges, &ChildAge{Age: childAge.Age})
	}
}

func mapping_searchresponse_to_hbe(
	searchResponse *SearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	hotelIdTransformer *HotelIdTransformer) {

	los := hbeHotelSearchRequest.GetLos()

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

	hbeHotel.HotelId = fmt.Sprintf("%d", hotelItem.Id) // hotelIdTransformer.GenerateJacHotelId(propertyResult.PropertyReferenceId)
	hbeHotel.CustomTag = hotelItem.CategoryId          // propertyResult.PropertyId
	hbeHotel.Notes = hotelItem.Description

	hbeRoomTypes := []*hbecommon.RoomType{}

	for _, Room := range hotelItem.Rooms {

		hbeRoomType := &hbecommon.RoomType{}
		hbeRoomType.Description = Room.Name

		for _, MP := range Room.MealPlans {
			var Total float32 = 0
			for _, NightlyRate := range MP.NightsDetail {
				Total += NightlyRate.Total
			}
			hbeRoomType.Rate = &hbecommon.Rate{}
			hbeRoomType.Rate.PerNight = Total / float32(los)
			hbeRoomType.Rate.PerNightBase = Total / float32(los)
			hbeRoomType.Rate.CurrecyCode = MP.Currency
			hbeRoomType.Rate.NightlyRateTotal = Total

			adults, kids, k1a, k2a, k3a := getResources(MP.RateDetails[0].PaxCount)
			roomRef := RoomRef{
				Amount:  Total,
				Status:  MP.RateDetails[0].Available.Status,
				RateKey: MP.RateDetails[0].RateKey,
				Adults:  adults,
				Kids:    kids,
			}

			if k1a != "" {
				roomRef.K1a = k1a
			}
			if k2a != "" {
				roomRef.K2a = k2a
			}
			if k3a != "" {
				roomRef.K3a = k3a
			}

			//additionals
			roomRef.MarketId = Room.MealPlans[0].MarketId
			roomRef.Contract = Room.MealPlans[0].Contract
			roomRef.Currency = Room.MealPlans[0].Currency
			roomRef.Mealplan = Room.MealPlans[0].Id
			roomRef.RoomType = Room.Id

			//check special refs
			if hbeHotelSearchRequest.SpecificRoomRefs != nil && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
				var isMatched = false
				for _, roomRefStr := range hbeHotelSearchRequest.SpecificRoomRefs {
					var specialRoomRef RoomRef
					e := json.Unmarshal([]byte(roomRefStr), &specialRoomRef)
					if e == nil {
						if specialRoomRef.RoomType != "" && specialRoomRef.RoomType != roomRef.RoomType {
							continue
						}
						if specialRoomRef.MarketId != "" && specialRoomRef.MarketId != roomRef.MarketId {
							continue
						}
						if specialRoomRef.Mealplan != "" && specialRoomRef.Mealplan != roomRef.Mealplan {
							continue
						}
						if specialRoomRef.RateKey != "" && specialRoomRef.RateKey != roomRef.RateKey {
							continue
						}
						if specialRoomRef.Amount != roomRef.Amount {
							continue
						}
						if specialRoomRef.Status != "" && specialRoomRef.Status != roomRef.Status {
							continue
						}
						if specialRoomRef.Adults != "" && specialRoomRef.Adults != roomRef.Adults {
							continue
						}
						if specialRoomRef.Kids != "" && specialRoomRef.Kids != roomRef.Kids {
							continue
						}
						if specialRoomRef.K1a != "" && specialRoomRef.K1a != roomRef.K1a {
							continue
						}
						if specialRoomRef.K2a != "" && specialRoomRef.K2a != roomRef.K2a {
							continue
						}
						if specialRoomRef.K3a != "" && specialRoomRef.K3a != roomRef.K3a {
							continue
						}
						isMatched = true
						break
					}
				}

				if !isMatched {
					return
				}
			}

			roomRefStr, err := json.Marshal(roomRef)
			if err == nil {
				hbeRoomType.Ref = string(roomRefStr)
				hbeRoomType.ShortRef = MP.RateDetails[0].RateKey
			}

			if hbeHotelSearchRequest.Details == true {
				hbeRoomType.Nights = len(MP.NightsDetail)

				if hbeRoomType.NonRefundable == true {
					hbeRoomType.FreeCancellationPolicy = ""
				}

				if (len(MP.RateDetails) > 0) && (MP.RateDetails[0].CancellationPolicy != nil) {
					hbeRoomType.CancellationPolicy = GenerateCancellationPolicyShort(MP.RateDetails[0].CancellationPolicy, hbeHotelSearchRequest.CheckIn)
					hbeRoomType.FreeCancellationPolicy = GetFreeCancelDate(MP.RateDetails[0].CancellationPolicy, hbeHotelSearchRequest.CheckIn)
				}

				hbeRoomType.Taxes = []*hbecommon.Tax{}
				hbeRoomType.Taxes = append(hbeRoomType.Taxes, &hbecommon.Tax{
					Type:   "Tax",
					Amount: (MP.RateDetails[0].Total - MP.RateDetails[0].GrossTotal),
				})

				if len(MP.Promotions) > 0 {
					hbeRoomType.Promo = true
				} else {
					hbeRoomType.Promo = false
				}
			}

			break //  to get 1 information for 1 room ; originally 1 exist
		}
		if hbeHotelSearchRequest.Details == true {
			hbeRoomType.NormalBeddingOccupancy = Room.CapacityAdults + Room.CapacityKids
			hbeRoomType.ExtraBeddingOccupancy = Room.CapacityExtras
			hbeRoomType.Notes = Room.RoomView
			hbeRoomType.BoardTypeText = Room.Bedding
			hbeRoomType.NonRefundable = IsNonRefundableRoom(Room.Name)
			hbeRoomType.RemainingRooms = len(hotelItem.Rooms)
		}

		hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)
	}

	From(hbeRoomTypes).OrderBy(func(roomType interface{}) interface{} {
		return roomType.(*hbecommon.RoomType).Rate.PerNight
	}).ToSlice(&hbeHotel.RoomTypes)

	if len(hbeHotel.RoomTypes) > 0 {
		hbeHotel.CheapestRoom = hbeHotel.RoomTypes[0]
	}
}

func makeBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, bookingRequest *BookingRequest) {
	bookingRequest.Firstname = hbeBookingRequest.Customer.FirstName
	bookingRequest.Lastname = hbeBookingRequest.Customer.LastName
	bookingRequest.Emailaddress = hbeBookingRequest.Customer.Email
	bookingRequest.Total = hbeBookingRequest.Total

	bookingRequest.Phones = []*Phone{}
	bookingRequest.Phones = append(bookingRequest.Phones, &Phone{
		Type:   "2",
		Number: hbeBookingRequest.Customer.PhoneNumber,
	})

	bookingRequest.Hotels = []*HotelBook{}
	//hotelId := hbeBookingRequest.Hotel.HotelId
	hotel := hbeBookingRequest.Hotel
	var roomRef0 RoomRef
	err := json.Unmarshal([]byte(hotel.Rooms[0].Ref), &roomRef0)
	if err != nil {
		fmt.Printf("Request parsing error")
		return
	}

	var payment CreditPayment
	err = json.Unmarshal([]byte(hbeBookingRequest.InternalRef), &payment)
	if err != nil {
		fmt.Printf("Request parsing error")
		return
	}
	bookingRequest.CreditPayment = &payment

	fd, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckIn)
	td, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckOut)

	hotelBook := &HotelBook{
		Hotelid:       hotel.HotelId,
		Roomtype:      roomRef0.RoomType,
		Mealplan:      roomRef0.Mealplan,
		Datearrival:   fd.Format("20060102"),
		Datedeparture: td.Format("20060102"),
		Marketid:      roomRef0.MarketId,
		Contractid:    roomRef0.Contract,
		Currency:      roomRef0.Currency,
		Rooms:         []*RoomRef{},
	}

	for _, room := range hotel.Rooms {
		var roomRef RoomRef
		json.Unmarshal([]byte(room.Ref), &roomRef)
		hotelBook.Rooms = append(hotelBook.Rooms, &roomRef)
	}

	bookingRequest.Hotels = append(bookingRequest.Hotels, hotelBook)
}

func mapping_bookresponse_to_hbe(
	bookResponse *BookingResponse,
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

	if bookingResult.Confirmationid != "0" {
		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		bookingResponse.Booking = &hbecommon.Booking{
			Ref:          bookingResult.Confirmationid,
			SupplierName: bookingResult.Operatorname,
			Total:        bookingResult.Total,
		}

		if bookingResult.Rooms != nil {
			room := bookingResult.Rooms[0]
			bookingResponse.Booking.CancellationPolicy = GenerateCancellationPolicyShort(room.CancellationPolicy, hbeBookingRequest.CheckIn)
			bookingResponse.Booking.FreeCancellationPolicy = GetFreeCancelDate(room.CancellationPolicy, hbeBookingRequest.CheckIn)
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

	if cancelResult.Status == "CA" {
		bookingCancelResponse.Status = "200"
		bookingCancelResponse.Ref = cancelResult.Number
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
	for counter < len(searchRequest.SearchDetails.HotelIds) {

		batch := searchRequest.Clone()

		var batchSize int = maxBatchSize
		if counter+batchSize > len(searchRequest.SearchDetails.HotelIds) {

			batchSize = len(searchRequest.SearchDetails.HotelIds) - counter
		}

		batch.SearchDetails.HotelIds = searchRequest.SearchDetails.HotelIds[counter : counter+batchSize]

		batches = append(batches, batch)

		counter += batchSize
	}

	return batches
}
