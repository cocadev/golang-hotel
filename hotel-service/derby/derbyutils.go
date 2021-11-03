package derby

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	//. "github.com/ahmetb/go-linq"
)

const (
	JacHotelPrefix   string = "jac_"
	LayoutYYYY_MM_DD string = "2006-01-02"
	LayoutYYYYMMDD   string = "20060102"
)

type HotelIdTransformer struct {
	HotelIdInfoItems map[string]*HotelIdInfo
}

func NewHotelIdTransformer() *HotelIdTransformer {
	return &HotelIdTransformer{HotelIdInfoItems: map[string]*HotelIdInfo{}}
}

type HotelIdInfo struct {
	SupplierHotelId string
	Region          string
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func convertTime(_time string) string {

	newTime, err := time.Parse(LayoutYYYY_MM_DD, _time)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d%02d%02d", newTime.Year(), newTime.Month(), newTime.Day())
}

func IsNonRefundableRoom(roomTypeName string) bool {
	return strings.Index(strings.ToLower(roomTypeName), "non-refundable") >= 0
}

func GenerateCancellationPolicyShort(cancelPolicy *CancelPolicy, checkInDate string, price float32, nights int) string {
	arrivalDate, _ := time.Parse("2006-01-02", checkInDate)

	var cxlPolicy bytes.Buffer
	totalRows := len(cancelPolicy.CancelPenalties)
	td := ""
	for i, cancelPenalty := range cancelPolicy.CancelPenalties {
		if i == 0 && cancelPenalty.Cancellable {
			cxlPolicy.WriteString("If you cancel this booking ")

			deadline := cancelPenalty.CancelDeadline
			hours := 0
			if deadline.OffsetTimeUnit == "D" {
				hours += 24 * deadline.OffsetTimeValue
			} else {
				hours += deadline.OffsetTimeValue
			}
			td = GetCalculatedDate(hours, arrivalDate)
			cxlPolicy.WriteString("from now ")
			cxlPolicy.WriteString(fmt.Sprintf(" until %s %s, cancellation charge=$%.2f.",
				td,
				deadline.DealineTime,
				0.0))
		}

		if cxlPolicy.Len() != 0 {
			cxlPolicy.WriteString("<br />")
		}

		if i > 0 && i < totalRows-1 {
			cxlPolicy.WriteString("If you cancel this booking ")

			deadline := cancelPenalty.CancelDeadline
			penaltyCharge := cancelPenalty.PenaltyCharge
			hours := 0
			if deadline.OffsetTimeUnit == "D" {
				hours += 24 * deadline.OffsetTimeValue
			} else {
				hours += deadline.OffsetTimeValue
			}

			fd := GetCalculatedDate(hours, arrivalDate)
			td = GetCalculatedDate(hours-24, arrivalDate)

			amount := float32(0.0)
			if penaltyCharge.ChargeBase == "FullStay" {
				amount = price
			} else if penaltyCharge.ChargeBase == "NightBase" {
				amount = price / float32(nights) * float32(penaltyCharge.Nights)
			} else {
				amount = penaltyCharge.Amount
			}

			cxlPolicy.WriteString(fmt.Sprintf("from %[1]s %s until %[2]s %s, cancellation charge=$%.2f.",
				fd,
				cancelPolicy.CancelPenalties[i].CancelDeadline.DealineTime,
				td,
				cancelPolicy.CancelPenalties[i+1].CancelDeadline.DealineTime,
				amount))
		}

		if i == totalRows-1 {
			cxlPolicy.WriteString("If you cancel this booking ")
			cxlPolicy.WriteString(fmt.Sprintf("from %s %s onwards, cancellation charge $%.2f.",
				td,
				cancelPolicy.CancelPenalties[i-1].CancelDeadline.DealineTime,
				price))
		}
	}

	return cxlPolicy.String()
}

func GetFreeCancelDate(cancelPolicy *CancelPolicy, arrivalDate string) string {
	if cancelPolicy == nil {
		return ""
	}

	message := ""

	cancelPanelties := cancelPolicy.CancelPenalties

	if len(cancelPanelties) > 0 && cancelPanelties[0].Cancellable {
		checkInDate, _ := time.Parse("2006-01-02", arrivalDate)
		deadline := cancelPanelties[0].CancelDeadline
		hours := 0
		if deadline.OffsetTimeUnit == "D" {
			hours += 24 * deadline.OffsetTimeValue
		} else {
			hours += deadline.OffsetTimeValue
		}
		freeCancellationDate := GetCalculatedDate(hours, checkInDate)
		message = fmt.Sprintf("FREE Cancellation until %s %s", freeCancellationDate, deadline.DealineTime)
	}

	return message
}

func GetCalculatedDate(hours int, baseDate time.Time) string {
	calcDate := baseDate.Add((-1) * time.Hour * time.Duration(hours))
	// if calcDate.After(time.Now()) {
	// 	return calcDate.Format("02/01/2006")
	// } else {
	// 	return ""
	// }
	return calcDate.Format("02/01/2006")
}

func getResources(pax string) (string, string, string, string, string) {
	adult := "0"
	children := "0"
	k1a := ""
	k2a := ""
	k3a := ""
	if pax == "" {
		return adult, children, k1a, k2a, k3a
	}
	items := strings.Split(pax, "A")
	if items[1] != "" {
		items = strings.Split(items[1], "C")
		adult = items[0]
		items = strings.Split(items[1], "-")
		children = items[0]
		if items[1] != "" {
			items = strings.Split(items[1], ",")
			if len(items) > 0 {
				k1a = items[0]
			}
			if len(items) > 1 {
				k2a = items[1]
			}
			if len(items) > 2 {
				k3a = items[2]
			}
		}
	}

	return adult, children, k1a, k2a, k3a
}

func ConvertObjectsToArray(objectNames []string, input string) string {

	output := input
	transformed := false

	for _, objectName := range objectNames {

		transformed = true

		for transformed {
			output, transformed = ConvertObjectToArray(objectName, output)
		}

	}

	return output
}

func ConvertObjectToArray(objectName string, input string) (string, bool) {

	objectNameSearch := fmt.Sprintf("\"%s\":{", objectName)

	i := strings.Index(input, objectNameSearch)

	if i >= 0 {

		var output bytes.Buffer

		output.WriteString(input[:i])
		output.WriteString(fmt.Sprintf("\"%s\":", objectName))
		output.WriteString("[{")

		var counter int = 1
		var prevCh rune

		for j, ch := range input[i+len(objectNameSearch):] {

			if ch == '{' && prevCh != '\\' {
				counter++
			} else if ch == '}' && prevCh != '\\' {
				counter--
			}

			if counter == 0 {
				output.WriteString("}]")
				output.WriteString(input[i+len(objectNameSearch)+j+1:])
				break
			}

			output.WriteString(string(ch))

			prevCh = ch
		}

		return output.String(), true

	} else {
		return input, false
	}
}

func ConvertStringToInt(s string) int {
	r, _ := strconv.Atoi(s)
	return r
}

func ConvertStringToFloat32(s string) float32 {
	r, _ := strconv.ParseFloat(s, 32)
	return float32(r)
}

func ConvertIntToString(s int) string {
	return strconv.Itoa(s)
}

func ConvertFloat32ToString(s float32) string {
	return fmt.Sprintf("%f", s)
}

type ClosableChannelPayload func(channel chan interface{})

type ClosableChannel struct {
	Channel  chan interface{}
	IsClosed bool
	Lock     sync.Mutex
}

func NewClosableChannel() *ClosableChannel {
	return &ClosableChannel{IsClosed: false, Channel: make(chan interface{})}
}

func (m *ClosableChannel) Close() {
	m.Lock.Lock()
	m.IsClosed = true
	close(m.Channel)
	m.Lock.Unlock()
}

func (m *ClosableChannel) Execute(payload ClosableChannelPayload) {

	m.Lock.Lock()
	if !m.IsClosed {
		payload(m.Channel)
	}
	m.Lock.Unlock()
}

// setField sets field of v with given name to given value.
func SetField(v interface{}, name string, value string) error {
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("not a field name: %s", name)
	}

	// Field must be exported
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}

	// We expect a string field
	if fv.Kind() != reflect.String {
		return fmt.Errorf("%s is not a string field", name)
	}

	// Set the value
	fv.SetString(value)
	return nil
}

const numberBytes = "01234567890"
const hexBytes = "01234567890ABCDEF"
const complexBytes = "01234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateDistributorResID() string {
	n := 10
	b := make([]byte, n)
	for i := range b {
		b[i] = complexBytes[rand.Intn(len(complexBytes))]
	}
	return string(b)
}
func GenerateDerbyResID() string {
	n := 12
	b := make([]byte, n)
	for i := range b {
		b[i] = hexBytes[rand.Intn(len(hexBytes))]
	}
	return string(b)
}
func GenerateSupplierResID() string {
	n := 8
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}
