package teamamerica

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"

	. "github.com/ahmetb/go-linq"
)

func HbeToTeamAmericaHotelSearchRequests(hbeSearchRequest *hbecommon.HotelSearchRequest) []*TeamAmericaHotelSearchRequest {
	tempCheckIn, _ := time.Parse("2006-01-02", hbeSearchRequest.CheckIn)
	tempCheckOut, _ := time.Parse("2006-01-02", hbeSearchRequest.CheckOut)

	NumberOfNights := int(math.Floor(tempCheckOut.Sub(tempCheckIn).Hours() / 24))

	hotelIds := hbeSearchRequest.HotelIds
	var VendorIDs *VendorIDs = &VendorIDs{
		VendorID: []string{},
	}

	results := []*TeamAmericaHotelSearchRequest{}
	for _, hotelId := range hotelIds {
		compoundhotelRawId := GetRawHotelId(hotelId)
		params := strings.Split(compoundhotelRawId, "|")
		VendorIDs.VendorID = append(VendorIDs.VendorID, params[0])

		roomRequests := hbeSearchRequest.RequestedRooms

		ret := &TeamAmericaHotelSearchRequest{
			CityCode:       params[1],
			ProductCode:    params[2],
			Type:           "Hotel",
			ArrivalDate:    hbeSearchRequest.CheckIn,
			VendorIDs:      VendorIDs,
			NumberOfNights: NumberOfNights,
			NumberOfRooms:  len(roomRequests),
		}

		//fmt.Printf("Count of room request : %d\n", len(roomRequests))

		//consider all room condition is same so only process first one

		if len(roomRequests) >= 1 {
			roomRequest := roomRequests[0]
			ret.Occupancy = makeOccupancyString(roomRequest.Adults, roomRequest.Children)
		}

		results = append(results, ret)
	}

	return results
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

	//fmt.Printf("CheckMaxNumberOfPax - %d, %d\n", totalPax, MaxNumberOfPaxPerBooking)
	return totalPax <= MaxNumberOfPaxPerBooking
}

func CheckMaxNumberOfRooms(hotelSearchRequest *hbecommon.HotelSearchRequest) bool {

	//fmt.Printf("CheckMaxNumberOfRooms - %d\n", len(hotelSearchRequest.RequestedRooms))
	return len(hotelSearchRequest.RequestedRooms) >= 1
}

func mapping_searchresponse_to_hbe(
	searchResponse *TeamAmericaHotelSearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	roomCxlResponses []*TeamAmericaCancelPolicyResponse,
	//logScope roomresutils.ILog
) {
	roomTypes := 0

	for _, hotel := range searchResponse.HotelSearchResponse.HotelDatas {

		hbeHotel := &hbecommon.Hotel{}

		hotelId := fmt.Sprintf("%s_%s", "ta", hotel.CompoundHotelId)
		isNewHotel := true

		//check if current hotel is registered already
		for _, t := range hbeHotelSearchResponse.Hotels {
			if t.HotelId == hotelId {
				hbeHotel = t
				isNewHotel = false
				break
			}
		}

		if hbeHotel.HotelId == "" {
			hbeHotel.HotelId = hotelId
		}

		hbeRoomTypes := []*hbecommon.RoomType{}
		if hbeHotel.RoomTypes != nil {
			hbeRoomTypes = hbeHotel.RoomTypes
		}

		//checkOut, _ := time.Parse("2006-01-02", hbeHotelSearchRequest.CheckOut)
		//checkIn, _ := time.Parse("2006-01-02", hbeHotelSearchRequest.CheckIn)

		//deltaDays := float32(math.Floor(checkOut.Sub(checkIn).Hours() / 24))

		//checkDetail := false

		var total float32 = 0
		for _, nightInfo := range hotel.NightlyInfo {
			total += nightInfo.Prices.AdultPrice
		}

		nonrefundable := false
		if hotel.NonRefundable == 1 {
			nonrefundable = true
		}

		//Removed by Li, 2018/10/20
		// var searchRoomRefs []*SearchRoomRef
		// if hbeHotelSearchRequest.Details && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
		// 	//checkDetail = true
		// 	for _, roomRefStr := range hbeHotelSearchRequest.SpecificRoomRefs {
		// 		var searchRoomRef SearchRoomRef
		// 		err := json.Unmarshal([]byte(roomRefStr), &searchRoomRef)
		// 		if err != nil {
		// 			fmt.Printf("Cannot parse specialRoomRef = %s\n", roomRefStr)
		// 		} else {
		// 			searchRoomRefs = append(searchRoomRefs, &searchRoomRef)
		// 		}
		// 	}
		// 	//check special condition, let's consider specialRoomRef is one
		// 	checkCondition := searchRoomRefs[0]
		// 	if checkCondition.ProductCode != "" && checkCondition.ProductCode != hotel.ProductCode {
		// 		continue
		// 	}
		// 	if checkCondition.ProductDate != "" && checkCondition.ProductDate != hotel.ProductDate {
		// 		continue
		// 	}
		// 	if checkCondition.RoomType != "" && checkCondition.RoomType != hotel.RoomType {
		// 		continue
		// 	}
		// 	if checkCondition.MealPlan != "" && checkCondition.MealPlan != hotel.MealPlan {
		// 		continue
		// 	}
		// 	if checkCondition.ChildAge != 0 && checkCondition.ChildAge != hotel.ChildAge {
		// 		continue
		// 	}
		// 	if checkCondition.FamilyPlan != "" && checkCondition.FamilyPlan != hotel.FamilyPlan {
		// 		continue
		// 	}
		// 	if checkCondition.NonRefundable != 0 && checkCondition.NonRefundable != hotel.NonRefundable {
		// 		continue
		// 	}
		// 	if checkCondition.MaxOccupancy != 0 && checkCondition.ProductCode != hotel.ProductCode {
		// 		continue
		// 	}
		// 	if checkCondition.MinPrice != 0 && checkCondition.MinPrice > total {
		// 		continue
		// 	}
		// 	if checkCondition.MaxPrice != 0 && checkCondition.MaxPrice < total {
		// 		continue
		// 	}
		// }
		refString := &RoomRef{
			ProductCode:        hotel.ProductCode,
			ProductDate:        hotel.ProductDate,
			MealPlan:           hotel.MealPlan,
			ChildAge:           hotel.ChildAge,
			RoomType:           hotel.RoomType,
			FamilyPlan:         hotel.FamilyPlan,
			NonRefundable:      hotel.NonRefundable,
			MaxOccupancy:       hotel.MaxOccupancy,
			AverageNightlyRate: hotel.AverageRate.AverageNightlyRate,
		}
		ref, _ := json.Marshal(refString)
		//fmt.Printf("Ref=%s\n", ref)
		roomType := &hbecommon.RoomType{
			Rate: &hbecommon.Rate{
				PerNight:     hotel.AverageRate.AverageNightlyRate,
				PerNightBase: hotel.AverageRate.AverageNightlyRate,
				Total:        total,
			},
			Description: hotel.RoomType,
			//ShortRef:      hotel.MealPlan,
			Ref:           string(ref),
			Nights:        len(hotel.NightlyInfo),
			NonRefundable: nonrefundable,
		}
		if hbeHotelSearchRequest.Details && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
			//set cancellation policy here
			roomCxlResponse := From(roomCxlResponses).
				SingleWith(
					func(f interface{}) bool {
						return f.(*TeamAmericaCancelPolicyResponse).CancellationPolicyResponse.ProductCode == hotel.ProductCode
					},
				).(*TeamAmericaCancelPolicyResponse)
			if roomCxlResponse != nil && len(roomCxlResponse.CancellationPolicyResponse.CancelPolicies) > 0 {
				roomType.FreeCancellationPolicy = GetFreeCancelDate(
					roomCxlResponse.CancellationPolicyResponse.CancelPolicies,
					hbeHotelSearchRequest.CheckIn)

				roomType.CancellationPolicy = GenerateCancellationPolicyShort(
					roomCxlResponse.CancellationPolicyResponse.CancelPolicies,
					hbeHotelSearchRequest.CheckIn,
					hotel.AverageRate.AverageNightlyRate,
					total,
				)
			} else {
				roomType.NonRefundable = true
			}
		}

		hbeRoomTypes = append(hbeRoomTypes, roomType)
		hbeHotel.RoomTypes = hbeRoomTypes

		if isNewHotel {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
		}
	}

	for _, hbeHotel := range hbeHotelSearchResponse.Hotels {
		hbeRoomTypes := hbeHotel.RoomTypes
		From(hbeRoomTypes).OrderBy(func(roomType interface{}) interface{} {
			return roomType.(*hbecommon.RoomType).Rate.PerNight
		}).ToSlice(&hbeHotel.RoomTypes)

		//ShortRefs set as incremental number
		for idx, roomType := range hbeHotel.RoomTypes {
			roomType.ShortRef = ConvertIntToString(idx + 1)
		}

		if len(hbeHotel.RoomTypes) > 0 {
			//fmt.Println(len(hbeHotel.RoomTypes))
			hbeHotel.CheapestRoom = hbeHotel.RoomTypes[0]
		}

		if !hbeHotelSearchRequest.Details {
			hbeHotel.RoomTypes = []*hbecommon.RoomType{}

			if hbeHotel.CheapestRoom != nil {
				hbeHotel.CheapestRoom.Ref = ""
				hbeHotel.CheapestRoom.ShortRef = ""
			}
		}

		roomTypes += len(hbeHotel.RoomTypes)
	}

	fmt.Printf("Result : %d Hotel(s), %d Room Type(s)\n", len(hbeHotelSearchResponse.Hotels), roomTypes)
}

func makeBookingRequest(hbeBookingRequest *hbecommon.BookingRequest, bookingRequest *TeamAmericaBookReserveRequest) {
	tempCheckIn, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckIn)
	tempCheckOut, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckOut)

	NumberOfNights := int(math.Floor(tempCheckOut.Sub(tempCheckIn).Hours() / 24))

	rooms := []*Item{}
	//hotelId := hbeBookingRequest.Hotel.HotelId
	for _, room := range hbeBookingRequest.Hotel.Rooms {
		var roomRef RoomRef
		json.Unmarshal([]byte(room.Ref), &roomRef)
		newPassengers := []*Passenger{}
		for _, guest := range room.Guests {
			PassengerType := "AD"
			if !guest.IsAdult {
				PassengerType = "CH"
			}
			newPassengers = append(newPassengers, &Passenger{
				Salutation:    guest.Title,
				FamilyName:    guest.LastName,
				FirstName:     guest.FirstName,
				PassengerType: PassengerType,
				PassengerAge:  guest.Age,
			})
		}

		rooms = append(rooms, &Item{
			ProductCode:    roomRef.ProductCode,
			ProductDate:    roomRef.ProductDate,
			Occupancy:      makeOccupancyString(room.Adults, room.Children),
			NumberOfNights: NumberOfNights,
			Language:       "ENG",
			Quantity:       room.Count,
			Passengers: &Passengers{
				NewPassengers: newPassengers,
			},
			BelongsToPackage: 0,
			PackageCode:      0,
			RateExpected:     roomRef.AverageNightlyRate,
		})
	}

	bookingRequest.Items = &Items{
		NewItems: rooms,
	}
}

// func makeOccupantsNames(room *hbecommon.BookingRoom) string {
// 	var guestNames []string //name1#surname1#.#age1#
// 	for _, guest := range room.Guests {
// 		guestNames = append(guestNames, guest.LastName+"#"+guest.FirstName+"#.#"+ConvertIntToString(guest.Age))
// 	}

// 	return strings.Join(guestNames, "@") + "@"
// }

// func mapping_hbe_bookingrequest_to_prebookrequest(
// 	hbeBookingRequest *hbecommon.BookingRequest) []*PreBookRequest {
// 	preBookRequests := []*PreBookRequest{}
// 	idx := 0
// 	fmt.Printf("Room counts = %d\n", len(hbeBookingRequest.Hotel.Rooms))
// 	for idx = 0; idx < len(hbeBookingRequest.Hotel.Rooms); idx++ {
// 		preBookRequests = append(preBookRequests, &PreBookRequest{
// 			TeamAmericaBookingParam: mapping_hbe_prebookrequest(hbeBookingRequest, idx),
// 		})
// 	}

// 	return preBookRequests
// }

// func makeBookingConfirmRequest(
// 	hbeBookingRequest *hbecommon.BookingRequest,
// 	preBookingResponse *PreBookResponse,
// 	confirmBookRequest *ConfirmBookRequest) {
// 	confirmBookingParam := confirmBookRequest.TeamAmericaConfirmBookingParam
// 	confirmBookingParam.BookingNumber = preBookingResponse.ParamPreBookResult.BookingNumber
// 	confirmBookingParam.Action = "AE"
// }

// func mapping_precancelresponse_to_hbe(
// 	preCancelResponse *ConfirmBookResponse,
// 	bookingCancelResponse *hbecommon.BookingCancelResponse) {

// 	if preCancelResponse.ConfirmBookResponseParam.Status == "00" {
// 		bookingCancelResponse.Status = "200"

// 		bookingCancelResponse.Ref = preCancelResponse.ConfirmBookResponseParam.BookingNumber
// 	} else {
// 		bookingCancelResponse.Status = "Failed"
// 		if preCancelResponse.ConfirmBookResponseParam.Error != nil {
// 			bookingCancelResponse.ErrorMessages = []*hbecommon.ErrorMessage{
// 				&hbecommon.ErrorMessage{
// 					Message: preCancelResponse.ConfirmBookResponseParam.Error.Description,
// 				}}
// 		}
// 	}
// }

func mapping_bookresponse_to_hbe(
	bookResponse *TeamAmericaBookReserveResponse,
	hbeBookingRequest *hbecommon.BookingRequest,
	bookingResponse *hbecommon.BookingResponse,
	cancelResponses []*TeamAmericaCancelPolicyResponse) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	if len(bookResponse.NewMultiItemReservationResponse.ReservationInformations) > 0 {
		bookItemResponse := bookResponse.NewMultiItemReservationResponse.ReservationInformations[0]
		if bookItemResponse.ReservationStatus == "OK" {
			bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum
			Ref := bookItemResponse.ReservationNumber
			SupplierReference, _ := json.Marshal(bookItemResponse.BookingAgentInfo)
			bookingResponse.Booking = &hbecommon.Booking{
				Ref:               Ref,
				ItineraryId:       bookItemResponse.ReservationNumber,
				SupplierReference: string(SupplierReference),
				SupplierName:      bookItemResponse.BookingAgentInfo.AgencyName,
				Total:             bookItemResponse.TotalResNetPrice}

			roomCxlResponse := From(cancelResponses).
				SingleWith(
					func(f interface{}) bool {
						return f.(*TeamAmericaCancelPolicyResponse).CancellationPolicyResponse.ProductCode == bookItemResponse.BookItems[0].ProductCode
					},
				).(*TeamAmericaCancelPolicyResponse)
			fmt.Printf("roomCxlResponse = %+v\n", roomCxlResponse)
			if roomCxlResponse != nil && len(roomCxlResponse.CancellationPolicyResponse.CancelPolicies) > 0 {
				bookingResponse.Booking.FreeCancellationPolicy = GetFreeCancelDate(
					roomCxlResponse.CancellationPolicyResponse.CancelPolicies,
					hbeBookingRequest.CheckIn)

				bookingResponse.Booking.CancellationPolicy = GenerateCancellationPolicyShort(
					roomCxlResponse.CancellationPolicyResponse.CancelPolicies,
					hbeBookingRequest.CheckIn,
					bookItemResponse.BookItems[0].AverageNetPricePerNight,
					bookItemResponse.TotalResNetPrice,
				)
			}
		} else {
			bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
			bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{Message: "Failed to booking reservation"}}
		}
	}
}

func mapping_cancelresponse_to_hbe(
	cancelResponse *TeamAmericaCancelReservationResponse,
	bookingCancelConfirmationResponse *hbecommon.BookingCancelResponse) {

	if cancelResponse != nil && cancelResponse.CancelReservationResp.ReservationStatusCode == "OK" {
		bookingCancelConfirmationResponse.Status = "200"
		ref, _ := json.Marshal(cancelResponse.CancelReservationResp.CancelItems)
		bookingCancelConfirmationResponse.Ref = string(ref)
	} else {
		bookingCancelConfirmationResponse.Status = "Failed"
	}
}

func mapping_roomcancelresponse_to_hbe(
	cancelResponse *TeamAmericaCancelPolicyResponse,
	hbeCancelResponse *hbecommon.RoomCxlResponse,
	arrivalDate string) {

	if cancelResponse != nil && cancelResponse.CancellationPolicyResponse.CancelPolicies != nil && len(cancelResponse.CancellationPolicyResponse.CancelPolicies) > 0 {
		hbeCancelResponse.Description = GetFreeCancelDate(cancelResponse.CancellationPolicyResponse.CancelPolicies, arrivalDate)
	} else {
		hbeCancelResponse.Description = "Failed to get room cancel policy"
	}
}

func GroupSearchRequestsByCity(searchRequests []*TeamAmericaHotelSearchRequest) []*TeamAmericaHotelSearchRequest {

	searchRequestsByCity := map[string]*TeamAmericaHotelSearchRequest{}

	for _, searchRequest := range searchRequests {

		key := strings.ToUpper(searchRequest.CityCode)
		if _, ok := searchRequestsByCity[key]; !ok {

			searchRequstByCity := searchRequest.Clone()
			searchRequstByCity.VendorIDs.VendorID = []string{}

			searchRequestsByCity[key] = searchRequstByCity
		}

		searchRequestsByCity[key].VendorIDs.VendorID = append(searchRequestsByCity[key].VendorIDs.VendorID, searchRequest.VendorIDs.VendorID[0])
	}

	results := []*TeamAmericaHotelSearchRequest{}

	for _, searchRequestByCity := range searchRequestsByCity {
		results = append(results, searchRequestByCity)
	}

	return results
}

func SplitBatchSearchRequests(searchRequests []*TeamAmericaHotelSearchRequest, maxBatchSize int) []*TeamAmericaHotelSearchRequest {

	batches := []*TeamAmericaHotelSearchRequest{}

	for _, searchRequest := range searchRequests {

		batches = append(batches, SplitBatchSearchRequest(searchRequest, maxBatchSize)...)
	}

	return batches
}

func SplitBatchSearchRequest(searchRequest *TeamAmericaHotelSearchRequest, maxBatchSize int) []*TeamAmericaHotelSearchRequest {

	batches := []*TeamAmericaHotelSearchRequest{}

	var counter int = 0
	for counter < len(searchRequest.VendorIDs.VendorID) {

		batch := searchRequest.Clone()

		var batchSize int = maxBatchSize
		if counter+batchSize > len(searchRequest.VendorIDs.VendorID) {

			batchSize = len(searchRequest.VendorIDs.VendorID) - counter
		}

		batch.VendorIDs.VendorID = searchRequest.VendorIDs.VendorID[counter : counter+batchSize]

		batches = append(batches, batch)

		counter += batchSize
	}

	return batches
}
