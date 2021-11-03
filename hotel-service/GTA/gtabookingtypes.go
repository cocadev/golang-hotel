package gta

import "encoding/xml"

type BookingRequest struct {
	XMLName xml.Name `xml:"Request"`

	BookingSource     *BookingSource     `xml:"Source"`
	AddBookingRequest *AddBookingRequest `xml:"RequestDetails>AddBookingRequest"`
}

type BookingResponse struct {
	XMLName xml.Name `xml:"Response"`

	BookingResponseData *BookingResponseData `xml:"ResponseDetails>BookingResponse"`
}

type CancelRequest struct {
	XMLName xml.Name `xml:"Request"`

	BookingSource       *BookingSource       `xml:"Source"`
	CancelRequestDetail *CancelRequestDetail `xml:"RequestDetails>CancelBookingRequest"`
}

type CancelResponse struct {
	XMLName xml.Name `xml:"Response"`

	BookingResponseData *BookingResponseData `xml:"ResponseDetails>BookingResponse"`
}

type BookingSource struct {
	RequestorId                 *RequestorId                 `xml:"RequestorID"`
	BookingRequestorPreferences *BookingRequestorPreferences `xml:"RequestorPreferences"`
}

type BookingRequestorPreferences struct {
	Language string `xml:"Language,attr"`
	Currency string `xml:"Currency,attr"`
	Country  string `xml:"Country,attr"`

	RequestMode string `xml:"RequestMode"`
	ResponseURL string `xml:"ResponseURL"`
}

type AddBookingRequest struct {
	Currency             string        `xml:"Currency,attr"`
	BookingReference     string        `xml:"BookingReference"`
	BookingDepartureDate string        `xml:"BookingDepartureDate"`
	PaxNames             *PaxNames     `xml:"PaxNames"`
	BookingItems         *BookingItems `xml:"BookingItems"`
}

type PaxNames struct {
	PaxNames []*PaxName `xml:"PaxName"`
}

type PaxName struct {
	PaxId   int    `xml:"PaxId,attr"`
	PaxName string `xml:",chardata"`
}

type BookingItems struct {
	BookingItems []*BookingRequestItem `xml:"BookingItem"`
}

type BookingRequestItem struct {
	ItemType      string      `xml:"ItemType,attr"`
	ExpectedPrice float32     `xml:"ExpectedPrice,attr"`
	ItemReference int         `xml:"ItemReference"`
	ItemCity      *ItemCity   `xml:"ItemCity"`
	ItemCode      *ItemCode2  `xml:"Item"`
	CheckInDate   string      `xml:"HotelItem>PeriodOfStay>CheckInDate"`
	CheckOutDate  string      `xml:"HotelItem>PeriodOfStay>CheckOutDate"`
	HotelRooms    *HotelRooms `xml:"HotelItem>HotelRooms"`
}

type ItemCity struct {
	Code string `xml:"Code,attr"`
}
type ItemCode2 struct {
	Code string `xml:"Code,attr"`
}
type HotelRooms struct {
	HotelRooms []*HotelRoom `xml:"HotelRoom"`
}
type HotelRoom struct {
	Code        string `xml:"Code,attr"`
	Id          string `xml:"Id,attr"`
	Description string `xml:"Description"`
	PaxIds      []int  `xml:"PaxIds>PaxId"`
}
type BookingResponseData struct {
	BookingReferences    *BookingReferences    `xml:"BookingReferences"`
	BookingCreationDate  string                `xml:"BookingCreationDate"`
	BookingDepartureDate string                `xml:"BookingDepartureDate"`
	BookingName          string                `xml:"BookingName"`
	BookingPrice         *BookingPrice         `xml:"BookingPrice"`
	BookingStatus        *BookingStatus        `xml:"BookingStatus"`
	PaxNames             *PaxNames             `xml:"PaxNames"`
	BookingItems         *BookingResponseItems `xml:"BookingItems"`
}

type BookingReferences struct {
	BookingReferences []*BookingReference `xml:"BookingReference"`
}

type BookingResponseItems struct {
	BookingResponseItems []*BookingResponseItem `xml:"BookingItem"`
}

type BookingResponseItem struct {
	ItemType                  string             `xml:"ItemType,attr"`
	ExpectedPrice             string             `xml:"ExpectedPrice,attr"`
	ItemReference             string             `xml:"ItemReference"`
	ItemCity                  *ItemCity          `xml:"ItemCity"`
	ItemCode                  *ItemCode2         `xml:"Item"`
	ItemPrice                 *ItemPrice         `xml:"ItemPrice"`
	ItemStatus                *SimpleCodeValue   `xml:"ItemStatus"`
	ItemConfirmationReference string             `xml:"ItemConfirmationReference"`
	HotelItem                 *HotelItem         `xml:"HotelItem"`
	ChargeConditions          []*ChargeCondition `xml:"ChargeConditions>ChargeCondition"`
}

type HotelItem struct {
	CheckInDate  string      `xml:"PeriodOfStay>CheckInDate"`
	CheckOutDate string      `xml:"PeriodOfStay>CheckOutDate"`
	HotelRooms   *HotelRooms `xml:"HotelRooms"`
	Meals        *Meals      `xml:"Meals"`
}
type Meals struct {
	Basis     *SimpleCodeValue `xml:"Basis"`
	Breakfast *SimpleCodeValue `xml:"Breakfast"`
}
type SimpleCodeValue struct {
	Code  string `xml:"Code,attr"`
	Value string `xml:",chardata"`
}

type ItemPrice struct {
	Commission string  `xml:"Commission,attr"`
	Currency   string  `xml:"Currency,attr"`
	Gross      float32 `xml:"Gross,attr"`
	Nett       float32 `xml:"Nett,attr"`
}

type CancelRequestDetail struct {
	BookingReference *BookingReference `xml:"BookingReference"`
}
