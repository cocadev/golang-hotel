package restel

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"strings"

	hbecommon "../roomres/hbe/common"

	. "github.com/ahmetb/go-linq"
)

func mapping_hbe_to_searchrequest(hotelSearchRequest *hbecommon.HotelSearchRequest) []*SearchRequest {
	var searchRequests []*SearchRequest
	idx := 0
	for idx = 0; idx < len(hotelSearchRequest.HotelIds); idx++ {
		searchRequests = append(searchRequests, &SearchRequest{
			RestelSearchParam: HbeToRestelHotelSearchRequest(hotelSearchRequest, idx),
		})
	}

	//fmt.Printf("multiple search requests count : %d\n", len(searchRequests))
	return searchRequests
}

func RestelSearchRequestToHbeSearchRequest(restelHotelSearchRequest *RestelHotelSearchRequest) *hbecommon.HotelSearchRequest {
	tempCheckIn, _ := time.Parse("01/02/2006", restelHotelSearchRequest.CheckInDate)
	tempCheckOut, _ := time.Parse("01/02/2006", restelHotelSearchRequest.CheckOutDate)

	var hotelIds []string
	hotelIds = append(hotelIds, restelHotelSearchRequest.Country)
	hotelIds = append(hotelIds, restelHotelSearchRequest.Province)
	hotelIds = append(hotelIds, restelHotelSearchRequest.Affiliation)
	hotelIds = append(hotelIds, restelHotelSearchRequest.UserCode)
	hotelIds = append(hotelIds, restelHotelSearchRequest.Hotel)

	var externalRefs []string
	externalRefs = append(externalRefs, ConvertIntToString(restelHotelSearchRequest.Category))
	externalRefs = append(externalRefs, ConvertIntToString(restelHotelSearchRequest.Radius))
	externalRefs = append(externalRefs, ConvertIntToString(restelHotelSearchRequest.Language))
	externalRefs = append(externalRefs, ConvertIntToString(restelHotelSearchRequest.CompressedXML))
	externalRefs = append(externalRefs, ConvertIntToString(restelHotelSearchRequest.Information))

	var roomRequests []*hbecommon.RoomRequest
	if restelHotelSearchRequest.Type1RoomNumber > 0 {
		tempStrs := strings.Split(restelHotelSearchRequest.Type1Guest, "-")
		childAgeStrs := strings.Split(restelHotelSearchRequest.Type1ChildAges, ",")
		var ChildAges []*hbecommon.ChildAge
		for _, ele := range childAgeStrs {
			ChildAges = append(ChildAges, &hbecommon.ChildAge{Age: ConvertStringToInt(ele)})
		}
		roomRequests = append(roomRequests, &hbecommon.RoomRequest{
			Adults:    ConvertStringToInt(tempStrs[0]),
			Children:  ConvertStringToInt(tempStrs[1]),
			ChildAges: ChildAges,
		})
	}

	if restelHotelSearchRequest.Type2RoomNumber > 0 {
		tempStrs := strings.Split(restelHotelSearchRequest.Type2Guest, "-")
		childAgeStrs := strings.Split(restelHotelSearchRequest.Type2ChildAges, ",")
		var ChildAges []*hbecommon.ChildAge
		for _, ele := range childAgeStrs {
			ChildAges = append(ChildAges, &hbecommon.ChildAge{Age: ConvertStringToInt(ele)})
		}
		roomRequests = append(roomRequests, &hbecommon.RoomRequest{
			Adults:    ConvertStringToInt(tempStrs[0]),
			Children:  ConvertStringToInt(tempStrs[1]),
			ChildAges: ChildAges,
		})
	}

	if restelHotelSearchRequest.Type3RoomNumber > 0 {
		tempStrs := strings.Split(restelHotelSearchRequest.Type3Guest, "-")
		childAgeStrs := strings.Split(restelHotelSearchRequest.Type3ChildAges, ",")
		var ChildAges []*hbecommon.ChildAge
		for _, ele := range childAgeStrs {
			ChildAges = append(ChildAges, &hbecommon.ChildAge{Age: ConvertStringToInt(ele)})
		}
		roomRequests = append(roomRequests, &hbecommon.RoomRequest{
			Adults:    ConvertStringToInt(tempStrs[0]),
			Children:  ConvertStringToInt(tempStrs[1]),
			ChildAges: ChildAges,
		})
	}

	return &hbecommon.HotelSearchRequest{
		//SessionId string `json:"SessionId"`
		//IpAddress string `json:"IpAddress"`
		//UserAgent string `json:"UserAgent"`

		//AgencyId   int `json:"AgencyId"`
		//ProviderId int `json:"ProviderId"`
		CheckIn:        tempCheckIn.Format("2006-01-02"),
		CheckOut:       tempCheckOut.Format("2006-01-02"),
		RequestedRooms: roomRequests,
		//Lat              float32        `json:"Lat"`
		//Lon              float32        `json:"Lon"`
		HotelIds:     []string{strings.Join(hotelIds, "|")},
		ExternalRefs: externalRefs,
		//ExternalRefs     []string       `json:"ExternalRefs"`
		//CurrencyCode     string         `json:"CurrencyCode"`
		//SpecificRoomRefs []string       `json:"SpecificRoomRefs"`
		//Details          bool           `json:"Details"`
		//Packaging        bool           `json:"Packaging"`

		//SortType string `json:"SortType"`

		//AutocompleteId     string `json:"AutocompleteId"`
		//HotelFilterId      string `json:"HotelFilterId"`
		//RecommendationOnly bool   `json:"RecommendationOnly"`

		SingleHotelSearch: false,
		SingleHotelFilter: false,
		//SingleHotelId     int  `json:"-"`

		//MaxPrice    float32   `json:"MaxPrice"`
		//MinPrice    float32   `json:"MinPrice"`
		//StarRatings []float32 `json:"StarRatings"`

		PageIndex: 1,
		PageSize:  10,
	}
}

func HbeToRestelHotelSearchRequest(hbeSearchRequest *hbecommon.HotelSearchRequest, idx int) *RestelHotelSearchRequest {
	tempCheckIn, _ := time.Parse("2006-01-02", hbeSearchRequest.CheckIn)
	tempCheckOut, _ := time.Parse("2006-01-02", hbeSearchRequest.CheckOut)

	hotelIds := hbeSearchRequest.HotelIds
	hotelIdIndexes := strings.Split(GetRawHotelId(hotelIds[idx]), "|")

	ret := &RestelHotelSearchRequest{
		Hotel:        hotelIdIndexes[0] + "#",
		Country:      hotelIdIndexes[1],
		Province:     hotelIdIndexes[2],
		Poblacion:    "",
		Category:     0, //ConvertStringToInt(hotelIdIndexes[5]),
		Radius:       9, //ConvertStringToInt(hotelIdIndexes[6]),
		CheckInDate:  tempCheckIn.Format("01/02/2006"),
		CheckOutDate: tempCheckOut.Format("01/02/2006"),
		//Group           string `xml:"marca"`
		Type1RoomNumber: 0,
		Type1Guest:      "0",
		//Type1ChildAges  string `xml:"edades1"`
		Type2RoomNumber: 0,
		Type2Guest:      "0",
		//Type2ChildAges  string `xml:"edades2"`
		Type3RoomNumber: 0,
		Type3Guest:      "0",
		//Type3ChildAges  string `xml:"edades3"`
		Language:        0,
		Duplicated:      1, // int    `xml:"duplicidad"`
		CompressedXML:   2,
		Information:     0,
		Refundable:      0, //ConvertStringToInt(hotelIdIndexes[7]),
		CompoundHotelId: GetRawHotelId(hotelIds[idx]),
	}

	roomRequests := hbeSearchRequest.RequestedRooms

	//fmt.Printf("Count of room request : %d\n", len(roomRequests))

	if len(roomRequests) >= 1 {
		roomRequest := roomRequests[0]
		ret.Type1RoomNumber = 1
		ret.Type1Guest = ConvertIntToString(roomRequest.Adults) + "-" + ConvertIntToString(roomRequest.Children)
		if roomRequest.ChildAges != nil {
			var childAges []string

			for _, chlidAge := range roomRequest.ChildAges {
				childAges = append(childAges, ConvertIntToString(chlidAge.Age))
			}

			ret.Type1ChildAges = strings.Join(childAges, ",")
		}
	}

	if len(roomRequests) >= 2 {
		roomRequest := roomRequests[1]
		ret.Type2RoomNumber = 1
		ret.Type2Guest = ConvertIntToString(roomRequest.Adults) + "-" + ConvertIntToString(roomRequest.Children)
		if roomRequest.ChildAges != nil {
			var childAges []string

			for _, chlidAge := range roomRequest.ChildAges {
				childAges = append(childAges, ConvertIntToString(chlidAge.Age))
			}

			ret.Type2ChildAges = strings.Join(childAges, ",")
		}
	}

	if len(roomRequests) >= 3 {
		roomRequest := roomRequests[2]
		ret.Type3RoomNumber = 1
		ret.Type3Guest = ConvertIntToString(roomRequest.Adults) + "-" + ConvertIntToString(roomRequest.Children)
		if roomRequest.ChildAges != nil {
			var childAges []string

			for _, chlidAge := range roomRequest.ChildAges {
				childAges = append(childAges, ConvertIntToString(chlidAge.Age))
			}

			ret.Type3ChildAges = strings.Join(childAges, ",")
		}
	}

	return ret
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
	searchResponse *SearchResponse,
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse,
	cancelFeeResponses []*CancelFeeDetailResponse,
	//logScope roomresutils.ILog
) {
	roomTypes := 0

	for _, hotel := range searchResponse.ParamResult.HotelResults.Hotels {

		hbeHotel := &hbecommon.Hotel{}

		mapping_hotel_to_hbe(
			hbeHotelSearchRequest,
			hotel,
			hbeHotel,
			cancelFeeResponses,
			//logScope,
		)

		if hbeHotel.CheapestRoom != nil {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
			roomTypes += len(hbeHotel.RoomTypes)
		}
	}

	fmt.Printf("Result : %d Hotel(s), %d Room Type(s)\n", len(hbeHotelSearchResponse.Hotels), roomTypes)
}

func mapping_hotel_to_hbe(
	hbeHotelSearchRequest *hbecommon.HotelSearchRequest,
	hotelResult *HotelResult,
	hbeHotel *hbecommon.Hotel,
	roomCxlResponses []*CancelFeeDetailResponse,
	//logScope roomresutils.ILog
) {
	hbeHotel.HotelId = hotelResult.HotelCode
	hbeHotel.CustomTag = hotelResult.HotelAffiliation

	hbeRoomTypes := []*hbecommon.RoomType{}
	checkOut, _ := time.Parse("2006-01-02", hbeHotelSearchRequest.CheckOut)
	checkIn, _ := time.Parse("2006-01-02", hbeHotelSearchRequest.CheckIn)

	deltaDays := float32(math.Floor(checkOut.Sub(checkIn).Hours() / 24))

	var searchRoomRefs []*ShortRoomRef

	if hbeHotelSearchRequest.Details && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
		for _, roomRefStr := range hbeHotelSearchRequest.SpecificRoomRefs {
			var shortRoomRef ShortRoomRef
			err := json.Unmarshal([]byte(roomRefStr), &shortRoomRef)
			if err != nil {
				fmt.Printf("Cannot parse specialRoomRef = %s\n", roomRefStr)
			} else {
				searchRoomRefs = append(searchRoomRefs, &shortRoomRef)
			}
		}
	}

	for _, roomPlan := range hotelResult.HotelRestrict.RoomPlans {
		//removed by li, 2018/10/20
		// if checkDetail && len(searchRoomRefs) > 0 {
		// 	//should check specifiedRoomRefs
		// 	isMatch := false
		// 	for _, searchRoomRef := range searchRoomRefs {
		// 		if searchRoomRef.AdultChildren != "" && hotelResult.HotelRestrict.AdultChldren != searchRoomRef.AdultChildren {
		// 			continue
		// 		}
		// 		if searchRoomRef.RoomType != "" && roomPlan.RoomType != searchRoomRef.RoomType {
		// 			continue
		// 		}
		// 		if searchRoomRef.MealPlanType != "" && roomPlan.RoomLine.MealPlanType != searchRoomRef.MealPlanType {
		// 			continue
		// 		}
		// 		if searchRoomRef.RefundableType != "" && roomPlan.RoomLine.FeeRate != searchRoomRef.RefundableType {
		// 			continue
		// 		}

		// 		if searchRoomRef.MinPrice != "" && ConvertStringToFloat32(roomPlan.RoomLine.PlanPrice) < ConvertStringToFloat32(searchRoomRef.MinPrice) {
		// 			continue
		// 		}

		// 		if searchRoomRef.MaxPrice != "" && ConvertStringToFloat32(roomPlan.RoomLine.PlanPrice) > ConvertStringToFloat32(searchRoomRef.MaxPrice) {
		// 			continue
		// 		}
		// 		isMatch = true
		// 		break
		// 	}

		// 	if !isMatch {
		// 		continue
		// 	}
		// }

		Ref, _ := json.Marshal(roomPlan)
		total := ConvertStringToFloat32(roomPlan.RoomLine.PlanPrice)
		perNight := total / deltaDays
		roomType := &hbecommon.RoomType{
			Rate: &hbecommon.Rate{
				PerNight:     perNight,
				PerNightBase: perNight,
				Total:        total,
				CurrecyCode:  roomPlan.RoomLine.Currency,
			},
			Description: roomPlan.RoomDesc,
			Nights:      int(deltaDays),
			Ref:         string(Ref),
		}

		if hbeHotelSearchRequest.Details && len(hbeHotelSearchRequest.SpecificRoomRefs) > 0 {
			//set cancellation policy here
			roomCxlResponse := From(roomCxlResponses).
				SingleWith(
					func(f interface{}) bool {
						hashText := GetStringArrayHash(roomPlan.RoomLine.CompressedAvailability)
						return f.(*CancelFeeDetailResponse).CancelFeeDetailResponseParam.RequestKey == hashText
					},
				).(*CancelFeeDetailResponse)
			if roomCxlResponse != nil && len(roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies) > 0 {
				roomType.FreeCancellationPolicy = GetFreeCancelDate(
					roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies,
					hbeHotelSearchRequest.CheckIn,
				)
				roomType.CancellationPolicy = GenerateCancellationPolicyShort(
					roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies,
					hbeHotelSearchRequest.CheckIn, perNight, total)
			} else {
				roomType.NonRefundable = true
			}
		}

		hbeRoomTypes = append(hbeRoomTypes, roomType)
	}

	if len(hbeRoomTypes) == 0 {
		return
	}

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
}

func mapping_hbe_prebookrequest(
	hbeBookingRequest *hbecommon.BookingRequest,
	index int) *RestelBookingRequest {

	//tempCheckIn, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckIn)
	//tempCheckOut, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckOut)

	hotelIdIndexes := strings.Split(GetRawHotelId(hbeBookingRequest.Hotel.HotelId), "|")
	room := hbeBookingRequest.Hotel.Rooms[index]
	creditCardInfo := hbeBookingRequest.CreditCardInfo
	customerInfo := hbeBookingRequest.Customer

	var roomRef RoomPlan
	err := json.Unmarshal([]byte(room.Ref), &roomRef)
	if err != nil {
		fmt.Errorf("Cannot parse room ref field : %s, err=%s\n", room.Ref, err)
		return nil
	}

	var internalRef InternalRef
	err = json.Unmarshal([]byte(hbeBookingRequest.InternalRef), &internalRef)
	if err != nil {
		fmt.Printf("Cannot parse InternalRef field : %s, err=%s\n", hbeBookingRequest.InternalRef, err)
		return nil
	}

	var ccExpireDate time.Time
	if creditCardInfo.ExpiryDate != "" {
		ccExpireDate, err = time.Parse("20060102", creditCardInfo.ExpiryDate)
		if err != nil {
			fmt.Printf("Cannot parse credit card expire date : %s, err=%s\n", creditCardInfo.ExpiryDate, err)
			return nil
		}
	}

	var primaryGuest *hbecommon.Guest = nil
	for _, guest := range room.Guests {
		if guest.Primary {
			primaryGuest = guest
			break
		}
	}

	if primaryGuest == nil {
		return nil
	}

	restelBookingParam := &RestelBookingRequest{
		HotelCode:         hotelIdIndexes[0],
		GuestName:         primaryGuest.FirstName + " " + primaryGuest.LastName,
		Remarks:           internalRef.Remark,
		InternalUse:       internalRef.InternalUse,
		PaymentForm:       internalRef.PaymentForm,
		GuestContactEmail: customerInfo.Email,
		GuestContactPhone: customerInfo.PhoneNumber,
		SimpleLines: &SimpleLines{
			CompressedAvailability: roomRef.RoomLine.CompressedAvailability,
			OccupantsNames:         makeOccupantsNames(room),
		},
	}

	fmt.Printf("Credit card information %+v, %s\n", creditCardInfo, creditCardInfo.CardType)

	if creditCardInfo.CardType != "" {
		//TODO : need to validate creditcard type etc ...

		restelBookingParam.CreditCardType = creditCardInfo.CardType
		restelBookingParam.CreditCadNumber = creditCardInfo.Number
		restelBookingParam.CVV = creditCardInfo.Cvc
		restelBookingParam.MonthOfExpire = ccExpireDate.Format("01")
		restelBookingParam.YearOfExpire = ccExpireDate.Format("2006")
		restelBookingParam.HolderName = creditCardInfo.HolderName
	}

	return restelBookingParam
}

func makeOccupantsNames(room *hbecommon.BookingRoom) string {
	var guestNames []string //name1#surname1#.#age1#
	for _, guest := range room.Guests {
		guestNames = append(guestNames, guest.LastName+"#"+guest.FirstName+"#.#"+ConvertIntToString(guest.Age))
	}

	return strings.Join(guestNames, "@") + "@"
}

func mapping_hbe_bookingrequest_to_prebookrequest(
	hbeBookingRequest *hbecommon.BookingRequest) []*PreBookRequest {
	preBookRequests := []*PreBookRequest{}
	idx := 0
	fmt.Printf("Room counts = %d\n", len(hbeBookingRequest.Hotel.Rooms))
	for idx = 0; idx < len(hbeBookingRequest.Hotel.Rooms); idx++ {
		preBookRequests = append(preBookRequests, &PreBookRequest{
			RestelBookingParam: mapping_hbe_prebookrequest(hbeBookingRequest, idx),
		})
	}

	return preBookRequests
}

func makeBookingConfirmRequest(
	hbeBookingRequest *hbecommon.BookingRequest,
	preBookingResponse *PreBookResponse,
	confirmBookRequest *ConfirmBookRequest) {
	confirmBookingParam := confirmBookRequest.RestelConfirmBookingParam
	confirmBookingParam.BookingNumber = preBookingResponse.ParamPreBookResult.BookingNumber
	confirmBookingParam.Action = "AE"
}

func mapping_precancelresponse_to_hbe(
	preCancelResponse *ConfirmBookResponse,
	bookingCancelResponse *hbecommon.BookingCancelResponse) {

	if preCancelResponse.ConfirmBookResponseParam.Status == "00" {
		bookingCancelResponse.Status = "200"

		bookingCancelResponse.Ref = preCancelResponse.ConfirmBookResponseParam.BookingNumber
	} else {
		bookingCancelResponse.Status = "Failed"
		if preCancelResponse.ConfirmBookResponseParam.Error != nil {
			bookingCancelResponse.ErrorMessages = []*hbecommon.ErrorMessage{
				&hbecommon.ErrorMessage{
					Message: preCancelResponse.ConfirmBookResponseParam.Error.Description,
				}}
		}
	}
}

func mapping_bookresponse_to_hbe(
	bookResponse *ConfirmBookResponse,
	preBookResponse *PreBookResponse,
	bookingResponse *hbecommon.BookingResponse,
	roomCxlResponses []*CancelFeeDetailResponse,
	hbeBookingRequest *hbecommon.BookingRequest) {

	const (
		BookingFailureExceptionType1 string = "Booking failed (Component Failure)"
		BookingFailureExceptionType2 string = "Failed to book third party component"
	)

	if bookResponse.ConfirmBookResponseParam.Status == "00" {

		bookingResponse.BookingStatus = hbecommon.BookingStatusConfirmedEnum

		bookingResponse.Booking = &hbecommon.Booking{
			Ref:               preBookResponse.ParamPreBookResult.BookingNumber,
			ItineraryId:       bookResponse.ConfirmBookResponseParam.ShortBookingNumber,
			SupplierReference: preBookResponse.ParamPreBookResult.Remarks,
			Total:             preBookResponse.ParamPreBookResult.TotalReservationAmount}

		if len(roomCxlResponses) > 0 {
			//set cancellation policy here
			roomCxlResponse := roomCxlResponses[0]
			if roomCxlResponse != nil && len(roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies) > 0 {
				bookingResponse.Booking.FreeCancellationPolicy = GetFreeCancelDate(
					roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies,
					hbeBookingRequest.CheckIn,
				)

				total := preBookResponse.ParamPreBookResult.TotalReservationAmount
				checkOut, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckOut)
				checkIn, _ := time.Parse("2006-01-02", hbeBookingRequest.CheckIn)

				deltaDays := float32(math.Floor(checkOut.Sub(checkIn).Hours() / 24))
				bookingResponse.Booking.CancellationPolicy = GenerateCancellationPolicyShort(
					roomCxlResponse.CancelFeeDetailResponseParam.CancelPolicies,
					hbeBookingRequest.CheckIn, total/deltaDays,
					preBookResponse.ParamPreBookResult.TotalReservationAmount)
			}
		}
	} else {
		bookingResponse.BookingStatus = hbecommon.BookingStatusFailedEnum
		bookingResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: "Failed to confirm book number(" + preBookResponse.ParamPreBookResult.BookingNumber + ")"}}

	}
}

func mapping_cancelresponse_to_hbe(
	cancelResponse *CancelConfirmedBookResponse,
	bookingCancelConfirmationResponse *hbecommon.BookingCancelConfirmationResponse) {

	if cancelResponse.CancelConfirmedBookResponseParam.Status == "00" {
		bookingCancelConfirmationResponse.Status = "200"
	} else {
		bookingCancelConfirmationResponse.Status = "Failed"
		bookingCancelConfirmationResponse.ErrorMessages = []*hbecommon.ErrorMessage{
			&hbecommon.ErrorMessage{Message: cancelResponse.CancelConfirmedBookResponseParam.Comment}}
	}
}

func mapping_cancelfeedetail_to_hbe(
	response *CancelFeeDetailResponse,
	hbeResponse *hbecommon.RoomCxlResponse,
	arrivalDate string) {
	hbeResponse.Description = GetFreeCancelDate(response.CancelFeeDetailResponseParam.CancelPolicies, arrivalDate)
}

func GroupSearchRequestsByCity(searchRequests []*SearchRequest) []*SearchRequest {

	searchRequestsByCity := map[string]*SearchRequest{}

	for _, searchRequest := range searchRequests {

		key := searchRequest.RestelSearchParam.Province
		if _, ok := searchRequestsByCity[key]; !ok {

			searchRequstByCity := searchRequest.Clone()
			searchRequstByCity.RestelSearchParam.Hotel = ""

			searchRequestsByCity[key] = searchRequstByCity
		}

		searchRequestsByCity[key].RestelSearchParam.Hotel += searchRequest.RestelSearchParam.Hotel
	}

	results := []*SearchRequest{}

	for _, searchRequestByCity := range searchRequestsByCity {
		results = append(results, searchRequestByCity)
	}

	return results
}

func SplitBatchSearchRequests(searchRequests []*SearchRequest, maxBatchSize int) []*SearchRequest {

	batches := []*SearchRequest{}

	for _, searchRequest := range searchRequests {

		batches = append(batches, SplitBatchSearchRequest(searchRequest, maxBatchSize)...)
	}

	return batches
}

func SplitBatchSearchRequest(searchRequest *SearchRequest, maxBatchSize int) []*SearchRequest {

	batches := []*SearchRequest{}

	var counter int = 0
	var hotels = delete_empty(strings.Split(searchRequest.RestelSearchParam.Hotel, "#"))
	for counter < len(hotels) {

		batch := searchRequest.Clone()

		var batchSize int = maxBatchSize
		if counter+batchSize > len(hotels) {

			batchSize = len(hotels) - counter
		}

		batchHotels := hotels[counter : counter+batchSize]
		batch.RestelSearchParam.Hotel = strings.Join(batchHotels, "#") + "#"

		batches = append(batches, batch)

		counter += batchSize
	}

	return batches
}
