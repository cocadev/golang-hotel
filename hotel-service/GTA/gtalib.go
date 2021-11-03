package gta

import (
	"fmt"

	hbecommon "../roomres/hbe/common"

	roomresutils "../roomres/utils"
	//repository "../roomres/repository"
	//. "github.com/ahmetb/go-linq"
)

type GtaClient struct {
	HotelsPerBatch  int
	NumberOfThreads int
	SmartTimeout    int

	UserName        string
	SiteId          string
	ApiKey          string
	ProfileCurrency string
	ProfileCountry  string

	SearchEndPoint string

	CommissionProvider hbecommon.ICommissionProvider
	//RepositoryFactory  *repository.RepositoryFactory

	Log roomresutils.ILog
}

func NewHotelBookingProvider(hotelProviderSettings *hbecommon.HotelProviderSettings, logging roomresutils.ILog) *GtaClient {

	gtaClient := &GtaClient{

		Log: logging,

		SearchEndPoint: hotelProviderSettings.SearchEndPoint,

		UserName:           hotelProviderSettings.UserName,
		SiteId:             hotelProviderSettings.SiteId,
		ApiKey:             hotelProviderSettings.ApiKey,
		ProfileCurrency:    hotelProviderSettings.ProfileCurrency,
		ProfileCountry:     hotelProviderSettings.ProfileCountry,
		HotelsPerBatch:     hotelProviderSettings.HotelsPerBatch,
		NumberOfThreads:    hotelProviderSettings.NumberOfThreads,
		SmartTimeout:       hotelProviderSettings.SmartTimeout,
		CommissionProvider: hotelProviderSettings.CommissionProvider,
		//RepositoryFactory:  hotelProviderSettings.RepositoryFactory,
	}

	return gtaClient
}

func (m *GtaClient) CreateHttpHeaders() map[string]string {
	return map[string]string{
		"Content-Type":     "text/xml; charset=UTF-8",
		"Content-Encoding": "UTF-8",
	}
}

func (m *GtaClient) NewSource(currencyCode string) *Source {
	return &Source{
		RequestorId: &RequestorId{
			Client:       m.SiteId,
			Password:     m.ApiKey,
			EmailAddress: m.UserName,
		},
		RequestorPreferences: &RequestorPreferences{
			Language:    GtaDefaultLanguageCode,
			Currency:    currencyCode,
			Country:     m.ProfileCountry,
			RequestMode: &RequestMode{RequestModeText: "SYNCHRONOUS"},
		}}
}

func (m *GtaClient) NewSearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *SearchRequest {

	searchRequest := &SearchRequest{
		//Source: m.NewSource(hotelSearchRequest.CurrencyCode),
		Source: m.NewSource(m.ProfileCurrency),
	}

	return searchRequest
}

func (m *GtaClient) SearchRequest(hotelSearchRequest *hbecommon.HotelSearchRequest) *hbecommon.HotelSearchResponse {

	if hotelSearchRequest.Details {

		// Hu B code here :   2) & 3)  + cxl policy
		//check hotel counts
		if len(hotelSearchRequest.HotelIds) > 1 {
			return nil
		}
		requests := []*SearchHotelPricePaxRequest{}
		mapping_hbe_to_SearchHotelPricePaxRequests(m.HotelsPerBatch, hotelSearchRequest, &requests)
		searchRequests := []*SearchRequest{}

		for _, request := range requests {
			searchRequest := m.NewSearchRequest(hotelSearchRequest)
			searchRequest.SearchHotelPricePaxRequest = request
			searchRequests = append(searchRequests, searchRequest)
		}

		searchResponse := m.SendSearchRequests(searchRequests)

		hotelSearchResponse := &hbecommon.HotelSearchResponse{}

		mapping_searchresponsedetail_to_hbe(searchResponse, hotelSearchRequest, hotelSearchResponse)

		return hotelSearchResponse

	} else {

		requests := []*SearchHotelPricePaxRequest{}
		mapping_hbe_to_SearchHotelPricePaxRequests(m.HotelsPerBatch, hotelSearchRequest, &requests)

		searchRequests := []*SearchRequest{}

		for _, request := range requests {

			searchRequest := m.NewSearchRequest(hotelSearchRequest)

			searchRequest.SearchHotelPricePaxRequest = request

			searchRequests = append(searchRequests, searchRequest)
		}

		searchResponse := m.SendSearchRequests(searchRequests)

		hotelSearchResponse := &hbecommon.HotelSearchResponse{}

		mapping_searchresponse_to_hbe(searchResponse, hotelSearchRequest, hotelSearchResponse)

		return hotelSearchResponse

	}
}

func (m *GtaClient) CancelBooking(bookingCancelRequest *hbecommon.BookingCancelRequest) *hbecommon.BookingCancelResponse {
	// Hu B code here
	cancelRequest := m.CancelRequestSource()
	cancelRequest.CancelRequestDetail = &CancelRequestDetail{
		BookingReference: &BookingReference{
			ReferenceSource: "client",
			Value:           bookingCancelRequest.Ref,
		},
	}

	cancelResponse := m.SendCancelRequest(cancelRequest)

	hbeCancelResponse := &hbecommon.BookingCancelResponse{}

	fmt.Printf("cancelResponse = %+v\n", cancelResponse)
	mapping_cancelresponse_to_hbe(cancelResponse, bookingCancelRequest, hbeCancelResponse)

	return hbeCancelResponse
}

func (m *GtaClient) MakeBooking(hbeBookingRequest *hbecommon.BookingRequest) *hbecommon.BookingResponse {
	// Hu B code here
	bookRequest := m.NewBookingSource()
	bookRequest.AddBookingRequest = makeBookingContent(hbeBookingRequest, m.ProfileCurrency)

	bookResponse := m.SendBookRequest(bookRequest)

	hbeBookingResponse := &hbecommon.BookingResponse{}

	fmt.Printf("bookResponse = %+v\n", bookResponse)
	mapping_bookresponse_to_hbe(bookResponse, hbeBookingRequest, hbeBookingResponse)

	return hbeBookingResponse
}

func (m *GtaClient) GetRoomCxl(roomCxlRequest *hbecommon.RoomCxlRequest) *hbecommon.RoomCxlResponse {
	return nil
}
func (m *GtaClient) CancelBookingConfirm(bookingCancelConfirmationRequest *hbecommon.BookingCancelConfirmationRequest) *hbecommon.BookingCancelConfirmationResponse {
	return nil
}

func (m *GtaClient) SendSpecialRequest(bookingSpecialRequestRequest *hbecommon.BookingSpecialRequestRequest) *hbecommon.BookingSpecialRequestResponse {
	return nil
}

func (m *GtaClient) NewBookingSource() *BookingRequest {
	BookingRequest := &BookingRequest{
		BookingSource: &BookingSource{
			RequestorId: &RequestorId{
				Client:       m.SiteId,
				Password:     m.ApiKey,
				EmailAddress: m.UserName,
			},
			BookingRequestorPreferences: &BookingRequestorPreferences{
				Language:    GtaDefaultLanguageCode,
				Currency:    m.ProfileCurrency,
				Country:     m.ProfileCountry,
				RequestMode: "SYNCHRONOUS",
				ResponseURL: "/ProcessResponse/GetXML",
			}},
	}

	return BookingRequest
}

func (m *GtaClient) CancelRequestSource() *CancelRequest {
	CancelRequest := &CancelRequest{
		BookingSource: &BookingSource{
			RequestorId: &RequestorId{
				Client:       m.SiteId,
				Password:     m.ApiKey,
				EmailAddress: m.UserName,
			},
			BookingRequestorPreferences: &BookingRequestorPreferences{
				Language:    GtaDefaultLanguageCode,
				Currency:    m.ProfileCurrency,
				Country:     m.ProfileCountry,
				RequestMode: "SYNCHRONOUS",
				ResponseURL: "/ProcessResponse/GetXML",
			}},
	}

	return CancelRequest
}
