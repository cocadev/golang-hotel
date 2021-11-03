package hoteldo

import (
	"encoding/xml"

	roomresutils "../roomres/utils"
)

/*
const (
	JacProviderId       int    = 7
	JacStandardLanguage string = "en-au"
)
*/

/*------------------------ Search Request -----------------------*/

type SearchRequest struct {
	SearchDetails *SearchDetails
}

func (m *SearchRequest) Clone() *SearchRequest {
	var searchRequest *SearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

type BookingRequest struct {
	XMLName       xml.Name       `xml:"Request"`
	Type          string         `xml:"Type,attr"`
	Version       string         `xml:"Version,attr"`
	Affiliateid   string         `xml:"affiliateid"`
	Language      string         `xml:"language"`
	Currency      string         `xml:"currency"`
	Uid           string         `xml:"uid"`
	Firstname     string         `xml:"firstname"`
	Lastname      string         `xml:"lastname"`
	Emailaddress  string         `xml:"emailaddress"`
	ClientCountry string         `xml:"clientcountry"`
	Country       string         `xml:"country"`
	Address       string         `xml:"address"`
	City          string         `xml:"city"`
	State         string         `xml:"state"`
	Zip           string         `xml:"zip"`
	Total         float32        `xml:"total"`
	Phones        []*Phone       `xml:"phones>phone"`
	Hotels        []*HotelBook   `xml:"hotels>hotel"`
	CreditPayment *CreditPayment `xml:"payments>agencycreditpayment"`
}

type BookingResponse struct {
	XMLName        xml.Name    `xml:"BookingResponse"`
	Confirmationid string      `xml:"confirmationid"`
	Currency       string      `xml:"currency"`
	Total          float32     `xml:"total"`
	Statusinternet string      `xml:"statusinternet"`
	Statusbooking  string      `xml:"statusbooking"`
	Statuspayment  string      `xml:"statuspayment"`
	Effective      string      `xml:"effective"`
	Operatorname   string      `xml:"operatorname"`
	Operatoremail  string      `xml:"operatoremail"`
	Rooms          []*RoomBook `xml:"Rooms>Room"`
}

type BookingCancelResponse struct {
	XMLName     xml.Name `xml:"Book"`
	Affiliate   string   `xml:"Affiliate"`
	Number      string   `xml:"Number"`
	Consecutive int      `xml:"Consecutive"`
	Status      string   `xml:"Status"`
}

type ChildAge struct {
	Age int `xml:"Age" json:"age"`
}

type RoomRequest struct {
	Adults   int
	Children int

	ChildAges []*ChildAge
}

type SearchDetails struct {
	AffiliateID  string
	CountryCode  string
	CurrencyCode string

	StartDate     string
	EndDate       string
	HotelIds      []string
	MealPlan      string
	NumberOfRooms int

	RequestedRooms  []*RoomRequest
	DestinationCode int
	LanguageId      string
	/*
		Cities 			 []string
		Geo 			 *Geo
	*/
	Order string

	Details          bool
	SpecificRoomRefs []string
}

/*----------------------  Search Response   ----------------------------*/

type AgencyPublic_Type struct {
	AgencyPublic      float32
	GrossAgencyPublic float32
}

type Available_Type struct {
	Id     int
	Status string
}

type Promotion struct {
	Id              int
	Name            string
	Saving          float32
	Rate            float32
	Destination     string
	PromotionTypeId string
}

type NoShow_Type struct {
	DateFrom string
	Amount   float32
}

type CancellationPolicy_Type struct {
	Id                      int
	Description             string
	Amount                  float32
	DaysToApplyCancellation int
	NightsPenalty           int
	PaymentLimitDay         string
	NoShow                  *NoShow_Type
}

type RateDetail struct {
	Id                 string
	AgencyPublic       *AgencyPublic_Type
	Available          *Available_Type
	AverageGrossNormal float32
	AverageGrossTotal  float32
	AverageNormal      float32
	AverageTotal       float32
	CancellationPolicy *CancellationPolicy_Type
	DutyAmount         float32
	GrossNormal        float32
	GrossTotal         float32
	Normal             float32
	PaxCount           string
	RateKey            string
	RoomsCount         int
	Total              float32
}

type NightlyRate struct {
	XMLName xml.Name `xml:"NightlyRate"`

	Available   string
	Date        string
	GrossNormal float32
	GrossTotal  float32
	Normal      float32
	Total       float32
	PromoId     string
}

type MealPlan struct {
	XMLName xml.Name `xml:"MealPlan"`

	Id                 string
	Name               string
	AgencyPublic       *AgencyPublic_Type
	Available          *Available_Type
	AverageGrossNormal float32
	AverageGrossTotal  float32
	AverageNormal      float32
	AverageTotal       float32
	DutyAmount         float32
	GrossNormal        float32
	GrossTotal         float32
	Normal             float32
	RoomsCount         int
	Total              float32
	Currency           string
	MarketId           string
	Contract           int
	Ratekey            string
	Promotions         []*Promotion   `xml:"Promotions>Promotion"`
	RateDetails        []*RateDetail  `xml:"RateDetails>RateDetail"`
	NightsDetail       []*NightlyRate `xml:"NightsDetail>NightlyRate"`
}

type Room struct {
	XMLName xml.Name `xml:"Room"`

	Id                    string
	Name                  string
	MealPlans             []*MealPlan `xml:"MealPlans>MealPlan"`
	CapacityAdults        int
	CapacityKids          int
	CapacityExtras        int
	CapacityTotal         int
	CapacityChildAgeFrom  int
	CapacityChildAgeTo    int
	CapacityJuniorAgeFrom int
	CapacityJuniorAgeTo   int
	RoomView              string
	Bedding               string
	ImageUrl              string
}

type Chain_Type struct {
	Id   string
	Name string
	Path string
}

type Destination_Type struct {
	Id    string
	Name  string
	Image *Image_Type
}

type Image_Type struct {
	Id     string
	Name   string
	Domain string
	URL    string
	Type   string
}

type Review struct {
	Rating float32
	Source string
	Count  int
}

type Service struct {
	Id          string
	Name        string
	Description string
	ExtraCharge bool
	Order       int
}

type Theme struct {
	Id   string
	Name string
	Path string
}

type InterfaceInfo_Type struct {
	Id string
}

type Hotel struct {
	XMLName xml.Name `xml:"Hotel"`

	Id               int
	Name             string
	CityId           string
	CityName         string
	CountryId        string
	CountryName      string
	Street           string
	ZipCode          string
	CategoryId       string
	LocationId       string
	LocationName     string
	Image            string
	Description      string
	Path             string
	Currency         string
	Status           string // *StatusCode
	Latitude         float32
	Longitude        float32
	AdditionalCharge string
	Order            string

	Chain         *Chain_Type
	Destination   *Destination_Type
	Reviews       []*Review           `xml:"Reviews>Review"`
	Services      []*Service          `xml:"Services>Service"`
	Themes        []*Theme            `xml:"Themes>Theme"`
	InterfaceInfo *InterfaceInfo_Type // InterfaceType - attribute

	Rooms []*Room `xml:"Rooms>Room"`
}

/*
type Filter_Type struct {
	Categories			[]*Category
	Chains 				[]*Chain
	Cities 				[]*City
	Destinations 		[]*Destination
	Locations 			[]*Location
	MealPlans 			[]*MealPlan
	Prices 				[]*Price
	Reviews 			[]*Review
	Services 			[]*Service
	Themes 				[]*Theme
}
*/

type SearchResponse struct {
	XMLName xml.Name `xml:"QuoteHotels"`

	QuoteId string   `xml:"QuoteId"`
	Hotels  []*Hotel `xml:"Hotels>Hotel"`
	//	Filters		  *Filter_Type	`xml:"Filters>Filter"`
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

type PropertyReferenceId struct {
	ReferenceId string `xml:",chardata"`
}

type RoomRef struct {
	Name     string  `json:"name,omitempty" xml:"name,omitempty"`
	Lastname string  `json:"lastname,omitempty" xml:"lastname,omitempty"`
	RoomType string  `json:"roomType,omitempty" xml:"roomType,omitempty"`
	Mealplan string  `json:"mealplan,omitempty" xml:"mealplan,omitempty"`
	MarketId string  `json:"marketid,omitempty" xml:"marketid,omitempty"`
	Contract int     `json:"contract,omitempty" xml:"contract,omitempty"`
	Currency string  `json:"currency,omitempty" xml:"currency,omitempty"`
	Amount   float32 `json:"amount,omitempty" xml:"amount,omitempty"`
	Status   string  `json:"status,omitempty" xml:"status,omitempty"`
	RateKey  string  `json:"ratekey,omitempty" xml:"ratekey,omitempty"`
	Adults   string  `json:"adults,omitempty" xml:"adults,omitempty"`
	Kids     string  `json:"kids,omitempty" xml:"kids,omitempty"`
	K1a      string  `json:"k1a,omitempty" xml:"k1a,omitempty"`
	K2a      string  `json:"k2a,omitempty" xml:"k2a,omitempty"`
	K3a      string  `json:"k3a,omitempty" xml:"k3a,omitempty"`
}

type Phone struct {
	Type   string `xml:"type"`
	Number string `xml:"number"`
}

type HotelBook struct {
	Hotelid       string     `xml:"hotelid"`
	Roomtype      string     `xml:"roomtype"`
	Mealplan      string     `xml:"mealplan"`
	Datearrival   string     `xml:"datearrival"`
	Datedeparture string     `xml:"datedeparture"`
	Marketid      string     `xml:"marketid"`
	Contractid    int        `xml:"contractid"`
	Currency      string     `xml:"currency"`
	Rooms         []*RoomRef `xml:"rooms>room"`
}

type CreditPayment struct {
	Type     string  `xml:"type"`
	Currency string  `xml:"currency"`
	Amount   float32 `xml:"amount"`
}

type RoomBook struct {
	Id                 string                   `xml:"Id"`
	MealPlanId         string                   `xml:"MealPlanId"`
	CancellationPolicy *CancellationPolicy_Type `xml:"CancellationPolicy"`
}
