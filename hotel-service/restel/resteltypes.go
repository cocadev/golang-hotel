package restel

import (
	"encoding/xml"
	"strings"

	roomresutils "../roomres/utils"
)

const (
	RestelProviderId       int    = 7
	RestelStandardLanguage string = "en-au"
)

type BaseRestelRequest struct {
	XMLName       xml.Name `xml:"peticion"`
	ServiceNumber string   `xml:"tipo"`
	RequestDesc   string   `xml:"nombre"`
	AgencyName    string   `xml:"agencia"`
}

type BaseRestelResponse struct {
	XMLName       xml.Name `xml:"respuesta"`
	ServiceNumber string   `xml:"tipo"`
	RequestDesc   string   `xml:"nombre"`
	AgencyName    string   `xml:"agencia"`
}

type RestelHotelSearchRequest struct {
	Hotel           string `xml:"hotel"`
	Country         string `xml:"pais"`
	Province        string `xml:"provincia"`
	Poblacion       string `xml:"poblacion"`
	Category        int    `xml:"categoria"`
	Radius          int    `xml:"radio"`
	CheckInDate     string `xml:"fechaentrada"`
	CheckOutDate    string `xml:"fechasalida"`
	Group           string `xml:"marca"`
	Affiliation     string `xml:"afiliacion"`
	UserCode        string `xml:"usuario"`
	Type1RoomNumber int    `xml:"numhab1"`
	Type1Guest      string `xml:"paxes1"`
	Type1ChildAges  string `xml:"edades1"`
	Type2RoomNumber int    `xml:"numhab2"`
	Type2Guest      string `xml:"paxes2"`
	Type2ChildAges  string `xml:"edades2"`
	Type3RoomNumber int    `xml:"numhab3"`
	Type3Guest      string `xml:"paxes3"`
	Type3ChildAges  string `xml:"edades3"`
	Language        int    `xml:"idioma"`
	Duplicated      int    `xml:"duplicidad"`
	CompressedXML   int    `xml:"comprimido"`
	Information     int    `xml:"informacion_hotel"`
	Refundable      int    `xml:"tarifas_reembolsables"`
	CompoundHotelId string `xml:"-"`
}

type SearchRequest struct {
	BaseRestelRequest
	RestelSearchParam *RestelHotelSearchRequest `xml:"parametros"`

	Details bool `xml:"-"`
}

func (m *SearchRequest) Clone() *SearchRequest {
	var searchRequest *SearchRequest
	roomresutils.Clone(m, &searchRequest)
	return searchRequest
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

type RoomLine struct {
	MealPlanType           string   `xml:"cod,attr" json:"MealPlanType"`
	PlanPrice              string   `xml:"prr,attr" json:"PlanPrice"`
	Currency               string   `xml:"div,attr" json:"Currency"`
	Status                 string   `xml:"esr,attr" json:"Status"`
	MinPrice               string   `xml:"pvp,attr" json:"MinPrice"`
	FeeRate                string   `xml:"nr,attr" json:"FeeRate"`
	CompressedAvailability []string `xml:"lin" json:"CompressedAvailability"`
}

type RoomPlan struct {
	RoomType string    `xml:"cod,attr"`
	RoomDesc string    `xml:"desc,attr"`
	RoomLine *RoomLine `xml:"reg"`
}

type HotelRestrict struct {
	//XMLName      xml.Name    `xml:"pax"`
	AdultChldren string      `xml:"cod,attr"`
	RoomPlans    []*RoomPlan `xml:"hab"`
}

type HotelResult struct {
	HotelCode              string         `xml:"cod"`
	HotelAffiliation       string         `xml:"Afi"`
	HotelName              string         `xml:"nom"`
	ProvinceCode           string         `xml:"Pro"`
	ProvinceName           string         `xml:"Prn"`
	TownName               string         `xml:"pob"`
	HotelCategory          string         `xml:"cat"`
	CheckInDate            string         `xml:"fen"`
	CheckOutDate           string         `xml:"fsa"`
	DirectPaymentAllowance string         `xml:"pdr"`
	WithCertificate        string         `xml:"cal"`
	HotelBrand             string         `xml:"mar"`
	HotelRestrict          *HotelRestrict `xml:"res>pax"`
	DoesAccet              string         `xml:"pns"`
	MinChildAge            string         `xml:"end"`
	MaxChildAge            string         `xml:"enh"`
	NewHotelCategory       string         `xml:"cat2"`
	CityTax                string         `xml:"city_tax"`
	EstablishmentType      string         `xml:"tipo_establecimiento"`
}

type HotelResults struct {
	XMLName     xml.Name       `xml:"hotls"`
	HotelNumber int            `xml:"num,attr"`
	Hotels      []*HotelResult `xml:"hot"`
}

type ParamResult struct {
	HotelResults *HotelResults `xml:"hotls"`
	ID           string        `xml:"id"`
}

type SearchResponse struct {
	XMLName       xml.Name     `xml:"respuesta"`
	ServiceNumber string       `xml:"tipo"`
	RequestDesc   string       `xml:"nombre"`
	AgencyName    string       `xml:"agencia"`
	ParamResult   *ParamResult `xml:"param"`
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

type ShortRoomRef struct {
	AdultChildren  string `json:"AdultChildren"`
	RoomType       string `json:"RoomType"`
	MealPlanType   string `json:"MealPlanType"`
	RefundableType string `json:"RefundableType"`
	MinPrice       string `json:"MinPrice"`
	MaxPrice       string `json:"MaxPrice"`
	Currency       string `json:"Currency"`
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

type SimpleLines struct {
	CompressedAvailability []string `xml:"lin"`
	OccupantsNames         string   `xml:"ocupantes"` //name1#surname1#.#age1#@name2#surname2#.#age2#@
}

type InternalRef struct {
	Remark      string `json:"Remark"`
	InternalUse string `json:"InternalUse"`
	PaymentForm int    `json:"PaymentForm"`
}

type RestelBookingRequest struct {
	HotelCode         string       `xml:"codigo_hotel"`
	GuestName         string       `xml:"nombre_cliente"`
	Remarks           string       `xml:"observaciones"`
	InternalUse       string       `xml:"num_mensaje"`
	AgencyBookingID   int          `xml:"num_expediente"`
	PaymentForm       int          `xml:"forma_pago"`   //12 - Hotel Payment, 25 - Credit, 44 - Prepayment
	CreditCardType    string       `xml:"tipo_targeta"` //MasterCard/VisaCard/AmExCard/DinersClubCard/enRouteCard/DiscoverCard/JCBCard
	CreditCadNumber   string       `xml:"num_targeta"`
	CVV               string       `xml:"cvv_targeta"`
	MonthOfExpire     string       `xml:"mes_expiracion_targeta"`
	YearOfExpire      string       `xml:"ano_expiracion_targeta"`
	HolderName        string       `xml:"titular_targeta"`
	GuestContactEmail string       `xml:"email"`
	GuestContactPhone string       `xml:"telefono"`
	SimpleLines       *SimpleLines `xml:"res"`
}

type PreBookRequest struct {
	BaseRestelRequest
	RestelBookingParam *RestelBookingRequest `xml:"parametros"`
	Index              int                   `xml:"-"`
}

type ParamPreBookResult struct {
	ReservationStatus      string  `xml:"estado"` //00:ok, otro valor:error
	BookingNumber          string  `xml:"N_localizador"`
	TotalReservationAmount float32 `xml:"importe_total_reseserva"`
	Currency               string  `xml:"divisa_total"`
	Reservado1             int     `xml:"N_mensaje"`
	AgencyBookingID        int     `xml:"N_expediente"`
	Remarks                string  `xml:"observaciones"`
	Reservado2             string  `xml:"datos"`
}

type PreBookResponse struct {
	BaseRestelResponse
	ParamPreBookResult *ParamPreBookResult `xml:"parametros"`
}

type RestelConfirmBookingParam struct {
	BookingNumber string `xml:"localizador"`
	Action        string `xml:"accion"` //AE : Confirm, AI : Cancel Prebooking
}

type ConfirmBookRequest struct {
	BaseRestelRequest
	RestelConfirmBookingParam *RestelConfirmBookingParam `xml:"parametros"`
	Index                     int                        `xml:"-"`
}

type RestelErrorMessage struct {
	Description string `xml:"descripcion"`
}

type ConfirmBookResponseParam struct {
	Status             string              `xml:"estado"`
	BookingNumber      string              `xml:"localizador"`
	ShortBookingNumber string              `xml:"localizador_corto"`
	Error              *RestelErrorMessage `xml:"error"`
}

type ConfirmBookResponse struct {
	BaseRestelResponse
	ConfirmBookResponseParam *ConfirmBookResponseParam `xml:"parametros"`
}

type CancelConfirmedBookRequestParam struct {
	BookingNumber      string `xml:"localizador_largo"`
	ShortBookingNumber string `xml:"localizador_corto"`
}

type CancelConfirmedBookRequest struct {
	BaseRestelRequest
	CancelConfirmedBookRequestParam *CancelConfirmedBookRequestParam `xml:"parametros"`
}

type CancelConfirmedBookResponseParam struct {
	Status        string `xml:"estado"`
	BookingNumber string `xml:"localizador"`
	Comment       string `xml:"localizador_baja"`
}

type CancelConfirmedBookResponse struct {
	BaseRestelResponse
	CancelConfirmedBookResponseParam *CancelConfirmedBookResponseParam `xml:"parametros"`
}

type BookingInfoRequestParam struct {
	HotelCode string `xml:"codigo"`
	Language  string `xml:"idioma"`
}

type BookingInfoRequest struct {
	BaseRestelRequest
	BookingInfoRequestParam *BookingInfoRequestParam `xml:"parametros"`
}

type ReservationLineDetail struct {
	FromDate                     string `xml:"Entrada"`
	ToDate                       string `xml:"Salida"`
	DailyPrice                   string `xml:"Precio"`
	SalePriceExcludingCommission string `xml:"Servicio"`
	CurrencyCode                 string `xml:"Divisa"`
	NumberOfRooms                string `xml:"Num_Hab"`
	RoomType                     string `xml:"Tipo_Hab"`
	PaymentForm                  string `xml:"Forma_pago"`
}

type ReservationLine struct {
	HotPro                string                 `xml:"HotPro"`
	HotelCode             string                 `xml:"Hotel"`
	HotelAffiliation      string                 `xml:"Afiliacion"`
	HotelName             string                 `xml:"Nombre_Hotel"`
	Country               string                 `xml:"Pais"`
	Province              string                 `xml:"Provincia"`
	ReservationLineDetail *ReservationLineDetail `xml:"Fecha"`
}

type BookingInfoDetail struct {
	Country          string `xml:"pais"`
	HotelCode        string `xml:"codigo_hotel"`
	HotelCode2       string `xml:"codigo"`
	HotelAffiliation string `xml:"hot_afiliacion"`
	HotelName        string `xml:"nombre_h"`
	HotelAddress     string `xml:"direccion"`
	ProvinceCode     string `xml:"codprovincia"`
	HotelProvince    string `xml:"provincia"`
	TownCode         string `xml:"codpoblacion"`
	HotelTown        string `xml:"poblacion"`
	ZipCode          string `xml:"cp"`
	CodeDuplicate    string `xml:"coddup"`
	Email            string `xml:"mail"`
	Webpage          string `xml:"web"`
	Phone            string `xml:"telefono"`
	MapPhoto         string `xml:"plano"`
	Description      string `xml:"desc_hotel"`
	RoomNumber       string `xml:"num_habitaciones"`
	HotelHowGet      string `xml:"como_llegar"`
	Category         string `xml:"categoria"`
	Logo             string `xml:"logo_h"`
	CheckInTime      string `xml:"checkin"`
	CheckOutTime     string `xml:"checkout"`
	MinChildAge      string `xml:"edadnindes"`
	MaxChildAge      string `xml:"edadninhas"`
	Currency         string `xml:"currency"`
	Brand            string `xml:"marca"`
}

type BookingInfoResponse struct {
	BaseRestelResponse
	BookingInfoDetail *BookingInfoDetail `xml:"parametros>hotel"`
}

type LastReservationsRequestParam struct {
	SearchMode     string `xml:"selector"`
	FromDate       string `xml:"dia"`
	FromMonth      string `xml:"mes"`
	FromYear       string `xml:"ano"`
	AgencyUserCode string `xml:"usuario"`
	HotelName      string `xml:"hotel"`
	CustomerName   string `xml:"cliente"`
	BookingNumber  string `xml:"localizador"`
	VoucherNumber  string `xml:"bono"`
	InternalUse    string `xml:"usuario"`
}

type LastReservationsRequest struct {
	BaseRestelRequest
	LastReservationsRequestParam *LastReservationsRequestParam `xml:"parametros"`
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

type CancelFeeDetailRequestParam struct {
	Hotel       string   `xml:"hotel"`
	BookingLine []string `xml:"lin"`
	Lang        int      `xml:"idioma"`
}

type CancelFeeDetailRequest struct {
	BaseRestelRequest
	CancelFeeDetailRequestParam *CancelFeeDetailRequestParam `xml:"parametros>datos_reserva"`
}

type CancelPolicy struct {
	RestrictNight                 string  `xml:"fecha,attr" json:"-"`
	NumberDaysPrior               int     `xml:"dias_antelacion"`
	Hour                          float32 `xml:"horas_antelacion"`
	NumberOfNights                int     `xml:"noches_gasto"`
	Percentage                    float32 `xml:"estCom_gasto"`
	Description                   string  `xml:"concepto"`
	BookingIncludeCancellationFee int     `xml:"entra_en_gastos"` //0 - No, 1 - Include
}

type CancelFeeDetailResponseParam struct {
	CancelPolicies []*CancelPolicy `xml:"politicaCanc"`
	RequestKey     string          `xml:"-" json:"-"`
}

type CancelFeeDetailResponse struct {
	BaseRestelResponse
	CancelFeeDetailResponseParam *CancelFeeDetailResponseParam `xml:"parametros"`
}
