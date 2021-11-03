package gta

import (
	"encoding/xml"
	"strings"

	. "github.com/ahmetb/go-linq"
)

const (
	GtaDefaultLanguageCode string = "en"
	GtaProviderId          int    = 2
)

type SearchBookingServiceRequest struct {
	XMLName xml.Name `xml:"Request"`

	Source               *Source               `xml:"Source"`
	SearchBookingRequest *SearchBookingRequest `xml:"RequestDetails>SearchBookingRequest"`
}

type SearchBookingRequest struct {
	BookingDateRange *BookingDateRange `xml:"BookingDateRange"`
}

type BookingDateRange struct {
	DateType string `xml:"DateType,attr"`
	FromDate string `xml:"FromDate"`
	ToDate   string `xml:"ToDate"`
}

type SearchBookingServiceReponse struct {
	XMLName xml.Name `xml:"Response"`

	SearchBookings []*SearchBooking `xml:"ResponseDetails>SearchBookingResponse>Bookings>Booking"`
}

type SearchBooking struct {
	BookingReferences []*BookingReference `xml:"BookingReferences>BookingReference"`
}

func (m *SearchBooking) GetReference(source string) *BookingReference {

	if ref := From(m.BookingReferences).Where(func(item interface{}) bool {
		return strings.ToUpper(item.(*BookingReference).ReferenceSource) == strings.ToUpper(source)
	}).First(); ref != nil {
		return ref.(*BookingReference)
	}

	return nil
}

type SearchBookingItemServiceRequest struct {
	XMLName xml.Name `xml:"Request"`

	Source                   *Source                   `xml:"Source"`
	SearchBookingItemRequest *SearchBookingItemRequest `xml:"RequestDetails>SearchBookingItemRequest"`
}

type SearchBookingItemRequest struct {
	BookingReference *BookingReference `xml:"BookingReference"`
}

type SearchBookingItemServiceResponse struct {
	XMLName xml.Name `xml:"Response"`

	SearchBookingItemResponse *SearchBookingItemResponse `xml:"ResponseDetails>SearchBookingItemResponse"`
}

type SearchBookingItemResponse struct {
	BookingStatus *BookingStatus `xml:"BookingStatus"`
	BookingPrice  *BookingPrice  `xml:"BookingPrice"`
	BookingItems  []*BookingItem `xml:"BookingItems>BookingItem"`
}

type BookingStatus struct {
	Code string `xml:"Code,attr"`
}

func (m *BookingStatus) IsConfirmed() bool {
	return strings.ToUpper(m.Code) == "C"
}

type BookingItem struct {
	ItemType  string `xml:"ItemType,attr"`
	CityName  string `xml:"ItemCity"`
	HotelName string `xml:"Item"`
	CheckIn   string `xml:"HotelItem>PeriodOfStay>CheckInDate"`
	CheckOut  string `xml:"HotelItem>PeriodOfStay>CheckOutDate"`
}

type BookingPrice struct {
	Gross float32 `xml:"Gross,attr"`
	Nett  float32 `xml:"Nett,attr"`
}

type BookingReference struct {
	ReferenceSource string `xml:"ReferenceSource,attr"`
	Value           string `xml:",chardata"`
}

type SearchItemInformationServiceRequest struct {
	XMLName xml.Name `xml:"Request"`

	Source                       *Source                       `xml:"Source"`
	SearchItemInformationRequest *SearchItemInformationRequest `xml:"RequestDetails>SearchItemInformationRequest"`
}

type SearchItemInformationRequest struct {
	ItemType        string           `xml:"ItemType,attr"`
	ItemDestination *ItemDestination `xml:"ItemDestination"`
	ItemCode        string           `xml:"ItemCode"`
}

type SearchItemInformationServiceResponse struct {
	XMLName xml.Name `xml:"Response"`

	SearchItems []*SearchItemInformation `xml:"ResponseDetails>SearchItemInformationResponse>ItemDetails>ItemDetail"`
}

type SearchItemInformation struct {
	HotelInformation *HotelInformation `xml:"HotelInformation"`
}

type HotelInformation struct {
	Telephone string `xml:"AddressLines>Telephone"`
}

type SearchRequest struct {
	XMLName xml.Name `xml:"Request"`

	Source                     *Source                     `xml:"Source"`
	SearchHotelPricePaxRequest *SearchHotelPricePaxRequest `xml:"RequestDetails>SearchHotelPricePaxRequest"`
}

type Source struct {
	RequestorId          *RequestorId          `xml:"RequestorID"`
	RequestorPreferences *RequestorPreferences `xml:"RequestorPreferences"`
}

type RequestorId struct {
	Client       string `xml:"Client,attr"`
	EmailAddress string `xml:"EMailAddress,attr"`
	Password     string `xml:"Password,attr"`
}

type RequestorPreferences struct {
	Language string `xml:"Language,attr"`
	Currency string `xml:"Currency,attr"`
	Country  string `xml:"Country,attr"`

	RequestMode *RequestMode `xml:"RequestMode"`
}

type RequestMode struct {
	RequestModeText string `xml:",chardata"`
}

type SearchHotelPricePaxRequest struct {
	ItemDestination           *ItemDestination           `xml:"ItemDestination"`
	ImmediateConfirmationOnly *ImmediateConfirmationOnly `xml:"ImmediateConfirmationOnly"`
	ItemCodes                 []*ItemCode                `xml:"ItemCodes>ItemCode"`
	CheckInDate               string                     `xml:"PeriodOfStay>CheckInDate"`
	CheckOutDate              string                     `xml:"PeriodOfStay>CheckOutDate"`
	IncludePriceBreakdown     *IncludePriceBreakdown     `xml:"IncludePriceBreakdown"`
	IncludeChargeConditions   *IncludeChargeConditions   `xml:"IncludeChargeConditions"`
	PaxRooms                  []*PaxRoom                 `xml:"PaxRooms>PaxRoom"`
}

type ItemDestination struct {
	DestinationType string `xml:"DestinationType,attr"`
	DestinationCode string `xml:"DestinationCode,attr"`
}

type ImmediateConfirmationOnly struct {
}

type ItemCode struct {
	ItemCodeText string `xml:",chardata"`
}

type IncludePriceBreakdown struct {
}

type IncludeChargeConditions struct {
	DateFormatResponse bool `xml:"DateFormatResponse,attr"`
}

type PaxRoom struct {
	Adults    int         `xml:"Adults,attr"`
	Cots      int         `xml:"Cots,attr"`
	RoomIndex int         `xml:"RoomIndex,attr"`
	ChildAges []*ChildAge `xml:"ChildAges>Age"`
}

type ChildAge struct {
	Age int `xml:",chardata"`
}

type DestinationGroup struct {
	DestinationCode string
	HotelIds        []string
}

type SearchResponse struct {
	XMLName xml.Name `xml:"Response"`

	Hotels []*Hotel `xml:"ResponseDetails>SearchHotelPricePaxResponse>HotelDetails>Hotel"`
}

type Hotel struct {
	HotelItemCode *HotelItemCode  `xml:"Item"`
	HotelCity     *HotelCity      `xml:"City"`
	HotelPaxRooms []*HotelPaxRoom `xml:"PaxRoomSearchResults>PaxRoom"`
}

type HotelItemCode struct {
	Code string `xml:"Code,attr"`
}

type HotelCity struct {
	Code string `xml:"Code,attr"`
}

type HotelPaxRoom struct {
	RoomIndex           int                  `xml:"RoomIndex,attr"`
	HotelRoomCategories []*HotelRoomCategory `xml:"RoomCategories>RoomCategory"`
}

type HotelRoomCategory struct {
	Id                string             `xml:"Id,attr"`
	Description       string             `xml:"Description"`
	RoomCategoryPrice *RoomCategoryPrice `xml:"ItemPrice"`
	ChargeConditions  *ChargeConditions  `xml:"ChargeConditions"` //Hu Bing Added
}

type RoomCategoryPrice struct {
	Price                float32 `xml:",chardata"`
	CommissionPercentage string  `xml:"CommissionPercentage"` //Hu Bing Added
	Currency             string  `xml:"Currency"`             //Hu Bing Added
}

//Hu Bing Added From Here
type ChargeConditions struct {
	ChargeConditions []*ChargeCondition `xml:"ChargeCondition"`
}

type ChargeCondition struct {
	Type       string       `xml:"Type,attr"`
	Conditions []*Condition `xml:"Condition"`
}

type Condition struct {
	Charge       bool    `xml:"Charge,attr"`
	ChargeAmount float32 `xml:"ChargeAmount,attr"`
	Currency     string  `xml:"Currency,attr"`
	FromDate     string  `xml:"FromDate,attr"`
	ToDate       string  `xml:"ToDate,attr"`
}
