package hb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../roomres/utils"

	hbecommon "../roomres/hbe/common"
	//repository "../roomres/repository"
)

func ProduceRawHotelId(combinedHotelId string) string {

	if strings.Index(combinedHotelId, "-") >= 0 {

		values := strings.Split(combinedHotelId, "-")

		return values[1]
	}

	return combinedHotelId
}

func mapping_hbe_to_bookingreportrequest(request *hbecommon.BookingReportRequest, reportRequest *BookingReportRequest) {

	reportRequest.DateFrom = request.BookingDateFrom
	reportRequest.DateTo = request.BookingDateTo

	reportRequest.From = request.BookingIndexFrom
	reportRequest.To = request.BookingIndexFrom + 10
}

func mapping_bookingreportresponse_to_hbe(responseReport *BookingReportResponse, response *hbecommon.BookingReportResponse, destinations map[string]*Destination) {

	if responseReport.Bookings == nil {
		return
	}

	response.Bookings = []*hbecommon.BookingReport{}
	response.MoreBookings = len(responseReport.Bookings.Items) > 0

	for _, booking := range responseReport.Bookings.Items {

		hbeBooking := &hbecommon.BookingReport{}

		hbeBooking.Ref = booking.Reference
		hbeBooking.InternalRef = booking.ClientReference

		hbeBooking.CheckIn, _ = time.Parse(utils.LayoutYYYYMMDD, booking.Hotel.CheckIn)
		hbeBooking.CheckOut, _ = time.Parse(utils.LayoutYYYYMMDD, booking.Hotel.CheckOut)

		hbeBooking.HotelSupplierRef = fmt.Sprintf("%s-%s", booking.Hotel.DestinationCode, booking.Hotel.Code)
		hbeBooking.HotelName = booking.Hotel.Name

		if destName, ok := destinations[strings.ToUpper(booking.Hotel.DestinationCode)]; ok {
			hbeBooking.CityName = destName.Name.Content
		} else {
			hbeBooking.CityName = booking.Hotel.DestinationCode
		}

		hbeBooking.Total = booking.TotalNet

		if booking.Hotel.IsConfirmed() {
			hbeBooking.Status = hbecommon.BookingReportStatusConfirmed
		} else {
			hbeBooking.Status = hbecommon.BookingReportStatusNotConfirmed
		}

		response.Bookings = append(response.Bookings, hbeBooking)
	}

}

// func mapping_hbe_to_hotelcontentrequest(request *hbecommon.BookingInfoRequest, hotelMappings []*repository.HotelMapping, hotelContentRequest *HotelContentRequest) {

// 	hotelMapping := hotelMappings[0]

// 	hotelContentRequest.HotelId = ProduceRawHotelId(hotelMapping.SupplierId)
// }

func mapping_hotelcontentresponse_to_hbe(hotelContentResponse *HotelContentResponse, response *hbecommon.BookingInfoResponse) {

	if hotelContentResponse.HotelContent == nil {
		return
	}

	for _, phone := range hotelContentResponse.HotelContent.Phones {

		if strings.ToUpper(phone.PhoneType) == "PHONEHOTEL" {

			response.HotelPhoneNumber = phone.PhoneNumber

			break
		}
	}

}

func mapping_hbe_to_searchrequest(
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	availabilityRequest *AvailabilityRequest) {

	availabilityRequest.Hotels = &Hotels{}
	for _, hotelId := range hotelSearchRequest.HotelIds {

		value := hotelId

		value = ProduceRawHotelId(value)

		/*
			if strings.Index(value, "-") >= 0 {

				values := strings.Split(value, "-")
				value = values[1]
			}
		*/

		id, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}

		availabilityRequest.Hotels.HotelIds = append(availabilityRequest.Hotels.HotelIds, id)
	}

	availabilityRequest.Stay = &Stay{
		CheckIn:  hotelSearchRequest.CheckIn,
		CheckOut: hotelSearchRequest.CheckOut,
	}

	availabilityRequest.Occupancies = []*Occupancy{}
	for _, requestedRoom := range hotelSearchRequest.RequestedRooms {

		occupancy := &Occupancy{}

		mapping_hbe_to_occupacy(requestedRoom, occupancy)

		availabilityRequest.Occupancies = append(availabilityRequest.Occupancies, occupancy)
	}

	availabilityRequest.Filter = &Filter{}

	if len(availabilityRequest.Hotels.HotelIds) != 1 && len(availabilityRequest.Occupancies) == 1 {

		availabilityRequest.Filter.MaxRooms = 1
		availabilityRequest.Filter.MaxRatesPerRoom = 1
	}

	availabilityRequest.Filter.PaymentType = "AT_WEB"
	availabilityRequest.Filter.HotelPackage = "NO"
	availabilityRequest.Filter.Packaging = hotelSearchRequest.Packaging
}

func mapping_hbe_to_occupacy(roomRequested *hbecommon.RoomRequest, occupancy *Occupancy) {

	occupancy.Adults = roomRequested.Adults
	occupancy.Children = roomRequested.Children
	occupancy.Rooms = 1
	occupancy.Paxes = []*Pax{}

	for i := 0; i < occupancy.Adults; i++ {

		occupancy.Paxes = append(occupancy.Paxes, &Pax{Type: "AD"})
	}

	for _, childAge := range roomRequested.ChildAges {
		pax := &Pax{}
		mapping_hbe_to_pax(childAge, pax)

		occupancy.Paxes = append(occupancy.Paxes, pax)
	}
}

func mapping_hbe_to_pax(child *hbecommon.ChildAge, pax *Pax) {

	pax.Type = "CH"
	pax.Age = child.Age
}

func mapping_searchresponse_to_hbe(
	availabilityResponse *AvailabilityResponse,
	hotelSearchRequest *hbecommon.HotelSearchRequest,
	hbeHotelSearchResponse *hbecommon.HotelSearchResponse) {

	if availabilityResponse.Hotels == nil {
		return
	}

	hbeHotelSearchResponse.Hotels = []*hbecommon.Hotel{}
	for _, hotel := range availabilityResponse.Hotels.Hotels {

		hbeHotel := &hbecommon.Hotel{HotelId: fmt.Sprintf("%s-%d", hotel.DestinationCode, hotel.HotelId)}

		mapping_hotel_to_hbe(hotelSearchRequest, hotel, hbeHotel)

		if hbeHotel.CheapestRoom != nil {
			hbeHotelSearchResponse.Hotels = append(hbeHotelSearchResponse.Hotels, hbeHotel)
		}
	}
}

func mapping_hotel_to_hbe(hotelSearchRequest *hbecommon.HotelSearchRequest, hotel *Hotel, hbeHotel *hbecommon.Hotel) {

	bestRoomRate := FindBestRoomRate(FindBestRoomRateCombinations(hotelSearchRequest.RequestedRooms, hotel.Rooms))

	hbeRoomTypes := []*hbecommon.RoomType{}

	hbeRoomType := &hbecommon.RoomType{}

	if bestRoomRate != nil {

		mapping_roomtype_to_hbe(
			hotelSearchRequest.GetLos(),
			bestRoomRate.Room,
			bestRoomRate.Rate,
			hbeRoomType,
		)

		hbeRoomTypes = append(hbeRoomTypes, hbeRoomType)
	}

	if len(hbeRoomTypes) > 0 {
		cheapestRoomType := hbeRoomTypes[0]

		hbeHotel.RoomTypes = []*hbecommon.RoomType{cheapestRoomType}
		hbeHotel.CheapestRoom = cheapestRoomType
	}
}

func mapping_roomtype_to_hbe(los int, room *Room, rate *Rate, hbeRoomType *hbecommon.RoomType) {

	hbeRoomType.ShortRef = room.Code
	hbeRoomType.Ref = room.Code

	hbeRoomType.Description = room.Name

	hbeRoomType.Taxes = []*hbecommon.Tax{}
	hbeRoomType.Surcharges = []*hbecommon.Surcharge{}

	hbeRoomType.Rate = &hbecommon.Rate{
		PerNight:     rate.NetValue / float32(los),
		PerNightBase: rate.NetValue / float32(los),
		Total:        rate.NetValue,
	}

}
