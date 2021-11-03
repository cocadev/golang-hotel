package derby

import (
	roomresutils "../roomres/utils"
)

type Header struct {
	DistributorId string `json:"distributorId"`
	SupplierId    string `json:"supplierId"`
	Version       string `json:"version"`
	Token         string `json:"token"`
}

type Hotel struct {
	SupplierId string `json:"supplierId"`
	HotelId    string `json:"hotelId"`
	Status     string `json:"status"`
}

type StayRange struct {
	CheckIn  string `json:"checkin"`
	CheckOut string `json:"checkout"`
}

type RoomCriteria struct {
	RoomCount  int   `json:"roomCount"`
	AdultCount int   `json:"adultCount"`
	ChildCount int   `json:"childCount"`
	ChildAges  []int `json:"childAges"`
}

type DateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type Fee struct {
	DateRange *DateRange `json:"dateRange"`
	Fee       *FeeMain   `json:"fee"`
}

type FeeMain struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Amount          float32 `json:"amount"`
	AmountType      string  `json:"amountType"`
	ChargeType      string  `json:"chargeType"`
	EffectivePerson int     `json:"effectivePerson"`
}

type CancelDeadline struct {
	OffsetTimeDropType string `json:"offsetTimeDropType"`
	OffsetTimeUnit     string `json:"offsetTimeUnit"`
	OffsetTimeValue    int    `json:"offsetTimeValue"`
	DealineTime        string `json:"dealineTime"`
}

type PenaltyCharge struct {
	ChargeBase string  `json:"chargeBase"`
	Nights     int     `json:"nights"`
	Amount     float32 `json:"amount"`
	Percent    int     `json:"percent"`
}

type CancelPenalty struct {
	NoShow         bool            `json:"noShow"`
	Cancellable    bool            `json:"cancellable"`
	CancelDeadline *CancelDeadline `json:"cancelDeadline"`
	PenaltyCharge  *PenaltyCharge  `json:"penaltyCharge"`
}

type CancelPolicy struct {
	Code            string           `json:"code"`
	Description     string           `json:"description"`
	CancelPenalties []*CancelPenalty `json:"cancelPenalties"`
}
type RoomRate struct {
	RoomId          string        `json:"roomId"`
	RateId          string        `json:"rateId"`
	Currency        string        `json:"currency"`
	AmountBeforeTax []float32     `json:"amountBeforeTax"`
	AmountAfterTax  []float32     `json:"amountAfterTax"`
	MealPlan        string        `json:"mealPlan"`
	Fees            []*Fee        `json:"fees"`
	CancelPolicy    *CancelPolicy `json:"cancelPolicy"`
	RoomCriteria    *RoomCriteria `json:"roomCriteria"`
	Inventory       int           `json:"inventory"`
}

type RoomRef struct {
	RoomId          string    `json:"roomId"`
	RateId          string    `json:"rateId"`
	Currency        string    `json:"currency"`
	AmountBeforeTax []float32 `json:"amountBeforeTax"`
	AmountAfterTax  []float32 `json:"amountAfterTax"`
}

type SpecialRoomRef struct {
	RoomId   string `json:"roomId"`
	RateId   string `json:"rateId"`
	MealPlan string `json:"mealPlan"`
	Currency string `json:"currency"`
}

type AvailHotel struct {
	SupplierId     string     `json:"supplierId"`
	HotelId        string     `json:"hotelId"`
	Status         string     `json:"status"`
	StayRange      *StayRange `json:"stayRange"`
	Iata           string     `json:"iata"`
	AvailRoomRates []RoomRate `json:"availRoomRates"`
}

type LoyaltyAccount struct {
	ProgramCode string `json:"programCode"`
	AccountId   string `json:"accountId"`
}

type ReservationIds struct {
	DistributorResId string `json:"distributorResId"`
	DerbyResId       string `json:"derbyResId"`
	SupplierResId    string `json:"supplierResId"`
}

type Total struct {
	AmountBeforeTax float32 `json:"amountBeforeTax"`
	AmountAfterTax  float32 `json:"amountAfterTax"`
}

type Payment struct {
	CardCode       string `json:"cardCode"`
	CardNumber     string `json:"cardNumber"`
	CardHolderName string `json:"cardHolderName"`
	ExpireDate     string `json:"expireDate"`
}

type Guest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Index     int    `json:"index"`
}

type SearchRequest struct {
	Header       *Header       `json:"header"`
	HotelId      string        `json:"hotelId"`
	StayRange    *StayRange    `json:"stayRange"`
	RoomCriteria *RoomCriteria `json:"roomCriteria"`
	Iata         string        `json:"iata"`
	// LoyaltyAccount     *LoyaltyAccount `json:"loyaltyAccount"`
}

type SearchResponse struct {
	HotelId      string        `json:"hotelId"`
	Header       *Header       `json:"header"`
	StayRange    *StayRange    `json:"stayRange"`
	RoomCriteria *RoomCriteria `json:"roomCriteria"`
	Iata         string        `json:"iata"`
	RoomRates    []*RoomRate   `json:"roomRates"`
}

func (m *SearchRequest) Clone() *SearchRequest {
	var searchRequest *SearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

type PreBookingRequest struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
	Iata           string          `json:"iata"`
	HotelId        string          `json:"hotelId"`
	StayRange      *StayRange      `json:"stayRange"`
	ContactPerson  *Guest          `json:"contactPerson"`
	RoomCriteria   *RoomCriteria   `json:"roomCriteria"`
	Total          *Total          `json:"total"`
	Payment        *Payment        `json:"payment"`
	LoyaltyAccount *LoyaltyAccount `json:"loyaltyAccount"`
	Guests         []*Guest        `json:"guests"`
	Comments       []string        `json:"comments"`
	RoomRates      []*RoomRate     `json:"roomRates"`
}

type BookingRequest struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
	Iata           string          `json:"iata"`
	HotelId        string          `json:"hotelId"`
	StayRange      *StayRange      `json:"stayRange"`
	ContactPerson  *Guest          `json:"contactPerson"`
	RoomCriteria   *RoomCriteria   `json:"roomCriteria"`
	Total          *Total          `json:"total"`
	Payment        *Payment        `json:"payment"`
	LoyaltyAccount *LoyaltyAccount `json:"loyaltyAccount"`
	Guests         []*Guest        `json:"guests"`
	Comments       []string        `json:"comments"`
	RoomRates      []*RoomRate     `json:"roomRates"`
	BookingToken   string          `json:"bookingToken"`
}

type CancelBookRequest struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
}

type PreBookResponse struct {
	Header       *Header `json:"header"`
	BookingToken string  `json:"bookingToken"`
}

type BookResponse struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
}

type BookingCancelResponse struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
	CancellationId string          `json:"cancellationId"`
}

type BookingDetailRequest struct {
	Header         *Header         `json:"header"`
	ReservationIds *ReservationIds `json:"reservationIds"`
}
