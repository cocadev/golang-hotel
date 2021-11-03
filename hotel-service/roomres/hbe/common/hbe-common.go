package common

import (
	"bytes"
	"fmt"

	//repository "roomres/repository"
	"strconv"
	"strings"
	"time"

	roomresutils "../../../roomres/utils"

	. "github.com/ahmetb/go-linq"
	"github.com/jinzhu/copier"
)

const LayoutYYYMMDD string = "2006-01-02"

const (
	GenderMaleEnum        = 1
	GenderFemaleEnum      = 2
	GenderUnspecifiedEnum = 3
)

const (
	BookingStatusConfirmedEnum     = 5
	BookingStatusPendingEnum       = 3
	BookingStatusFailedEnum        = 7
	BookingStatusFailedRestartEnum = 8
)

type ICommissionProvider interface {
	ReloadCommissions()
	CalculateSellRate(netRate float32) float32
	CalculateCommission(commission float32) float32
}

type HotelProviderSettings struct {
	ProviderId int
	UserName   string
	SecretKey  string
	SiteId     string
	ApiKey     string

	ProfileCurrency string
	ProfileCountry  string

	AuthEndPoint       string
	AuthRefreshEndPont string
	ServiceEndPoint    string

	Metadata string

	SearchEndPoint                          string
	PreBookEndPoint                         string
	BookingConfirmationEndPoint             string
	BookingCancellationEndPoint             string
	BookingCancellationConfirmationEndPoint string
	BookingSpecialRequestEndPoint           string
	BookingDetailsEndpoint                  string
	PropertyDetailsEndpoint                 string
	HotelsPerBatch                          int
	NumberOfThreads                         int
	SmartTimeout                            int

	CreditCardInfo *CreditCardInfo

	CommissionProvider ICommissionProvider
	//RepositoryFactory  *repository.RepositoryFactory
}

type IHotelBookingProvider interface {
	Init()
	SearchRequest(hotelSearchRequest *HotelSearchRequest) *HotelSearchResponse
	MakeBooking(bookingRequest *BookingRequest) *BookingResponse
	CancelBooking(bookingCancelRequest *BookingCancelRequest) *BookingCancelResponse
	CancelBookingConfirm(bookingCancelConfirmationRequest *BookingCancelConfirmationRequest) *BookingCancelConfirmationResponse
	SendSpecialRequest(specialRequestRequest *BookingSpecialRequestRequest) *BookingSpecialRequestResponse
	GetRoomCxl(roomCxlRequest *RoomCxlRequest) *RoomCxlResponse
	GetBookingInfo(request *BookingInfoRequest) *BookingInfoResponse
	ListBookingReport(request *BookingReportRequest) *BookingReportResponse
}

type BookingReportRequest struct {
	BaseRequest

	BookingDateFrom time.Time
	BookingDateTo   time.Time

	BookingIndexFrom int
	BookingIndexTo   int
}

type BookingReportResponse struct {
	Bookings     []*BookingReport
	MoreBookings bool
}

type BookingReport struct {
	Ref         string
	InternalRef string

	CheckIn     time.Time
	CheckOut    time.Time
	BookingDate time.Time

	HotelSupplierRef string
	HotelName        string
	CityName         string

	StatusRaw string
	Status    BookingReportStatus

	NoRooms int

	Total    float32
	Currency string
}

type BookingReportStatus int

const (
	BookingReportStatusConfirmed    BookingReportStatus = 0
	BookingReportStatusNotConfirmed BookingReportStatus = 1
)

type BaseRequest struct {
	SessionId string `json:"SessionId"`
	IpAddress string `json:"IpAddress"`
	UserAgent string `json:"UserAgent"`

	AgencyId   int `json:"AgencyId"`
	ProviderId int `json:"ProviderId"`
}

type BookingInfoRequest struct {
	BaseRequest

	BookingId    string `json:"BookingId"`
	BookingEmail string `json:"BookingEmail"`
	HotelId      int    `json:"HotelId"`
	CurrencyCode string `json:"CurrencyCode"`
}

type BookingInfoResponse struct {
	HotelSupplierId  string `json:"HotelSupplierId"`
	HotelPhoneNumber string `json:"HotelPhoneNumber"`
}

type RoomCxlRequest struct {
	BaseRequest

	HotelId string `json:"HotelId"`
	RoomRef string `json:"RoomRef"`
}

type RoomCxlResponse struct {
	Description string `json:"Description"`
}

type BookingCancelRequest struct {
	BaseRequest

	InternalRef string `json:"InternalRef"`
	Ref         string `json:"Ref"`
}

type BookingCancelResponse struct {
	Status        string          `json:"Status"`
	InternalRef   string          `json:"InternalRef"`
	Ref           string          `json:"Ref"`
	CancelRef     string          `json:"CancelRef"`
	PolicyText    string          `json:"PolicyText"`
	Payment       *Payment        `json:"Payment"`
	Refund        *Payment        `json:"Refund"`
	ErrorMessages []*ErrorMessage `json:"ErrorMessages"`
}

type Payment struct {
	Currency        string  `json:"Currency"`
	AmountInclusive float32 `json:"AmountInclusive"`
}

type BookingCancelConfirmationRequest struct {
	BaseRequest

	InternalRef  string   `json:"InternalRef"`
	Ref          string   `json:"Ref"`
	CancelRef    string   `json:"CancelRef"`
	CancelReason int      `json:"CancelReason"`
	Refund       *Payment `json:"Refund"`
}

type BookingCancelConfirmationResponse struct {
	Status        string          `json:"Status"`
	ErrorMessages []*ErrorMessage `json:"ErrorMessages"`
}

type BookingRequest struct {
	BaseRequest

	InternalRef string        `json:"InternalRef"`
	CheckIn     string        `json:"CheckIn"`
	CheckOut    string        `json:"CheckOut"`
	Total       float32       `json:"Total"`
	Hotel       *BookingHotel `json:"Hotel"`

	IdAddress      string          `json:"IP"`
	Customer       *Customer       `json:"Customer"`
	CreditCardInfo *CreditCardInfo `json:"CreditCardInfo"`
}

func (m *BookingRequest) GetLos() int {
	//los, _ := roomresutils.CalculateLos(m.CheckIn, m.CheckOut)

	los := 0
	return los
}

type BookingResponse struct {
	BookingStatus int             `json:"BookingStatus"`
	Booking       *Booking        `json:"Booking"`
	ErrorMessages []*ErrorMessage `json:"ErrorMessages"`
}

type BookingSpecialRequestRequest struct {
	BaseRequest

	InternalRef    string `json:"InternalRef"`
	Ref            string `json:"Ref"`
	SpecialRequest string `json:"SpecialRequest"`
}

type BookingSpecialRequestResponse struct {
	Status        string          `json:"Status"`
	ErrorMessages []*ErrorMessage `json:"ErrorMessages"`
}

type ErrorMessage struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type Booking struct {
	Ref               string  `json:"BookingRef"`
	ItineraryId       string  `json:"ItineraryId"`
	SupplierReference string  `json:"SupplierReference"`
	SupplierName      string  `json:"SupplierName"`
	SelfServiceUrl    string  `json:"SelfServiceUrl"`
	Total             float32 `json:"Total"`

	CancellationPolicy     string `json:"CancellationPolicy"`
	FreeCancellationPolicy string `json:"FreeCancellationPolicy"`
}

type CreditCardInfo struct {
	CardType       string `json:"CardType"`
	Number         string `json:"Number"`
	ExpiryDate     string `json:"ExpiryDate"`
	Cvc            string `json:"Cvc"`
	HolderName     string `json:"HolderName"`
	CountryOfIssue string `json:"CountryOfIssue"`
	IssuingBank    string `json:"IssuingBank"`
}

type BookingHotel struct {
	HotelId string `json:"HotelId"`

	Rooms []*BookingRoom `json:"Rooms"`
}

type BookingRoom struct {
	Ref string `json:"Ref"`

	Count          int    `json:"Count"`
	Adults         int    `json:"Adults"`
	Children       int    `json:"Children"`
	SpecialRequest string `json:"SpecialRequest"`

	Guests []*Guest `json:"Guests"`
}

type Guest struct {
	Primary              bool   `json:"Primary"`
	Title                string `json:"Title"`
	FirstName            string `json:"FirstName"`
	LastName             string `json:"LastName"`
	CountryOfPassport    string `json:"CountryOfPassport"`
	CountryOfNationality string `json:"CountryOfNationality"`
	Gender               int    `json:"Gender"`
	IsAdult              bool   `json:"IsAdult"`
	Age                  int    `json:"Age"`

	CountryOfResidencyInternalCode   string `json:"-"`
	CountryOfNationalityInternalCode string `json:"-"`
}

type Customer struct {
	Title     string `json:"Title"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`

	PhoneCountryCode string `json:"PhoneCountryCode"`
	PhoneAreaCode    string `json:"PhoneAreaCode"`
	PhoneNumber      string `json:"PhoneNumber"`
}

type Rate struct {
	PerNight     float32
	PerNightBase float32
	Total        float32
	CurrecyCode  string

	NightlyRateTotal    float32
	SurchargeTotal      float32
	CommissionableTotal float32
	GrossProfitOnline   float32

	MinimumSellingTotal float32
}

type Tax struct {
	Type   string
	Amount float32
	Name   string
}

type Surcharge struct {
	Type   string
	Amount float32
	Name   string
}

type RoomType struct {
	Rate     *Rate
	ShortRef string // non-unique reference
	Ref      string

	Description            string
	Nights                 int
	Promo                  bool
	Breakfast              bool
	NonRefundable          bool
	CancellationPolicy     string
	FreeCancellationPolicy string
	Notes                  string

	NormalBeddingOccupancy int
	ExtraBeddingOccupancy  int

	SurchargeInfo string
	Inclusions    string
	Exclusions    string
	Benefits      string

	Taxes      []*Tax
	Surcharges []*Surcharge

	RemainingRooms int
	IsPackage      int

	BoardTypeText string
}

type Hotel struct {
	HotelId                string
	AllPaxNamesRequired    bool
	CheapestRoom           *RoomType
	RoomTypes              []*RoomType
	Notes                  string
	TripAdvisorRating      float32
	TripAdvisorRatingUrl   string
	TripAdvisorReviewCount int

	CustomTag string

	ProviderId int
}

func (m *Hotel) GetTripAdvisorRatingUrl() string {
	return strings.Replace(m.TripAdvisorRatingUrl, "http://", "https://", 1)
}

func (m *Hotel) IsRateDefined() bool {
	return m.CheapestRoom != nil && m.CheapestRoom.Rate != nil
}

type HotelSearchResponse struct {
	Hotels        []*Hotel
	ErrorMessages []*ErrorMessage `json:"ErrorMessages"`
}

type RoomRequest struct {
	Adults    int         `json:"adults"`
	Children  int         `json:"children"`
	Cots      int         `json:"cots"`
	ChildAges []*ChildAge `json:"childages"`
}

func CheckDifferentCapacities(roomRequests []*RoomRequest) bool {

	roomRequest := roomRequests[0]

	for i := 1; i < len(roomRequests); i++ {

		item := roomRequests[i]

		if item.Adults != roomRequest.Adults ||
			item.Children != roomRequest.Children ||
			item.Cots != roomRequest.Cots ||
			len(item.ChildAges) != len(roomRequest.ChildAges) ||
			func() bool {

				var itemAges, requestAges []*ChildAge

				From(item.ChildAges).Sort(func(i interface{}, j interface{}) bool {
					return i.(*ChildAge).Age > j.(*ChildAge).Age
				}).ToSlice(&itemAges)

				From(roomRequest.ChildAges).Sort(func(i interface{}, j interface{}) bool {
					return i.(*ChildAge).Age > j.(*ChildAge).Age
				}).ToSlice(&requestAges)

				for i := 0; i < len(itemAges); i++ {
					if itemAges[i].Age != requestAges[i].Age {
						return true
					}
				}
				return false

			}() {

			return true
		}
	}

	return false
}

func FindMaxCapacityRoomRequest(roomRequests []*RoomRequest) *RoomRequest {

	if len(roomRequests) == 0 {
		return nil
	}

	roomRequest := roomRequests[0]

	for i := 1; i < len(roomRequests); i++ {

		tmp := roomRequests[i]

		if roomRequest.Adults < tmp.Adults ||
			(roomRequest.Adults == tmp.Adults && roomRequest.Children < tmp.Children) ||
			(roomRequest.Adults == tmp.Adults && roomRequest.Children == tmp.Children &&
				len(roomRequest.ChildAges) < len(tmp.ChildAges)) ||
			(roomRequest.Adults == tmp.Adults && roomRequest.Children == tmp.Children &&
				len(roomRequest.ChildAges) == len(tmp.ChildAges) && roomRequest.Cots < tmp.Cots) {

			roomRequest = tmp
		}
	}

	return roomRequest
}

type ChildAge struct {
	Age int `json:"age"`
}

type HotelSearchRequest struct {
	BaseRequest

	CheckIn          string         `json:"CheckIn"`
	CheckOut         string         `json:"CheckOut"`
	RequestedRooms   []*RoomRequest `json:"Rooms"`
	Lat              float32        `json:"Lat"`
	Lon              float32        `json:"Lon"`
	HotelIds         []string       `json:"HotelIds"`
	ExternalRefs     []string       `json:"ExternalRefs"`
	CurrencyCode     string         `json:"CurrencyCode"`
	SpecificRoomRefs []string       `json:"SpecificRoomRefs"`
	Details          bool           `json:"Details"`
	Packaging        bool           `json:"Packaging"`

	SortType string `json:"SortType"`

	AutocompleteId     string `json:"AutocompleteId"`
	HotelFilterId      string `json:"HotelFilterId"`
	RecommendationOnly bool   `json:"RecommendationOnly"`

	SingleHotelSearch bool `json:"-"`
	SingleHotelFilter bool `json:"-"`
	SingleHotelId     int  `json:"-"`

	MaxPrice    float32   `json:"MaxPrice"`
	MinPrice    float32   `json:"MinPrice"`
	StarRatings []float32 `json:"StarRatings"`

	PageIndex int `json:"PageIndex"`
	PageSize  int `json:"PageSize"`
}

func (m *HotelSearchRequest) GetCacheKeyLevel1() string {

	keyItems := []string{
		strconv.Itoa(m.AgencyId),
		m.AutocompleteId,
		m.HotelFilterId,
		strconv.FormatBool(m.RecommendationOnly),
		m.SortType,
		m.CheckIn,
		m.CheckOut,
		strconv.Itoa(m.PageIndex),
		strconv.Itoa(m.PageSize),
		fmt.Sprintf("%.2f", m.MinPrice),
		fmt.Sprintf("%.2f", m.MaxPrice),
		m.GetStarRatingsKey(),
		m.GetRequestedRoomsKey(),
	}

	return strings.Join(keyItems, "-")
}

func (m *HotelSearchRequest) GetCacheKeyLevel2() string {

	keyItems := []string{
		strconv.Itoa(m.AgencyId),
		m.AutocompleteId,
		m.CheckIn,
		m.CheckOut,
		m.GetRequestedRoomsKey(),
	}

	return strings.Join(keyItems, "-")
}

func (m *HotelSearchRequest) GetStarRatingsKey() string {

	var key bytes.Buffer

	for _, starRating := range m.StarRatings {
		key.WriteString(fmt.Sprintf("%.2f-", starRating))
	}

	return key.String()
}

func (m *HotelSearchRequest) GetRequestedRoomsKey() string {

	var key bytes.Buffer

	for _, room := range m.RequestedRooms {
		key.WriteString(strconv.Itoa(room.Adults))
		key.WriteString("-")
		key.WriteString(strconv.Itoa(room.Children))
		key.WriteString("-")
		key.WriteString(strconv.Itoa(room.Cots))
		key.WriteString("-")
		for _, age := range room.ChildAges {
			key.WriteString(strconv.Itoa(age.Age))
			key.WriteString("-")
		}
	}

	return key.String()
}

func (m *HotelSearchRequest) GetLos() int {
	los, _ := roomresutils.CalculateLos(m.CheckIn, m.CheckOut)
	return los
	//return 0
}

func (m *HotelSearchRequest) Clone() *HotelSearchRequest {
	var searchRequest *HotelSearchRequest
	//roomresutils.Clone(m, &searchRequest)
	copier.Copy(&searchRequest, m)
	return searchRequest
}

func (m *HotelSearchRequest) CheckTreshold(hours int) bool {

	var checkInTime time.Time
	var err error

	checkInTime, err = time.Parse(LayoutYYYMMDD, m.CheckIn)

	if err != nil {
		panic(err.Error())
	}

	loc, err1 := time.LoadLocation("Australia/Sydney")
	if err1 != nil {
		panic(err1.Error())
	}

	now := time.Now().In(loc)

	//fmt.Printf("CheckTreshold : %+v, %+v, %+v\n", hours, checkInTime, now)

	return checkInTime.After(now.Add(time.Duration(int64(time.Hour) * int64(hours))))
}

func PopulateHotelTotals(hotels []*Hotel, numOfRooms int, numOfNights int) {

	for _, hotel := range hotels {

		if hotel.CheapestRoom != nil {
			PopulateRoomTotals([]*RoomType{hotel.CheapestRoom}, numOfRooms, numOfNights)
		}
		PopulateRoomTotals(hotel.RoomTypes, numOfRooms, numOfNights)
	}
}

func PopulateRoomTotals(roomTypes []*RoomType, numOfRooms int, numOfNights int) {

	for _, roomType := range roomTypes {

		roomType.Nights = numOfNights

		PopulateRateTotals(roomType.Rate, numOfRooms, numOfNights)
	}
}

func PopulateRateTotals(rate *Rate, numOfRooms int, numOfNights int) {

	rate.Total = CalculateTotal(numOfRooms, numOfNights, rate.PerNight)
}

func CalculateTotal(numOfRooms int, numOfNights int, ratePerNight float32) float32 {
	return float32(numOfRooms*numOfNights) * ratePerNight
}
