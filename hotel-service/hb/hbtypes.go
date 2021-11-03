package hb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	roomresutils "../roomres/utils"
)

const (
	HBProviderId int = 3
)

type DestinationRequest struct {
	IndexFrom int
	IndexTo   int
}

func (m *DestinationRequest) GenerateUrlParams() string {

	return fmt.Sprintf("/locations/destinations?fields=all&from=%d&to=%d", m.IndexFrom, m.IndexTo)
}

type DestinationResponse struct {
	Destinations []*Destination `json:"destinations"`
}

type Destination struct {
	Code    string       `json:"code"`
	Name    *ContentName `json:"name"`
	IsoCode string       `json:"isoCode"`
}

type ContentName struct {
	Content string `json:"content"`
}

type BookingReportRequest struct {
	DateFrom time.Time
	DateTo   time.Time

	From int
	To   int
}

func (m *BookingReportRequest) GenerateUrlParams() string {

	return fmt.Sprintf("/bookings?start=%s&end=%s&from=%d&to=%d",
		m.DateFrom.Format(roomresutils.LayoutYYYYMMDD),
		m.DateTo.Format(roomresutils.LayoutYYYYMMDD),
		m.From,
		m.To)
}

type BookingReportResponse struct {
	Bookings *BookingReportItems `json:"bookings"`
}

type BookingReportItems struct {
	Items []*BookingReportItem `json:"bookings"`
}

type BookingReportItem struct {
	Reference       string `json:"reference"`
	ClientReference string `json:"clientReference"`

	CreationDate string `json:"creationDate"`

	Hotel *BookingReportHotelItem `json:"hotel"`

	TotalNet float32 `json:"totalNet"`
	Currency string  `json:"currency"`
}

type BookingReportHotelItem struct {
	CheckOut string `json:"checkOut"`
	CheckIn  string `json:"checkIn"`

	Name            string `json:"name"`
	Code            int64  `json:"code"`
	DestinationCode string `json:"destinationCode"`

	Rooms []*BookingReportHotelRoomItem `json:"rooms"`
}

func (m *BookingReportHotelItem) IsConfirmed() bool {

	if len(m.Rooms) == 0 {
		return false
	}

	for _, room := range m.Rooms {
		if !room.IsConfirmed() {
			return false
		}
	}

	return true
}

type BookingReportHotelRoomItem struct {
	Status string `json:"status"`
}

func (m *BookingReportHotelRoomItem) IsConfirmed() bool {
	return strings.ToUpper(m.Status) == "CONFIRMED"
}

type HotelContentRequest struct {
	HotelId string
}

func (m *HotelContentRequest) GenerateUrlParams() string {

	return fmt.Sprintf("/hotels/%s", m.HotelId)
}

type HotelContentResponse struct {
	HotelContent *HotelContent `json:"hotel"`
}

type HotelContent struct {
	Phones []*Phone `json:"phones"`
}

type Phone struct {
	PhoneNumber string `json:"phoneNumber"`
	PhoneType   string `json:"phoneType"`
}

type AvailabilityRequest struct {
	Stay        *Stay        `json:"stay"`
	Occupancies []*Occupancy `json:"occupancies"`
	Hotels      *Hotels      `json:"hotels"`
	Filter      *Filter      `json:"filter"`
}

func (m *AvailabilityRequest) Clone() *AvailabilityRequest {
	var searchRequest *AvailabilityRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

// 2018/09/17, Added by Li
type CheckRateRequest struct {
	Language  string        `json:"language"`
	Upselling string        `json:"upselling"`
	Rooms     []*SimpleRate `json:"rooms"`
}

type Stay struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
}

type Occupancy struct {
	Rooms    int `json:"rooms"`
	Adults   int `json:"adults"`
	Children int `json:"children"`

	Paxes []*Pax `json:"paxes"`
}

type Pax struct {
	Type string `json:"type"`
	Age  int    `json:"age"`
}

type Hotels struct {
	HotelIds []int `json:"hotel"`
}

type Filter struct {
	MaxRooms        int    `json:"maxRooms,omitempty"`
	MaxRatesPerRoom int    `json:"maxRatesPerRoom,omitempty"`
	MinCategory     int    `json:"minCategory,omitempty"`
	MaxCategory     int    `json:"maxCategory,omitempty"`
	PaymentType     string `json:"paymentType"`
	HotelPackage    string `json:"hotelPackage"`
	Packaging       bool   `json:"packaging"`
}

type AvailabilityResponse struct {
	Hotels *ResponseHotels `json:"hotels"`
}

//Added by Li, 20180917
type CheckRateResponse struct {
	Hotel *Hotel `json:"hotel"`
}

type ResponseHotels struct {
	Hotels []*Hotel `json:"hotels"`
}

type Hotel struct {
	HotelId         int     `json:"code"`
	DestinationCode string  `json:"destinationCode"`
	Rooms           []*Room `json:"rooms"`
	CurrencyCode    string  `json:"currency"` //Added by Hu Bing
}

type Room struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Rates []*Rate `json:"rates"`
}

type Rate struct {
	Adults               int                  `json:"adults"`
	Children             int                  `json:"children"`
	RateType             string               `json:"rateType"`
	RateKey              string               `json:"rateKey"`
	Net                  string               `json:"net"`
	CancellationPolicies []*CancelationPolicy `json:"cancellationPolicies"` //added by Hu Bing
	Taxes                *Taxes               `json:"taxes"`                //added by Li, 20180917
	NetValue             float32              `json:"-"`
}

//added by Li, 20180917
type SimpleRate struct {
	RateKey string `json:"rateKey"`
}

//added by Li, 20180917
type Taxes struct {
	AllIncluded bool       `json:"allIncluded"`
	Taxes       []*TaxItem `json:"taxes"`
}

//added by Li, 20180917
type TaxItem struct {
	Included       bool    `json:"included"`
	Percent        float32 `json:"percent"`
	Type           string  `json:"type"` //TAX or FEE
	Amount         float32 `json:"amount"`
	Currency       string  `json:"currency"`
	ClientAmount   float32 `json:"clientAmount"`
	ClientCurrency string  `json:"clientCurrency"`
}

//added by Hu bing
type CancelationPolicy struct {
	Amount string `json:"amount"`
	From   string `json:"from"`
}

func (m *Rate) GetNet() float32 {

	value, err := strconv.ParseFloat(m.Net, 32)

	if err != nil {
		panic(err)
	}

	return float32(value)
}

func (m *Rate) GetOccupancyHash() string {
	return m.GetHash(9)
}

func (m *Rate) GetHash(index int) string {

	return strings.Split(m.RateKey, "|")[index]
}
