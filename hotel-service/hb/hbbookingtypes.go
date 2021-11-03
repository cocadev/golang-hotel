package hb

type BookingRequest struct {
	Holder          *Holder        `json:"holder"`
	BookingRooms    []*BookingRoom `json:"rooms"`
	ClientReference string         `json:"clientReference"`
}

type Holder struct {
	FirstName string `json:"name"`
	LastName  string `json:"surname"`
}

type BookingRoom struct {
	RateKey string        `json:"rateKey"`
	Paxes   []*BookingPax `json:"paxes"`
}

type BookingPax struct {
	RoomID    string `json:"roomId"`
	Type      string `json:"type"`
	FirstName string `json:"name"`
	LastName  string `json:"surname"`
}

type BookingResponse struct {
	BookingResult *BookingResult `json:"booking"`
}

type BookingResult struct {
	Reference             string          `json:"reference"`
	CancellationReference string          `json:"cancellationReference"`
	ClientReference       string          `json:"clientReference"`
	CreationDate          string          `json:"creationDate"`
	Status                string          `json:"status"`
	Holder                *Holder         `json:"holder"`
	Hotel                 *BookingHotel   `json:"hotel"`
	InvoiceCompany        *InvoiceCompany `json:"invoiceCompany"`
	TotalNet              float32         `json:"totalNet"`
	PendingAmount         float32         `json:"pendingAmount"`
	Currency              string          `json:"currency"`
}

type BookingHotel struct {
	CheckOut        string    `json:"checkOut"`
	CheckIn         string    `json:"checkIn"`
	Code            int       `json:"code"`
	Name            string    `json:"name"`
	CategoryCode    string    `json:"categoryCode"`
	CategoryName    string    `json:"categoryName"`
	DestinationCode string    `json:"destinationCode"`
	DestinationName string    `json:"destinationName"`
	Latitude        string    `json:"latitude"`
	Longitude       string    `json:"longitude"`
	TotalNet        string    `json:"totalNet"`
	Currency        string    `json:"currency"`
	Supplier        *Supplier `json:"supplier"`
	Rooms           []*Room   `json:"Rooms"`
}

type InvoiceCompany struct {
	Code               string `json:"code"`
	Company            string `json:"company"`
	RegistrationNumber string `json:"registrationNumber"`
}

type Supplier struct {
	Name      string `json:"name"`
	VatNumber string `json:"vatNumber"`
}
