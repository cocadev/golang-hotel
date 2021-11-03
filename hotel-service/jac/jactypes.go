package jac

import (
	"encoding/xml"
	roomresutils "roomres/utils"
	"strings"
)

const (
	JacProviderId       int    = 7
	JacStandardLanguage string = "en-au"
)

type LoginDetails struct {
	Login          string `xml:"Login"`
	Password       string `xml:"Password"`
	Locale         string `xml:"Locale"`
	AgentReference string `xml:"AgentReference"`
	CurrencyId     string `xml:"CurrencyID"`
}

type SearchRequest struct {
	XMLName xml.Name `xml:"SearchRequest"`

	LoginDetails  *LoginDetails  `xml:"LoginDetails"`
	SearchDetails *SearchDetails `xml:"SearchDetails"`

	Details bool `xml:"-"`
}

func (m *SearchRequest) Clone() *SearchRequest {
	var searchRequest *SearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

type SearchDetails struct {
	ArrivalDate string `xml:"ArrivalDate"`
	Duration    int    `xml:"Duration"`

	PropertyId           string                 `xml:"PropertyID,omitempty"`
	PropertyReferenceIds []*PropertyReferenceId `xml:"PropertyReferenceIDs>PropertyReferenceID,omitempty"`

	MealBasisId   int `xml:"MealBasisID"`
	MinStarRating int `xml:"MinStarRating"`

	RoomRequests []*RoomRequest `xml:"RoomRequests>RoomRequest"`
}

type PropertyReferenceId struct {
	ReferenceId string `xml:",chardata"`
}

type RoomRequest struct {
	Adults   int `xml:"Adults" json:"adults"`
	Children int `xml:"Children" json:"children"`
	Infants  int `xml:"Infants" json:"infants"`

	ChildAges []*ChildAge `xml:"ChildAges>ChildAge" json:"childages"`
}

type ChildAge struct {
	Age int `xml:"Age" json:"age"`
}

type Exception struct {
}

type ReturnStatus struct {
	Success string `xml:"Success"`
	//Exception *Exception `xml:"Exception"`
	Exception string `xml:"Exception"`
}

func (m *ReturnStatus) IsSuccess() bool {
	return strings.ToLower(m.Success) == "true"
}

type SearchResponse struct {
	XMLName xml.Name `xml:"SearchResponse"`

	PropertyResults []*PropertyResult `xml:"PropertyResults>PropertyResult"`
}

type PropertyResult struct {
	PropertyId          string `xml:"PropertyID"`
	PropertyReferenceId string `xml:"PropertyReferenceID"`

	RoomTypes []*RoomType `xml:"RoomTypes>RoomType"`
}

type RoomType struct {
	Seq                 int        `xml:"Seq" json:"Seq"`
	PropertyRoomTypeId  string     `xml:"PropertyRoomTypeID" json:"PropertyRoomTypeId"`
	BookingToken        string     `xml:"BookingToken" json:"BookingToken"`
	MealBasisId         int        `xml:"MealBasisID" json:"MealBasisId"`
	MealBasis           string     `xml:"MealBasis" json:"MealBasis"`
	RoomType            string     `xml:"RoomType" json:"RoomType"`
	SubTotal            float32    `xml:"SubTotal" json:"SubTotal"`
	Discount            float32    `xml:"Discount" json:"Discount"`
	OnRequest           string     `xml:"OnRequest" json:"OnRequest"`
	Total               float32    `xml:"Total" json:"Total"`
	RSP                 string     `xml:"RSP" json:"RSP"`
	Errata              []*Erratum `xml:"Errata>Erratum" json:"Errata"`
	SpecialOfferApplied string     `xml:"SpecialOfferApplied" json:"SpecialOfferApplied"`
	Adults              int        `xml:"Adults" json:"Adults"`
	Children            int        `xml:"Children" json:"Children"`
	Infants             int        `xml:"Infants" json:"Infants"`

	RequestedRoom *RoomRequest `xml:"-" json:"-"`
}

type Erratum struct {
	Subject     string `xml:"Subject" json:"Subject"`
	Description string `xml:"Description" json:"Description"`
}

type CombinedRoomType struct {
	RoomTypes []*RoomType
}

type RoomTypeRef struct {
	RoomTypeId    string       `json:"roomtypeid"`
	BookingToken  string       `json:"bookingtoken"`
	MealBasisId   int          `json:"mealbasisid"`
	Seq           int          `json:"seq"`
	RequestedRoom *RoomRequest `json:"requestedroom"`
	RoomType      *RoomType    `json:"roomtype"`
}

type CombinedRoomRef struct {
	PropertyId          string         `json:"propertyid"`
	PropertyReferenceId string         `json:"propertyreferenceid"`
	ArrivalDate         string         `json:"checkin"`
	Duration            int            `json:"duration"`
	RoomTypeRefs        []*RoomTypeRef `json:"roomtyperefs"`
}

type PreBookRequest struct {
	XMLName      xml.Name      `xml:"PreBookRequest"`
	LoginDetails *LoginDetails `xml:"LoginDetails"`

	BookingDetails *BookingDetails `xml:"BookingDetails"`

	Index int `xml:"-"`
}

type BookingDetails struct {
	PropertyId  string `xml:"PropertyID"`
	ArrivalDate string `xml:"ArrivalDate"`
	Duration    int    `xml:"Duration"`

	RoomBookings []*RoomBooking `xml:"RoomBookings>RoomBooking"`

	CombinedRoomRef *CombinedRoomRef `xml:"-"`
}

type RoomBooking struct {
	PropertyRoomTypeId string `xml:"PropertyRoomTypeID"`
	BookingToken       string `xml:"BookingToken"`
	MealBasisId        int    `xml:"MealBasisID"`

	Adults   int `xml:"Adults"`
	Children int `xml:"Children"`
	Infants  int `xml:"Infants"`

	ChildAges []*ChildAge `xml:"ChildAges>ChildAge"`
}

type PreBookResponse struct {
	XMLName xml.Name `xml:"PreBookResponse"`

	ReturnStatus    *ReturnStatus `xml:"ReturnStatus"`
	PreBookingToken string        `xml:"PreBookingToken"`
	TotalPrice      float32       `xml:"TotalPrice"`
	TotalCommission float32       `xml:"TotalCommission"`
	VATOnCommission float32       `xml:"VATOnCommission"`

	Cancellations []*Cancellation `xml:"Cancellations>Cancellation"`
}

type Cancellation struct {
	StartDate string  `xml:"StartDate"`
	EndDate   string  `xml:"EndDate"`
	Penalty   float32 `xml:"Penalty"`
}

type BookRequest struct {
	XMLName xml.Name `xml:"BookRequest"`

	LoginDetails *LoginDetails `xml:"LoginDetails"`

	BookDetails *BookDetails `xml:"BookingDetails"`
}

type BookDetails struct {
	PropertyId      string `xml:"PropertyID"`
	PreBookingToken string `xml:"PreBookingToken"`

	ArrivalDate string `xml:"ArrivalDate"`
	Duration    int    `xml:"Duration"`

	LeadGuestTitle     string `xml:"LeadGuestTitle"`
	LeadGuestFirstName string `xml:"LeadGuestFirstName"`
	LeadGuestLastName  string `xml:"LeadGuestLastName"`
	LeadGuestAddress1  string `xml:"LeadGuestAddress1"`
	LeadGuestTownCity  string `xml:"LeadGuestTownCity"`
	LeadGuestPostcode  string `xml:"LeadGuestPostcode"`
	LeadGuestPhone     string `xml:"LeadGuestPhone"`
	LeadGuestEmail     string `xml:"LeadGuestEmail"`
	SpecialRequest     string `xml:"Request"`

	TradeReference string `xml:"TradeReference"`

	BookRooms []*RoomBook `xml:"RoomBookings>RoomBooking"`
}

type RoomBook struct {
	PropertyRoomTypeId string   `xml:"PropertyRoomTypeID"`
	BookingToken       string   `xml:"BookingToken"`
	MealBasisId        int      `xml:"MealBasisID"`
	Adults             int      `xml:"Adults"`
	Children           int      `xml:"Children"`
	Infants            int      `xml:"Infants"`
	Guests             []*Guest `xml:"Guests>Guest"`
}

type Guest struct {
	Type        string `xml:"Type"`
	Title       string `xml:"Title"`
	FirstName   string `xml:"FirstName"`
	LastName    string `xml:"LastName"`
	Age         int    `xml:"Age"`
	Nationality string `xml:"Nationality"`
}

type BookResponse struct {
	XMLName          xml.Name      `xml:"BookResponse"`
	BookingReference string        `xml:"BookingReference"`
	TradeReference   string        `xml:"TradeReference"`
	ReturnStatus     *ReturnStatus `xml:"ReturnStatus"`
	//Exception        string             `xml:"Exception"`
	TotalPrice       float32            `xml:"TotalPrice"`
	PropertyBookings []*PropertyBooking `xml:"PropertyBookings>PropertyBooking"`
}

type PropertyBooking struct {
	PropertyBookingReference string `xml:"PropertyBookingReference"`
	Supplier                 string `xml:"Supplier"`
	SupplierReference        string `xml:"SupplierReference"`
}

type PreCancelRequest struct {
	XMLName      xml.Name      `xml:"PreCancelRequest"`
	LoginDetails *LoginDetails `xml:"LoginDetails"`

	BookingReference string `xml:"BookingReference"`
}

type PreCancelResponse struct {
	XMLName      xml.Name      `xml:"PreCancelResponse"`
	ReturnStatus *ReturnStatus `xml:"ReturnStatus"`

	BookingReference  string  `xml:"BookingReference"`
	CancellationCost  float32 `xml:"CancellationCost"`
	CancellationToken string  `xml:"CancellationToken"`
}

type CancelRequest struct {
	XMLName      xml.Name      `xml:"CancelRequest"`
	LoginDetails *LoginDetails `xml:"LoginDetails"`

	BookingReference  string `xml:"BookingReference"`
	CancellationCost  string `xml:"CancellationCost"`
	CancellationToken string `xml:"CancellationToken"`
}

type CancelResponse struct {
	XMLName      xml.Name      `xml:"CancelResponse"`
	ReturnStatus *ReturnStatus `xml:"ReturnStatus"`
}

type PropertyDetailsRequest struct {
	XMLName      xml.Name      `xml:"PropertyDetailsRequest"`
	LoginDetails *LoginDetails `xml:"LoginDetails"`

	PropertyId string `xml:"PropertyID"`
}

type PropertyDetailsResponse struct {
	XMLName      xml.Name      `xml:"PropertyDetailsResponse"`
	ReturnStatus *ReturnStatus `xml:"ReturnStatus"`

	Description string `xml:"Description"`
}
