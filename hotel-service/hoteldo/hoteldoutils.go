package hoteldo

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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

func GenerateCancellationPolicyShort(cancelPolicy *CancellationPolicy_Type, checkInDate string) string {
	arrivalDate, _ := time.Parse("2006-01-02", checkInDate)
	noShowDate, _ := time.Parse("2006-01-02", cancelPolicy.NoShow.DateFrom)

	var cxlPolicy bytes.Buffer
	cxlPolicy.WriteString(cancelPolicy.Description)

	freeCancelDate := GetCalculatedDate(cancelPolicy.DaysToApplyCancellation+1, arrivalDate)
	if freeCancelDate != "" {
		cxlPolicy.WriteString("<br />")
		cxlPolicy.WriteString("If you cancel this booking ")
		cxlPolicy.WriteString("from now ")
		cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
			freeCancelDate,
			0.0))
	}
	cancelFromDate := GetCalculatedDate(cancelPolicy.DaysToApplyCancellation, arrivalDate)
	cancelToDate := GetCalculatedDate(1, noShowDate)
	cxlPolicy.WriteString("<br />")
	cxlPolicy.WriteString("If you cancel this booking ")

	cxlPolicy.WriteString(fmt.Sprintf("from %s until %s, cancellation charge=$%.2f.",
		cancelFromDate,
		cancelToDate,
		cancelPolicy.Amount))
	cxlPolicy.WriteString("<br />")
	cxlPolicy.WriteString(fmt.Sprintf("From %s No show charge=$%.2f.",
		cancelPolicy.NoShow.DateFrom,
		cancelPolicy.NoShow.Amount))

	return cxlPolicy.String()
}

func GetFreeCancelDate(freeCondition *CancellationPolicy_Type, arrivalDate string) string {
	message := ""
	if freeCondition != nil {
		checkInDate, _ := time.Parse("2006-01-02", arrivalDate)
		freeCancellationDate := checkInDate.AddDate(0, 0, freeCondition.DaysToApplyCancellation*(-1)-1)
		if freeCancellationDate.After(time.Now()) {
			message = fmt.Sprintf("FREE Cancellation until %s", freeCancellationDate.Format("02/01/2006"))
		}
	}

	return message
}

func GetCalculatedDate(dayPrior int, baseDate time.Time) string {
	calcDate := baseDate.AddDate(0, 0, dayPrior*(-1))
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
