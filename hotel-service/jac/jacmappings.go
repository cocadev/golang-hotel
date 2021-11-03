package jac

import (
	"bytes"
	"fmt"
	roomresutils "roomres/utils"

	"strings"

	hbecommon "../roomres/hbe/common"

	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_searchrequest(hotelSearchRequest *hbecommon.HotelSearchRequest, searchRequest *SearchRequest, hotelIdTransformer *HotelIdTransformer) {

	searchRequest.SearchDetails = &SearchDetails{}
	searchRequest.SearchDetails.ArrivalDate = hotelSearchRequest.CheckIn
	searchRequest.SearchDetails.Duration = hotelSearchRequest.GetLos()

	searchRequest.Details = hotelSearchRequest.Details

	externalRefs := []string{}

	for _, externalRef := range hotelSearchRequest.ExternalRefs {
		if len(strings.Trim(externalRef, " ")) > 0 {
			externalRefs = append(externalRefs, strings.Trim(externalRef, " "))
		}
	}

	if len(externalRefs) > 0 {

		searchRequest.SearchDetails.PropertyId = externalRefs[0]

	} else {

		searchRequest.SearchDetails.PropertyReferenceIds = []*PropertyReferenceId{}
		for _, hotelId := range hotelSearchRequest.HotelIds {

			searchRequest.SearchDetails.PropertyReferenceIds = append(
				searchRequest.SearchDetails.PropertyReferenceIds,
				&PropertyReferenceId{ReferenceId: hotelIdTransformer.ExtractJacHotelId(hotelId)})
		}
	}

	for _, hbeRequestedRoom := range hotelSearchRequest.RequestedRooms {

		roomRequest := &RoomRequest{}

		mapping_hbe_roomrequest(hbeRequestedRoom, roomRequest)

		searchRequest.SearchDetails.RoomRequests = append(searchRequest.SearchDetails.RoomRequests, roomRequest)
	}
}

const (
	MinChildAge              int = 2
	MaxNumberOfPaxPerBooking int = 9
)

func CheckMaxNumberOfPax(hotelSearchRequest *hbecommon.HotelSearchRequest) bool {

	totalPax := 0

	for _, hbeRequestedRoom := range hotelSearchRequest.RequestedRooms {

		totalPax += hbeRequestedRoom.Adults + len(hbeRequestedRoom.ChildAges)
	}

	return totalPax <= MaxNumberOfPaxPerBooking
}

func CheckMaxNumberOfRooms(hotelSearchRequest *hbecommon.HotelSearchRequest) bool {

	return len(hotelSearchRequest.RequestedRooms) <= 1
}

func mapping_hbe_roomrequest(hbeRequestedRoom *hbecommon.RoomRequest, roomRequest *RoomRequest) {

	roomRequest.Adults = hbeRequestedRoom.Adults

	//roomRequest.Infants = hbeRequestedRoom.Cots
	//roomRequest.Children = hbeRequestedRoom.Children

	roomRequest.Infants = 0
	roomRequest.Children = 0

	for _, childAge := range hbeRequestedRoom.ChildAges {

		if childAge.Age < MinChildAge {
			roomRequest.Infants++
		} else {
			roomRequest.Children++
			roomRequest.ChildAges = append(roomRequest.ChildAges, &ChildAge{Age: childAge.Age})
		}
	}
}

func mapping_searchrequest_to_hbe_searchresponse(
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	hotelIdTransformer *HotelIdTransformer,
) {

	hbeHotel := &hbecommon.Hotel{}

	hbeRoomTypes := []*hbecommon.RoomType{}

	for i, specificRoomRef := range hbeHotelSearchRequest.SpecificRoomRefs {

		var combinedRoomRef CombinedRoomRef
		DecodeCombinedRoomRef(specificRoomRef, &combinedRoomRef)

		if i == 0 {
			hbeHotel.HotelId = hotelIdTransformer.GenerateJacHotelId(combinedRoomRef.PropertyReferenceId)
			hbeHotel.CustomTag = combinedRoomRef.PropertyId
		}

		roomTypes := []*RoomType{}

		for _, roomTypeRef := range combinedRoomRef.RoomTypeRefs {
			roomTypes = append(roomTypes, roomTypeRef.RoomType)
		}

		hbeRoomType := &hbecommon.RoomType{}

		hbeRoomType.Ref = specificRoomRef
		hbeRoomType.ShortRef = specificRoomRef

		mapping_roomtypedetails_to_hbe(roomTypes, hbeRoomType, hbeHotelSearchRequest.GetLos())

		hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)

		if hbeRoomType.Notes != "" {

			if hbeHotel.Notes != "" {
				hbeHotel.Notes += "<br />"
			}

			hbeHotel.Notes += hbeRoomType.Notes
		}
	}

	hbeHotel.RoomTypes = hbeRoomTypes

	hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
}

func mapping_searchresponse_to_hbe(
	searchResponse *SearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	hotelIdTransformer *HotelIdTransformer, logScope roomresutils.ILog) {

	for _, propertyResult := range searchResponse.PropertyResults {

		hbeHotel := &hbecommon.Hotel{}

		mapping_hotel_to_hbe(
			hbeHotelSearchRequest,
			propertyResult,
			hbeHotel,
			hbeHotelSearchRequest.GetLos(),
			hotelIdTransformer,
			logScope,
		)

		if hbeHotel.CheapestRoom != nil {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
		}
	}
}

func mapping_hotel_to_hbe(
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	propertyResult *PropertyResult,
	hbeHotel *hbecommon.Hotel,
	los int, hotelIdTransformer *HotelIdTransformer, logScope roomresutils.ILog) {

	hbeHotel.HotelId = hotelIdTransformer.GenerateJacHotelId(propertyResult.PropertyReferenceId)
	hbeHotel.CustomTag = propertyResult.PropertyId

	hbeRoomTypes := []*hbecommon.RoomType{}

	combinations := FilterRoomTypeCombinations(
		GenerateRoomTypeCombinations(len(hbeHotelSearchRequest.RequestedRooms), propertyResult.RoomTypes))

	if len(combinations) > 200 {

		logScope.LogEvent(
			roomresutils.EventTypeInfo2,
			"JAC mapping_hotel_to_hbe - too many room combinations found",
			fmt.Sprintf("%s - %d", hbeHotel.HotelId, len(combinations)))
	}

	if len(combinations) > 50000 {
		return
	}

	var hotelNotes bytes.Buffer
	hotelNotes.WriteString("")

	for _, combination := range combinations {

		skipCombination := false
		for _, roomType := range combination.RoomTypes {

			if len(roomType.RSP) > 0 {
				skipCombination = true
				break
			}
		}

		if skipCombination {
			continue
		}

		hbeRoomType := &hbecommon.RoomType{}

		mapping_roomtype_to_hbe(hbeHotelSearchRequest, propertyResult, combination.RoomTypes, hbeRoomType, los)

		var includeRoomType bool = true
		if len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {

			includeRoomType = false
			for _, roomTypeRef := range hbeHotelSearchRequest.SpecificRoomRefs {
				if CompareRoomRefs(roomTypeRef, hbeRoomType.Ref) {

					includeRoomType = true
					break
				}
			}
		}

		if includeRoomType {
			hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)

			if hbeHotelSearchRequest.Details && hbeRoomType.Notes != "" && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {

				if hotelNotes.Len() > 0 {
					hotelNotes.WriteString("<br />")
				}

				hotelNotes.WriteString(hbeRoomType.Notes)
				//hbeHotel.Notes += hbeRoomType.Notes
			}
		}

	}

	hbeHotel.Notes = hotelNotes.String()
	hotelNotes.Reset()

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

func mapping_hbe_prebookrequest(
	roomRef string,
	preBookRequest *PreBookRequest) {

	var combinedRoomRef CombinedRoomRef

	DecodeCombinedRoomRef(roomRef, &combinedRoomRef)

	mapping_hbe_to_prebookrequest(combinedRoomRef.ArrivalDate,
		combinedRoomRef.Duration,
		&combinedRoomRef,
		preBookRequest)
}

/*
func mapping_hbe_bookingrequest_to_prebookrequest(
	hbeBookingRequest *hbecommon.BookingRequest,
	preBookRequest *PreBookRequest) {

	var combinedRoomRef CombinedRoomRef

	DecodeCombinedRoomRef(hbeBookingRequest.Hotel.Rooms[0].Ref, &combinedRoomRef)

	mapping_hbe_to_prebookrequest(hbeBookingRequest.CheckIn,
		hbeBookingRequest.GetLos(),
		&combinedRoomRef,
		preBookRequest)
}
*/

func mapping_hbe_bookingrequest_to_prebookrequest(
	hbeBookingRequest *hbecommon.BookingRequest,
	preBookRequest *PreBookRequest) {

	mapping_hbe_prebookrequest(hbeBookingRequest.Hotel.Rooms[0].Ref, preBookRequest)
}

func mapping_hbe_to_prebookrequest(
	checkIn string,
	duration int,
	combinedRoomRef *CombinedRoomRef,
	preBookRequest *PreBookRequest,
) {

	preBookRequest.BookingDetails = &BookingDetails{
		PropertyId:  combinedRoomRef.PropertyId,
		ArrivalDate: checkIn,
		Duration:    duration}

	for _, roomTypeRef := range combinedRoomRef.RoomTypeRefs {

		roomBooking := &RoomBooking{}

		mapping_roomref_to_roombooking(roomTypeRef, roomBooking)

		preBookRequest.BookingDetails.RoomBookings = append(
			preBookRequest.BookingDetails.RoomBookings,
			roomBooking,
		)
	}
}

func mapping_roomref_to_roombooking(roomTypeRef *RoomTypeRef, roomBooking *RoomBooking) {

	roomBooking.PropertyRoomTypeId = roomTypeRef.RoomTypeId
	roomBooking.BookingToken = roomTypeRef.BookingToken
	roomBooking.MealBasisId = roomTypeRef.MealBasisId

	roomBooking.Adults = roomTypeRef.RequestedRoom.Adults
	roomBooking.Children = roomTypeRef.RequestedRoom.Children
	roomBooking.Infants = roomTypeRef.RequestedRoom.Infants

	roomBooking.ChildAges = roomTypeRef.RequestedRoom.ChildAges
}

func mapping_prebookresponse_to_hbe(
	commissionProvider hbecommon.ICommissionProvider,
	preBookResponse *PreBookResponse,
	hbeRoomType *hbecommon.RoomType) {

	//hbeRoomType.NonRefundable = IsNonRefundableRoom(hbeRoomType.Name) //IsNonRefundableCancellation(preBookResponse.Cancellations)
	hbeRoomType.FreeCancellationPolicy = GenerateCancellationPolicyShort(preBookResponse.Cancellations)
	hbeRoomType.CancellationPolicy = GenerateCancellationPolicy(
		commissionProvider,
		preBookResponse.Cancellations,
		0)
}

func mapping_prebookresponse_to_roomcxlresponse(
	preBookResponse *PreBookResponse,
	hbeRoomCxlResponse *hbecommon.RoomCxlResponse) {

	hbeRoomCxlResponse.Description = GenerateCancellationPolicyShort(preBookResponse.Cancellations)
}

func mapping_hbe_to_bookingrequest(hbeBookingRequest *hbecommon.BookingRequest, preBookResponse *PreBookResponse, bookRequest *BookRequest) {

	bookRequest.BookDetails = &BookDetails{}

	mapping_hbe_to_bookdetails(hbeBookingRequest, preBookResponse, bookRequest.BookDetails)
}

func mapping_hbe_to_bookdetails(hbeBookingRequest *hbecommon.BookingRequest, preBookResponse *PreBookResponse, bookDetails *BookDetails) {

	bookDetails.PreBookingToken = preBookResponse.PreBookingToken

	bookDetails.TradeReference = hbeBookingRequest.InternalRef
	bookDetails.ArrivalDate = hbeBookingRequest.CheckIn
	bookDetails.Duration = hbeBookingRequest.GetLos()

	bookDetails.LeadGuestTitle = hbeBookingRequest.Customer.Title
	bookDetails.LeadGuestFirstName = hbeBookingRequest.Customer.FirstName
	bookDetails.LeadGuestLastName = hbeBookingRequest.Customer.LastName
	bookDetails.LeadGuestEmail = hbeBookingRequest.Customer.Email
	bookDetails.SpecialRequest = CombinedSpecialRequests(hbeBookingRequest)

	for i, hbeBookingRoom := range hbeBookingRequest.Hotel.Rooms {

		var combinedRoomRef CombinedRoomRef
		DecodeCombinedRoomRef(hbeBookingRoom.Ref, &combinedRoomRef)

		if i == 0 {
			bookDetails.PropertyId = combinedRoomRef.PropertyId
		}

		bookRoom := &RoomBook{}

		mapping_hbe_to_bookroom(hbeBookingRequest, &combinedRoomRef, hbeBookingRoom, bookRoom)

		bookDetails.BookRooms = append(bookDetails.BookRooms, bookRoom)
	}

}

func mapping_hbe_to_bookroom(
	hbeBookingRequest *hbecommon.BookingRequest,
	combinedRoomRef *CombinedRoomRef,
	hbeBookingRoom *hbecommon.BookingRoom,
	bookRoom *RoomBook) {

	bookRoom.PropertyRoomTypeId = combinedRoomRef.RoomTypeRefs[0].RoomTypeId
	bookRoom.MealBasisId = combinedRoomRef.RoomTypeRefs[0].MealBasisId
	bookRoom.BookingToken = combinedRoomRef.RoomTypeRefs[0].BookingToken

	bookRoom.Adults = hbeBookingRoom.Adults
	bookRoom.Children = 0
	bookRoom.Infants = 0

	for _, hbeGuest := range hbeBookingRoom.Guests {

		if !hbeGuest.IsAdult && hbeGuest.Age < MinChildAge {
			bookRoom.Infants++
		} else if !hbeGuest.IsAdult {
			bookRoom.Children++
		}

		if hbeGuest.IsAdult || hbeGuest.Age >= MinChildAge {

			guest := &Guest{}
			mapping_hbe_to_guest(hbeBookingRequest, hbeGuest, guest)
			bookRoom.Guests = append(bookRoom.Guests, guest)
		}

	}
}

func mapping_hbe_to_guest(
	hbeBookingRequest *hbecommon.BookingRequest,
	hbeGuest *hbecommon.Guest,
	guest *Guest) {

	if hbeGuest.IsAdult {
		guest.Type = "Adult"
	} else {
		guest.Type = "Child"
	}
	guest.Title = hbeGuest.Title
	guest.FirstName = hbeGuest.FirstName
	guest.LastName = hbeGuest.LastName
	guest.Age = hbeGuest.Age
}

func mapping_hbe_to_precancelrequest(
	bookingCancelRequest *hbecommon.BookingCancelRequest,
	preCancelRequest *PreCancelRequest) {

	preCancelRequest.BookingReference = bookingCancelRequest.Ref
}

func mapping_precancelresponse_to_hbe(
	preCancelResponse *PreCancelResponse,
	bookingCancelRequest *hbecommon.BookingCancelRequest,
	bookingCancelResponse *hbecommon.BookingCancelResponse) {

	if preCancelResponse.ReturnStatus.IsSuccess() {
		bookingCancelResponse.Status = "200"

		bookingCancelResponse.Ref = preCancelResponse.BookingReference
		bookingCancelResponse.InternalRef = bookingCancelRequest.InternalRef

		bookingCancelResponse.CancelRef =
			fmt.Sprintf("%s|%.2f", preCancelResponse.CancellationToken, preCancelResponse.CancellationCost)

		bookingCancelResponse.PolicyText = fmt.Sprintf("Internal Cancellation Cost: $%.2f", preCancelResponse.CancellationCost)
	} else {

		bookingCancelResponse.Status = "Failed"
	}
}

func mapping_bookresponse_to_hbe(bookResponse *BookResponse, bookingResponse *hbecommon.BookingResponse) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	if bookResponse.ReturnStatus.IsSuccess() {

		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		var supplierReference, supplierName string = "", ""

		if len(bookResponse.PropertyBookings) > 0 {
			supplierReference = bookResponse.PropertyBookings[0].SupplierReference
			supplierName = bookResponse.PropertyBookings[0].Supplier
		}

		bookingResponse.Booking = &hbecommon.Booking{
			Ref:               bookResponse.BookingReference,
			ItineraryId:       "",
			SupplierReference: supplierReference,
			SupplierName:      supplierName,
			Total:             bookResponse.TotalPrice}

	} else {

		if strings.Index(strings.ToLower(bookResponse.ReturnStatus.Exception), strings.ToLower(BookingFailureExceptionType1)) >= 0 ||
			strings.Index(strings.ToLower(bookResponse.ReturnStatus.Exception), strings.ToLower(BookingFailureExceptionType2)) >= 0 {

			bookingResponse.BookingStatus = hbecommon.BookingStatusFailedRestartEnum
			bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{Message: bookResponse.ReturnStatus.Exception}}

		} else {

			bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
			bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{Message: bookResponse.ReturnStatus.Exception}}
		}

	}
}

func mapping_propertyid_to_propertydetailsrequest(propertyId string, propertyDetailsRequest *PropertyDetailsRequest) {

	propertyDetailsRequest.PropertyId = propertyId
}

func mapping_cancelresponse_to_hbe(
	cancelResponse *CancelResponse,
	bookingCancelConfirmationResponse *hbecommon.BookingCancelConfirmationResponse) {

	if cancelResponse.ReturnStatus.IsSuccess() {
		bookingCancelConfirmationResponse.Status = "200"
	} else {
		bookingCancelConfirmationResponse.Status = "Failed"
	}
}

func mapping_hbe_to_cancelrequest(
	bookingCancelConfirmationRequest *hbecommon.BookingCancelConfirmationRequest,
	cancelRequest *CancelRequest) {

	cancelRequest.BookingReference = bookingCancelConfirmationRequest.Ref

	values := strings.Split(bookingCancelConfirmationRequest.CancelRef, "|")

	cancelRequest.CancellationToken = values[0]
	cancelRequest.CancellationCost = values[1]
}

func mapping_roomtypedetails_to_hbe(roomTypes []*RoomType, hbeRoomType *hbecommon.RoomType, los int) {

	hbeRoomType.Description = roomTypes[0].RoomType //GenerateCombinedRoomTypeName(roomTypes)
	hbeRoomType.Promo = false
	hbeRoomType.Breakfast = IsBB(roomTypes)
	hbeRoomType.NonRefundable = IsNonRefundableRoom(hbeRoomType.Description)

	hbeRoomType.BoardTypeText = roomTypes[0].MealBasis

	hbeRoomType.RemainingRooms = len(roomTypes)

	//// mapping roomtype occupancy need to get from db

	hbeRoomType.CancellationPolicy = "" // to get it from prebook request
	hbeRoomType.Notes = GenerateRoomTypeNotes(roomTypes[0])

	hbeRoomType.Rate = &hbecommon.Rate{}

	// roomtypes are identical
	/*
		for _, roomType := range roomTypes {

			hbeRoomType.Rate.PerNight += roomType.Total / float32(los)
			hbeRoomType.Rate.PerNightBase += roomType.SubTotal / float32(los)
		}
	*/

	hbeRoomType.Rate.PerNight = roomTypes[0].Total / float32(los)
	hbeRoomType.Rate.PerNightBase = roomTypes[0].Total / float32(los)

	hbeRoomType.Surcharges = []*hbecommon.Surcharge{}
	hbeRoomType.Taxes = []*hbecommon.Tax{}

	hbeRoomType.NormalBeddingOccupancy = roomTypes[0].Adults + roomTypes[0].Children
	hbeRoomType.ExtraBeddingOccupancy = 0

}

func mapping_roomtype_to_hbe(hbeHotelSearchRequest *hbecommon.HotelSearchRequest, hotel *PropertyResult, roomTypes []*RoomType, hbeRoomType *hbecommon.RoomType, los int) {

	combinedRoomRef := EncodeCombinedRoomRef(
		GenerateCombinedRoomRef(
			hotel.PropertyId,
			hotel.PropertyReferenceId,
			hbeHotelSearchRequest.CheckIn,
			hbeHotelSearchRequest.GetLos(),
			roomTypes,
		))

	hbeRoomType.ShortRef = combinedRoomRef
	hbeRoomType.Ref = combinedRoomRef

	mapping_roomtypedetails_to_hbe(roomTypes, hbeRoomType, los)
}
