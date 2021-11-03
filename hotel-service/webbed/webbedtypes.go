package webbed

import (
	"encoding/xml"

	hbecommon "../roomres/hbe/common"
	roomresutils "../roomres/utils"
)

type SearchRequest struct {
	UserName                      string
	Password                      string
	Language                      string
	Currencies                    string
	CheckInDate                   string
	CheckOutDate                  string
	NumberOfRooms                 string
	Destination                   string
	DestinationID                 string
	HotelIDs                      []string
	ResortIDs                     string
	AccommodationTypes            string
	NumberOfAdults                string
	NumberOfChildren              string
	ChildrenAges                  string
	Infant                        string
	SortBy                        string
	SortOrder                     string
	ExactDestinationMatch         string
	BlockSuperdeal                string
	ShowTransfer                  string
	MealIds                       string
	ShowCoordinates               string
	ShowReviews                   string
	ReferencePointLatitude        string
	ReferencePointLongitude       string
	MaxDistanceFromReferencePoint string
	MinStarRating                 string
	MaxStarRating                 string
	FeatureIds                    string
	MinPrice                      string
	MaxPrice                      string
	ThemeIds                      string
	ExcludeSharedRooms            string
	ExcludeSharedFacilities       string
	PrioritizedHotelIds           string
	TotalRoomsInBatch             string
	PaymentMethodId               string
	CustomerCountry               string
	B2c                           string
}

func (m *SearchRequest) Clone() *SearchRequest {
	var searchRequest *SearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

type PreBookingRequest struct {
	UserName        string
	Password        string
	Currency        string
	Language        string
	CheckInDate     string
	CheckOutDate    string
	RoomId          string
	Rooms           int
	Adults          int
	Children        int
	ChildrenAges    string
	Infant          int
	MealId          string
	CustomerCountry string
	B2c             int
	SearchPrice     string
	Guests          []*hbecommon.Guest
}

type BookingRequest struct {
	UserName                        string
	Password                        string
	Currency                        string
	Language                        string
	CheckInDate                     string
	CheckOutDate                    string
	RoomId                          string
	Rooms                           int
	Email                           string
	Adults                          int
	Children                        int
	Infant                          int
	YourRef                         string
	Specialrequest                  string
	MealId                          string
	AdultGuest1FirstName            string
	AdultGuest2FirstName            string
	AdultGuest3FirstName            string
	AdultGuest4FirstName            string
	AdultGuest5FirstName            string
	AdultGuest6FirstName            string
	AdultGuest7FirstName            string
	AdultGuest8FirstName            string
	AdultGuest9FirstName            string
	AdultGuest1LastName             string
	AdultGuest2LastName             string
	AdultGuest3LastName             string
	AdultGuest4LastName             string
	AdultGuest5LastName             string
	AdultGuest6LastName             string
	AdultGuest7LastName             string
	AdultGuest8LastName             string
	AdultGuest9LastName             string
	ChildrenGuest1FirstName         string
	ChildrenGuest2FirstName         string
	ChildrenGuest3FirstName         string
	ChildrenGuest4FirstName         string
	ChildrenGuest5FirstName         string
	ChildrenGuest6FirstName         string
	ChildrenGuest7FirstName         string
	ChildrenGuest8FirstName         string
	ChildrenGuest9FirstName         string
	ChildrenGuest1LastName          string
	ChildrenGuest2LastName          string
	ChildrenGuest3LastName          string
	ChildrenGuest4LastName          string
	ChildrenGuest5LastName          string
	ChildrenGuest6LastName          string
	ChildrenGuest7LastName          string
	ChildrenGuest8LastName          string
	ChildrenGuest9LastName          string
	ChildrenGuestAge1               string
	ChildrenGuestAge2               string
	ChildrenGuestAge3               string
	ChildrenGuestAge4               string
	ChildrenGuestAge5               string
	ChildrenGuestAge6               string
	ChildrenGuestAge7               string
	ChildrenGuestAge8               string
	ChildrenGuestAge9               string
	PaymentMethodId                 string
	CreditCardType                  string
	CreditCardNumber                string
	CreditCardHolder                string
	CreditCardCVV2                  string
	CreditCardExpYear               string
	CreditCardExpMonth              string
	CustomerEmail                   string
	InvoiceRef                      string
	CommissionAmountInHotelCurrency string
	CustomerCountry                 string
	B2c                             int
	PreBookCode                     string
}

type CancelBookRequest struct {
	UserName  string
	Password  string
	BookingID string
	Language  string
}

type SearchResponse struct {
	XMLName xml.Name `xml:"searchresult"`
	Hotels  []*Hotel `xml:"hotels>hotel"`
}

type PreBookResponse struct {
	XMLName              xml.Name              `xml:"preBookResult"`
	Notes                []*PreBookNote        `xml:"Notes>Note"`
	PreBookCode          string                `xml:"PreBookCode"`
	Price                string                `xml:"Price"`
	CancellationPolicies []*CancellationPolicy `xml:"CancellationPolicies>CancellationPolicy"`
	RoomId               string                `xml:"-"`
	Rooms                int                   `xml:"-"`
	Adults               int                   `xml:"-"`
	Children             int                   `xml:"-"`
	ChildrenAges         string                `xml:"-"`
	Infant               int                   `xml:"-"`
	MealId               string                `xml:"-"`
	Guests               []*hbecommon.Guest    `xml:"-"`
}

type BookResponse struct {
	XMLName xml.Name `xml:"bookResult"`
	Booking *Booking `xml:"booking"`
}

type BookingCancelResponse struct {
	Code             int      `xml:"Code"`
	Cancellationfees []*Price `xml:"CancellationPaymentMethod>cancellationfee"`
}

type Booking struct {
	BookingNumber                        string                `xml:"bookingnumber"`
	HotelId                              string                `xml:"hotel.id"`
	HotelName                            string                `xml:"hotel.name"`
	HotelPhone                           string                `xml:"hotel.phone"`
	NumOfRooms                           string                `xml:"numberofrooms"`
	RoomType                             string                `xml:"room.type	"`
	MealId                               string                `xml:"mealId"`
	Meal                                 string                `xml:"meal"`
	CheckIn                              string                `xml:"checkindate"`
	CheckOut                             string                `xml:"checkoutdate"`
	Prices                               []*Price              `xml:"prices>price"`
	Currency                             string                `xml:"currency"`
	bookingdate                          string                `xml:"bookingdate"`
	Timezone                             string                `xml:"bookingdate.timezone"`
	CancellationPolicies                 []*CancellationPolicy `xml:"cancellationpolicies"`
	EarliestNonFreeCancellationCETDate   string                `xml:"earliestNonFreeCancellationDate.CET"`
	EarliestNonFreeCancellationDateLocal string                `xml:"earliestNonFreeCancellationDate.Local"`
	YourRef                              string                `xml:"yourref"`
	Voucher                              string                `xml:"voucher"`
	BookedBy                             string                `xml:"bookedBy"`
	TransferBooked                       string                `xml:"transferbooked"`
	HotelNotes                           []*PreBookNote        `xml:"hotelNotes>hotelNote"`
	RoomNotes                            []*PreBookNote        `xml:"hotelNotes>roomNotes"`
}

type Hotel struct {
	Id             string      `xml:"hotel.id"`
	DestinationId  string      `xml:"destination_id"`
	ResortId       string      `xml:"resort_id"`
	Transfer       string      `xml:"transfer"`
	RoomTypes      []*RoomType `xml:"roomtypes>roomtype"`
	Type           string      `xml:"type"`
	Name           string      `xml:"name"`
	Address1       string      `xml:"hotel.addr.1"`
	Address2       string      `xml:"hotel.addr.2"`
	ZipCode        string      `xml:"hotel.addr.zip"`
	City           string      `xml:"hotel.addr.city"`
	State          string      `xml:"hotel.addr.state"`
	Country        string      `xml:"hotel.addr.country"`
	CountryCode    string      `xml:"hotel.addr.countrycode"`
	Address        string      `xml:"hotel.address"`
	Headline       string      `xml:"headline"`
	Description    string      `xml:"description"`
	Resort         string      `xml:"resort"`
	Destination    string      `xml:"destination"`
	Classification string      `xml:"classification"`
	TimeZone       string      `xml:"timeZone"`
}

type RoomType struct {
	RoomTypeID       string  `xml:"roomtype.ID"`
	RoomType         string  `xml:"room.type"`
	sharedRoom       bool    `xml:"sharedRoom"`
	sharedFacilities bool    `xml:"sharedFacilities"`
	Rooms            []*Room `xml:"rooms>room"`
}

type Room struct {
	Id                   string                `xml:"id"`
	Beds                 int                   `xml:"beds"`
	Extrabeds            int                   `xml:"extrabeds"`
	Meals                []*Meal               `xml:"meals>meal"`
	CancellationPolicies []*CancellationPolicy `xml:"cancellation_policies>cancellation_policy"`
}

type Meal struct {
	Id        string    `xml:"id"`
	Prices    []*Price  `xml:"prices>price"`
	Name      string    `xml:"name"`
	Discounts []*Amount `xml:"discount>amounts>amount"`
}

type Price struct {
	Currency       string  `xml:"currency,attr"`
	PaymentMethods string  `xml:"paymentMethods,attr"`
	Price          float32 `xml:",chardata"`
}

type Amount struct {
	Currency       string  `xml:"currency,attr"`
	PaymentMethods string  `xml:"paymentMethods,attr"`
	Amount         float32 `xml:",chardata"`
}

type CancellationPolicy struct {
	Deadline   string `xml:"deadline"`
	Percentage string `xml:"percentage"`
	Text       string `xml:"text"`
}

type PreBookNote struct {
	start_date string `xml:"start_date,attr"`
	end_date   string `xml:"end_date,attr"`
	text       string `xml:"text"`
}

type SpecialRequest struct {
	MealId string `json:"mealId"`
}
