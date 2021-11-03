package hb

import (
	//repository "../roomres/repository"

	"fmt"

	roomresutils "../roomres/utils"

	hbecommon "../roomres/hbe/common"
	//. "github.com/ahmetb/go-linq"
)

type HBClient struct {
	ProviderId      int
	HotelsPerBatch  int
	NumberOfThreads int

	ApiKey    string
	SecretKey string

	SearchEndPoint        string
	CheckRateEndPoint     string //added by Li, 20180917
	BookingEndPoint       string
	BookingDetailEndPoint string
	ContentEndPoint       string

	CommissionProvider hbecommon.ICommissionProvider
	//RepositoryFactory  *repository.RepositoryFactory

	Log roomresutils.ILog

	Destinations map[string]*Destination
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings, logging roomresutils.ILog) *HBClient {

	hbClient := &HBClient{

		Log: logging,

		SearchEndPoint:        hotelProviderSettings.SearchEndPoint,
		CheckRateEndPoint:     hotelProviderSettings.BookingSpecialRequestEndPoint,
		BookingEndPoint:       hotelProviderSettings.BookingConfirmationEndPoint,
		ContentEndPoint:       hotelProviderSettings.PropertyDetailsEndpoint,
		BookingDetailEndPoint: hotelProviderSettings.BookingDetailsEndpoint,

		ProviderId:         hotelProviderSettings.ProviderId,
		ApiKey:             hotelProviderSettings.ApiKey,
		SecretKey:          hotelProviderSettings.SecretKey,
		HotelsPerBatch:     hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads:    hotelProviderSettings.NumberOfThreads,
		CommissionProvider: hotelProviderSettings.CommissionProvider,
		//RepositoryFactory:  hotelProviderSettings.RepositoryFactory,
	}

	return hbClient
}

func (m *HBClient) GetProviderName() string {

	return "HotelBeds"
}

func (m *HBClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"Api-Key":      m.ApiKey,
		"X-Signature":  GenerateSignature(m.ApiKey, m.SecretKey),
	}
}

func (m *HBClient) NewSearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *AvailabilityRequest {

	return &AvailabilityRequest{}
}

func (m *HBClient) NewCheckRateRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *CheckRateRequest {

	return &CheckRateRequest{}
}

func CheckSearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) bool {

	return hotelSearchRequest.GetLos() == 1
}

func (m *HBClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {

	if hotelSearchRequest.Details {

		// HU B code here
		//check hotel counts
		if len(hotelSearchRequest.HotelIds) > 1 {
			return nil
		}

		hotelSearchResponse := &hbecommon.HotelSearchResponse{}
		if hotelSearchRequest.SpecificRoomRefs != nil && len(hotelSearchRequest.SpecificRoomRefs) > 0 {
			request := m.NewCheckRateRequest(hotelSearchRequest)
			mapping_hbe_to_checkraterequest(hotelSearchRequest, request)
			response := m.SendCheckRateRequest(request)
			mapping_searchresponsedetail_to_hbe(response, hotelSearchRequest, hotelSearchResponse)
		} else {
			request := m.NewSearchRequest(hotelSearchRequest)
			mapping_hbe_to_searchrequest(hotelSearchRequest, request)
			response := m.SendSearchRequest(request)
			mapping_searchresponsedetail_to_hbe(response, hotelSearchRequest, hotelSearchResponse)
		}

		return hotelSearchResponse
	} else {
		//if CheckSearchRequest(hotelSearchRequest) {

		request := m.NewSearchRequest(hotelSearchRequest)
		mapping_hbe_to_searchrequest(hotelSearchRequest, request)

		requests := m.SplitSearchRequests(request)
		response := m.SendSearchRequests(requests)

		hotelSearchResponse := &hbecommon.HotelSearchResponse{}
		mapping_searchresponse_to_hbe(response, hotelSearchRequest, hotelSearchResponse)

		return hotelSearchResponse

		//} else {
		//	return &hbecommon.HotelSearchResponse{}
		//}
	}
}

func (m *HBClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {
	// Hu B code here
	cancelResponse := m.SendCancelRequest(bookingCancelRequest.Ref)

	hbeCancelResponse := &hbecommon.BookingCancelResponse{}

	fmt.Printf("cancelResponse = %+v\n", cancelResponse)
	mapping_cancelresponse_to_hbe(cancelResponse, bookingCancelRequest, hbeCancelResponse)

	return hbeCancelResponse
}

func (m *HBClient) CancelBookingConfirm(bookingCancelConfirmationRequest *hbecommon.BookingCancelConfirmationRequest) *hbecommon.BookingCancelConfirmationResponse {
	return nil
}

func (m *HBClient) SendSpecialRequest(bookingSpecialRequestRequest *hbecommon.BookingSpecialRequestRequest) *hbecommon.BookingSpecialRequestResponse {
	return nil
}

func (m *HBClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {
	bookRequest := &BookingRequest{}
	makeBookingRequest(hbeBookingRequest, bookRequest)

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse)

	return hbeBookingResponse
}

func (m *HBClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {
	return nil
}
