package teamamerica

import (
	"encoding/xml"

	roomresutils "../roomres/utils"
)

const (
	TeamAmericaProviderId       int    = 9
	TeamAmericaStandardLanguage string = "en-au"
)

type BaseTeamAmericaRequest struct {
	UserName string `xml:"xsd:UserName"`
	Password string `xml:"xsd:Password"`
}

type BaseTeamAmericaResponse struct {
}

type TeamAmericaRequest struct {
	XMLName    xml.Name    `xml:"soapenv:Envelope"`
	SoapEnv    string      `xml:"xmlns:soapenv,attr"`
	Xsd        string      `xml:"xmlns:xsd,attr"`
	SoapHeader *SoapHeader `xml:"soapenv:Header"`
	SoapBody   *SoapBody   `xml:"soapenv:Body"`
}

type SearchResponse struct {
	XMLName  xml.Name                        `xml:"Envelope"`
	SoapBody *TeamAmericaHotelSearchResponse `xml:"Body>PriceSearchResponse"`
}

type BookReserveResponse struct {
	XMLName  xml.Name                        `xml:"Envelope"`
	SoapBody *TeamAmericaBookReserveResponse `xml:"Body>NewMultiItemReservationResponse"`
}

type CancelReservationResponse struct {
	XMLName  xml.Name                              `xml:"Envelope"`
	SoapBody *TeamAmericaCancelReservationResponse `xml:"Body>CancelReservationResponse"`
}

type CancelPolicyResponse struct {
	XMLName  xml.Name                         `xml:"Envelope"`
	SoapBody *TeamAmericaCancelPolicyResponse `xml:"Body>CancellationPolicyResponse"`
}

type FaultResponse struct {
	XMLName  xml.Name `xml:"Envelope"`
	SoapBody *Fault   `xml:"Body>Fault"`
}

type SoapHeader struct {
}
type SoapBody struct {
	BodyData interface{}
}

type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
	Detail      string `xml:"detail"`
}

type TeamAmericaHotelSearchRequest struct {
	XMLName xml.Name `xml:"xsd:PriceSearch"`
	BaseTeamAmericaRequest
	CityCode         string     `xml:"xsd:CityCode"`
	ProductCode      string     `xml:"xsd:ProductCode"`
	Type             string     `xml:"xsd:Type"`
	Occupancy        string     `xml:"xsd:Occupancy"`
	ArrivalDate      string     `xml:"xsd:ArrivalDate"`
	NumberOfNights   int        `xml:"xsd:NumberOfNights"`
	NumberOfRooms    int        `xml:"xsd:NumberOfRooms"`
	DisplayCloseOut  string     `xml:"xsd:DisplayCloseOut"`
	DisplayOnRequest string     `xml:"xsd:DisplayOnRequest"`
	VendorIDs        *VendorIDs `xml:"xsd:VendorIDs"`
}

func (m *TeamAmericaHotelSearchRequest) Clone() *TeamAmericaHotelSearchRequest {
	var searchRequest *TeamAmericaHotelSearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
}

type TeamAmericaHotelSearchResponse struct {
	XMLName             xml.Name             `xml:"PriceSearchResponse"`
	HotelSearchResponse *HotelSearchResponse `xml:"hotelSearchResponse"`
}

type TeamAmericaProductInfoRequest struct {
	XMLName xml.Name `xml:"xsd:ProductInfov2"`
	BaseTeamAmericaRequest
	CityCode    string `xml:"xsd:CityCode"`
	ProductCode string `xml:"xsd:ProductCode"`
	VendorName  string `xml:"xsd:VendorName"`
	VendorID    string `xml:"xsd:VendorID"`
}

type TeamAmericaProductInfoResponse struct {
	XMLName     xml.Name       `xml:"ProductInfov2Response"`
	ProductInfo []*ProductInfo `xml:"response>body"`
}

type TeamAmericaServiceSearchRequest struct {
	XMLName xml.Name `xml:"xsd:ServiceSearch"`
	BaseTeamAmericaRequest
	CityCode         string `xml:"xsd:CityCode"`
	ServiceDate      string `xml:"xsd:ServiceDate"`
	ServiceType      string `xml:"xsd:ServiceType"`
	DisplayClodeOut  string `xml:"xsd:DisplayClodeOut"`
	DisplayOnRequest string `xml:"xsd:DisplayOnRequest"`
}

type TeamAmericaServiceSearchResponse struct {
	XMLName     xml.Name       `xml:"ServiceSearchResponse"`
	ServiceInfo []*ServiceInfo `xml:"serviceSearchResponse>body"`
}

type TeamAmericaBookingReportRequest struct {
	XMLName xml.Name `xml:"xsd:BookingReport"`
	BaseTeamAmericaRequest
	FromDate string `xml:"xsd:FromDate"`
	ToDate   string `xml:"xsd:ToDate"`
}

type TeamAmericaBookingReportResponse struct {
	XMLName           xml.Name           `xml:"BookingReportResponse"`
	BookingReportInfo *BookingReportInfo `xml:"bookingReportResponse"`
}

type TeamAmericaBookReserveRequest struct {
	XMLName xml.Name `xml:"xsd:NewMultiItemReservation"`
	BaseTeamAmericaRequest
	AgentName       string `xml:"xsd:AgentName"`
	AgentEmail      string `xml:"xsd:AgentEmail"`
	ClientReference string `xml:"xsd:ClientReference"`
	Items           *Items `xml:"xsd:Items"`
}

type TeamAmericaBookReserveResponse struct {
	NewMultiItemReservationResponse *NewMultiItemReservationResponse `xml:"newMultiItemReservationResponse"`
}

type NewMultiItemReservationResponse struct {
	ReservationInformations []*ReservationInformation `xml:"ReservationInformation"`
}

type TeamAmericaRetrieveReservationRequest struct {
	XMLName xml.Name `xml:"xsd:RetrieveReservation"`
	BaseTeamAmericaRequest
	ReservationNumber string `xml:"xsd:ReservationNumber"`
}

type TeamAmericaRetrieveReservationResponse struct {
	XMLName             xml.Name                `xml:"RetrieveReservationResponse"`
	ReservationResponse *ReservationInformation `xml:"retrievelReservationResp"`
}

type TeamAmericaCancelReservationRequest struct {
	XMLName xml.Name `xml:"xsd:CancelReservation"`
	BaseTeamAmericaRequest
	ReservationNumber string `xml:"xsd:ReservationNumber"`
}

type TeamAmericaCancelReservationResponse struct {
	CancelReservationResp *CancelReservationResp `xml:"cancelReservationResp"`
}

type TeamAmericaDeleteItemRequest struct {
	XMLName xml.Name `xml:"xsd:DeleteItem"`
	BaseTeamAmericaRequest
	ItemID string `xml:"xsd:ItemID"`
	BKGSrc string `xml:"xsd:BKGSrc"`
}

type TeamAmericaDeleteItemResponse struct {
	XMLName        xml.Name        `xml:"DeleteItemResponse"`
	DeleteItemResp *DeleteItemResp `xml:"deleteItemResp"`
}

type TeamAmericaCancelPolicyRequest struct {
	XMLName xml.Name `xml:"xsd:CancellationPolicy"`
	BaseTeamAmericaRequest
	ProductCode string `xml:"xsd:ProductCode"`
}

type TeamAmericaCancelPolicyResponse struct {
	CancellationPolicyResponse *CancellationPolicyResponse `xml:"cancellationPolicyResp"`
}

type VendorIDs struct {
	VendorID []string `xml:"xsd:VendorID"`
}

type HotelData struct {
	ProductCode     string         `xml:"ProductCode"`
	ProductType     string         `xml:"ProductType"`
	TeamVendorID    int            `xml:"TeamVendorID"`
	ProductDate     string         `xml:"ProductDate"`
	MealPlan        string         `xml:"MealPlan"`
	RoomType        string         `xml:"RoomType"`
	ChildAge        int            `xml:"ChildAge"`
	FamilyPlan      string         `xml:"FamilyPlan"`
	NonRefundable   int            `xml:"NonRefundable"`
	MaxOccupancy    int            `xml:"MaxOccupancy"`
	NightlyInfo     []*NightlyInfo `xml:"NightlyInfo"`
	AverageRate     *AverageRate   `xml:"AverageRate"`
	WebPriority     int            `xml:"WebPriority"`
	CompoundHotelId string         `xml:"-"`
}

type NightlyInfo struct {
	Dates              string  `xml:"Dates"`
	Status             string  `xml:"Status"`
	PromoMessage       string  `xml:"PromoMessage"`
	MinStay            string  `xml:"MinStay"`
	MaxStay            string  `xml:"MaxStay"`
	ArrivalRestriction string  `xml:"ArrivalRestriction"`
	Prices             *Prices `xml:"Prices"`
}

type Prices struct {
	Occupancy  string  `xml:"Occupancy"`
	AdultPrice float32 `xml:"AdultPrice"`
}

type AverageRate struct {
	Occupancy          string  `xml:"Occupancy"`
	AverageNightlyRate float32 `xml:"AverageNightlyRate"`
}

type HotelSearchResponse struct {
	HotelDatas []*HotelData `xml:"body"`
}

type ProductInfo struct {
	ProductCode       string           `xml:"ProductCode"`
	ProductType       string           `xml:"ProductType"`
	ProductName       string           `xml:"ProductName"`
	VendorName        string           `xml:"VendorName"`
	TeamVendorID      string           `xml:"TeamVendorID"`
	CityName          string           `xml:"CityName"`
	ChildAge          string           `xml:"ChildAge"`
	MaximumOccupancy  string           `xml:"MaximumOccupancy"`
	HotelRating       string           `xml:"HotelRating"`
	Latitude          string           `xml:"Latitude"`
	Longitude         string           `xml:"Longitude"`
	VendorImage       string           `xml:"VendorImage"`
	VendorDescription string           `xml:"VendorDescription"`
	VendorAddress1    string           `xml:"VendorAddress1"`
	VendorAddress2    string           `xml:"VendorAddress2"`
	VendorState       string           `xml:"VendorState"`
	MealPlan          string           `xml:"MealPlan"`
	RoomType          string           `xml:"RoomType"`
	NorthStarCode     string           `xml:"NorthStarCode"`
	PropertyImage     []*PropertyImage `xml:"PropertyImage"`
	ProductDetail     string           `xml:"ProductDetail"`
	PermRoomCode      string           `xml:"PermRoomCode"`
	ResortFee         string           `xml:"ResortFee"`
	ResortFeeType     string           `xml:"ResortFeeType"`
	NonRefundable     string           `xml:"NonRefundable"`
	VendorCountryISO  string           `xml:"VendorCountryISO"`
	VendorZip         string           `xml:"VendorZip"`
}

type PropertyImage struct {
	VendorID   string `xml:"xsd:VendorID"`
	Caption    string `xml:"xsd:Caption"`
	Thumbnail  string `xml:"xsd:Thumbnail"`
	ActualSize string `xml:"xsd:ActualSize"`
}

type ServiceInfo struct {
	ProductCode        string `xml:"xsd:ProductCode"`
	ProductType        string `xml:"xsd:ProductType"`
	TransferType       string `xml:"xsd:TransferType"`
	MaximumOccupancy   string `xml:"xsd:MaximumOccupancy"`
	ProductName        string `xml:"xsd:ProductName"`
	ProductDescription string `xml:"xsd:ProductDescription"`
	VendorID           string `xml:"xsd:VendorID"`
	CityName           string `xml:"xsd:CityName"`
	Status             string `xml:"xsd:Status"`
	ProductDate        string `xml:"xsd:ProductDate"`
	Price              string `xml:"xsd:Price"`
	ChildPrice         string `xml:"xsd:ChildPrice"`
	ChildAge           string `xml:"xsd:ChildAge"`
	PromoMessage       string `xml:"xsd:PromoMessage"`
	WebPriority        string `xml:"xsd:WebPriority"`
	NumPax             string `xml:"xsd:NumPax"`
	Latitude           string `xml:"xsd:Latitude"`
	Longitude          string `xml:"xsd:Longitude"`
	VendorImage        string `xml:"xsd:VendorImage"`
	VendorName         string `xml:"xsd:VendorName"`
	VendorDescription  string `xml:"xsd:VendorDescription"`
	VendorAddress1     string `xml:"xsd:VendorAddress1"`
	VendorAddress2     string `xml:"xsd:VendorAddress2"`
}

type BookingReportInfo struct {
	FromDate     string        `xml:"xsd:FromDate"`
	ToDate       string        `xml:"xsd:ToDate"`
	AgentDetails *AgentDetails `xml:"xsd:AgentDetails"`
	Reservations *Reservations `xml:"xsd:Reservations"`
}

type AgentDetails struct {
	AgCode    string `xml:"xsd:AgCode"`
	AgentName string `xml:"xsd:AgentName"`
	Address1  string `xml:"xsd:Address1"`
	Address2  string `xml:"xsd:Address2"`
	City      string `xml:"xsd:City"`
	State     string `xml:"xsd:State"`
	Zip       string `xml:"xsd:Zip"`
	Country   string `xml:"xsd:Country"`
}

type Reservations struct {
	ReservationNumber string `xml:"xsd:ReservationNumber"`
	LastName          string `xml:"xsd:LastName"`
	FirstName         string `xml:"xsd:FirstName"`
	ClientRef         string `xml:"xsd:ClientRef"`
	DepDate           string `xml:"xsd:DepDate"`
	Quantity          string `xml:"xsd:Quantity"`
	ProductCode       string `xml:"xsd:ProductCode"`
	PickUpLocation    string `xml:"xsd:PickUpLocation"`
	Misc1             string `xml:"xsd:Misc1"`
	Misc2             string `xml:"xsd:Misc2"`
}

type Items struct {
	NewItems []*Item `xml:"xsd:NewItem"`
}

type Item struct {
	ProductCode      string      `xml:"xsd:ProductCode"`
	ProductDate      string      `xml:"xsd:ProductDate"`
	Occupancy        string      `xml:"xsd:Occupancy"`
	NumberOfNights   int         `xml:"xsd:NumberOfNights"`
	Language         string      `xml:"xsd:Language"`
	PickUpLocation   string      `xml:"xsd:PickUpLocation"`
	PickUpTime       string      `xml:"xsd:PickUpTime"`
	Quantity         int         `xml:"xsd:Quantity"`
	ItemRemarks      string      `xml:"xsd:ItemRemarks"`
	Passengers       *Passengers `xml:"xsd:Passengers"`
	BelongsToPackage int         `xml:"xsd:BelongsToPackage"`
	PackageCode      int         `xml:"xsd:PackageCode"`
	RateExpected     float32     `xml:"xsd:RateExpected"`
}

type Passengers struct {
	NewPassengers []*Passenger `xml:"xsd:NewPassenger"`
}

type Passenger struct {
	Salutation    string `xml:"xsd:Salutation"`
	FamilyName    string `xml:"xsd:FamilyName"`
	FirstName     string `xml:"xsd:FirstName"`
	PassengerType string `xml:"xsd:PassengerType"`
	PassengerAge  int    `xml:"xsd:PassengerAge"`
}

type ReservationInformation struct {
	ReservationNumber            string            `xml:"ReservationNumber"`
	BookingAgentReferenceNumber  string            `xml:"BookingAgentReferenceNumber"`
	DateBooked                   string            `xml:"DateBooked"`
	ReservationStatus            string            `xml:"ReservationStatus"`
	FaxStatus                    string            `xml:"FaxStatus"`
	ReservationStatusDescription string            `xml:"ReservationStatusDescription"`
	GeneralComments              string            `xml:"GeneralComments"`
	AgencyCommissionAmount       string            `xml:"AgencyCommissionAmount"`
	BookingAgentInfo             *BookingAgentInfo `xml:"BookingAgentInfo"`
	BookItems                    []*BookItems      `xml:"Items"`
	Passengers                   []*BookPassengers `xml:"Passengers"`
	TotalResNetPrice             float32           `xml:"TotalResNetPrice"`
	gross_total                  float32           `xml:"gross_total"`
	paid                         float32           `xml:"paid"`
	gross_due                    float32           `xml:"gross_due"`
	ag_comm_level                float32           `xml:"ag_comm_level"`
	ag_comm_amt                  float32           `xml:"ag_comm_amt"`
	net_due                      float32           `xml:"net_due"`
	cc_fees                      float32           `xml:"cc_fees"`
	ttl_expenses                 float32           `xml:"ttl_expenses"`
	ttl_expenses_paid            float32           `xml:"ttl_expenses_paid"`
	commission_amt               float32           `xml:"commission_amt"`
	commission_rcvd              float32           `xml:"commission_rcvd"`
	amt_comm_amt                 float32           `xml:"amt_comm_amt"`
}

type BookingAgentInfo struct {
	AgencyName    string `xml:"AgencyName"`
	Address1      string `xml:"Address1"`
	Address2      string `xml:"Address2"`
	City          string `xml:"City"`
	State         string `xml:"State"`
	PostalCode    string `xml:"PostalCode"`
	Country       string `xml:"Country"`
	AgentUserName string `xml:"AgentUserName"`
	AgentEmail    string `xml:"AgentEmail"`
}

type BookItems struct {
	UniqueItemID                   string  `xml:"UniqueItemID"`
	ProductCode                    string  `xml:"ProductCode"`
	ItemDate                       string  `xml:"ItemDate"`
	Description                    string  `xml:"Description"`
	Occupancy                      string  `xml:"Occupancy"`
	NumberOfNights                 string  `xml:"NumberOfNights"`
	Language                       string  `xml:"Language"`
	PickUpLocation                 string  `xml:"PickUpLocation"`
	PickUpTime                     string  `xml:"PickUpTime"`
	FlightInfo                     string  `xml:"FlightInfo"`
	AverageNetPricePerNight        float32 `xml:"AverageNetPricePerNight"`
	TotalItemNetPrice              float32 `xml:"TotalItemNetPrice"`
	Quantity                       string  `xml:"Quantity"`
	ItemSupplierConfirmationNumber string  `xml:"ItemSupplierConfirmationNumber"`
	ItemStatusCode                 string  `xml:"ItemStatusCode"`
	ItemRemarks                    string  `xml:"ItemRemarks"`
	ItemInformation                string  `xml:"ItemInformation"`
	ItemCommissionable             string  `xml:"ItemCommissionable"`
	ItemComments                   string  `xml:"ItemComments"`
	SubTotal                       string  `xml:"SubTotal"`
	MealPlan                       string  `xml:"MealPlan"`
	MealPlanID                     string  `xml:"MealPlanID"`
	RoomType                       string  `xml:"RoomType"`
}

type BookPassengers struct {
	UniquePassengerID string `xml:"UniquePassengerID"`
	FamilyName        string `xml:"FamilyName"`
	FirstName         string `xml:"FirstName"`
	Type              string `xml:"Type"`
	Age               string `xml:"Age"`
}

type CancelReservationResp struct {
	ReservationStatusCode string       `xml:"ReservationStatusCode"`
	CancelItems           *CancelItems `xml:"Items"`
}

type CancelItems struct {
	ItemId        string `xml:"ItemId"`
	Description   string `xml:"Description"`
	Status        string `xml:"Status"`
	PenaltyAmount string `xml:"PenaltyAmount"`
	ItemDate      string `xml:"ItemDate"`
}

type DeleteItemResp struct {
	ItemStatusCode string `xml:"ItemStatusCode"`
	Message        string `xml:"Message"`
}

type SearchRoomRef struct {
	ProductCode   string  `json:"ProductCode"`
	ProductDate   string  `json:"ProductDate"`
	MealPlan      string  `json:"MealPlan"`
	RoomType      string  `json:"RoomType"`
	ChildAge      int     `json:"ChildAge"`
	FamilyPlan    string  `json:"FamilyPlan"`
	NonRefundable int     `json:"NonRefundable"`
	MaxOccupancy  int     `json:"MaxOccupancy"`
	MinPrice      float32 `json:"MinPrice"`
	MaxPrice      float32 `json:"MaxPrice"`
}

type RoomRef struct {
	ProductCode        string  `json:"ProductCode"`
	ProductDate        string  `json:"ProductDate"`
	MealPlan           string  `json:"MealPlan"`
	RoomType           string  `json:"RoomType"`
	ChildAge           int     `json:"ChildAge"`
	FamilyPlan         string  `json:"FamilyPlan"`
	NonRefundable      int     `json:"NonRefundable"`
	MaxOccupancy       int     `json:"MaxOccupancy"`
	AverageNightlyRate float32 `json:"AverageNightlyRate"`
}

type CancellationPolicyResponse struct {
	CancelPolicies []*CancelPolicy `xml:"body"`
	ProductCode    string          `xml:"-",json:"-"`
}

type CancelPolicy struct {
	NumberDaysPrior int     `xml:"NumberDaysPrior"`
	PenaltyType     string  `xml:"PenaltyType"`
	PenaltyAmount   float32 `xml:"PenaltyAmount"`
}
