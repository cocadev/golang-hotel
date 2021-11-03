package webbed

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/ahmetb/go-linq"
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

func GenerateCancellationPolicyShort(cancelPolicies []*CancellationPolicy, checkInDate string, price float32) string {
	arrivalDate, _ := time.Parse("2006-01-02", checkInDate)

	var cxlPolicy bytes.Buffer
	From(cancelPolicies).OrderByDescending(func(t interface{}) interface{} {
		return t.(*CancellationPolicy).Deadline
	}).ToSlice(&cancelPolicies)
	totalRows := len(cancelPolicies)

	for i, cancelPolicy := range cancelPolicies {
		if i == 0 && cancelPolicy.Deadline != "" {
			cxlPolicy.WriteString("If you cancel this booking ")
			td := GetCalculatedDate(ConvertStringToInt(cancelPolicy.Deadline), arrivalDate)
			cxlPolicy.WriteString("from now ")
			cxlPolicy.WriteString(fmt.Sprintf(" until %s 12:00, cancellation charge=$%.2f.",
				td,
				0.0))
		}

		if cxlPolicy.Len() != 0 {
			cxlPolicy.WriteString("<br />")
		}

		if i < totalRows-1 {
			cxlPolicy.WriteString("If you cancel this booking ")
			fd := GetCalculatedDate(ConvertStringToInt(cancelPolicies[i].Deadline), arrivalDate)
			td := GetCalculatedDate(ConvertStringToInt(cancelPolicies[i+1].Deadline), arrivalDate)

			cxlPolicy.WriteString(fmt.Sprintf("from %[1]s 12:00 until %[2]s 12:00, cancellation charge=$%.2f.",
				fd,
				td,
				ConvertStringToFloat32(cancelPolicy.Percentage)*price/100.0))
		} else {
			cxlPolicy.WriteString("If you cancel this booking ")
			fd := GetCalculatedDate(ConvertStringToInt(cancelPolicies[i].Deadline), arrivalDate)
			cxlPolicy.WriteString(fmt.Sprintf("from %s 12:00 onwards, cancellation charge $%.2f.",
				fd,
				ConvertStringToFloat32(cancelPolicy.Percentage)*price/100.0))
		}
	}

	return cxlPolicy.String()
}

func GetFreeCancelDate(cancelPolicies []*CancellationPolicy, arrivalDate string) string {
	if cancelPolicies[0].Deadline == "" {
		return ""
	}

	message := ""

	From(cancelPolicies).OrderByDescending(func(t interface{}) interface{} {
		return t.(*CancellationPolicy).Deadline
	}).ToSlice(&cancelPolicies)

	if cancelPolicies != nil && len(cancelPolicies) > 0 {
		checkInDate, _ := time.Parse("2006-01-02", arrivalDate)
		freeCancellationDate := GetCalculatedDate(ConvertStringToInt(cancelPolicies[0].Deadline), checkInDate)
		message = fmt.Sprintf("FREE Cancellation until %s 12:00", freeCancellationDate)
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

func GetRawHotelId(hotelId string) string {
	if strings.Index(hotelId, "_") >= 0 {

		values := strings.Split(hotelId, "_")

		return values[1]
	}

	return hotelId
}
