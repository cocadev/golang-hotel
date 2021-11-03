package restel

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	hbecommon "../roomres/hbe/common"
	. "github.com/ahmetb/go-linq"

	//roomresutils "roomres/utils"
	"strings"
	"time"
)

const (
	RestelHotelPrefix string = "restel_"
)

func CompareRoomRefs(roomRefA, roomRefB string) bool {

	return roomRefA == roomRefB
}

func GenerateRoomTypeCombinations(requestedRooms int, roomTypes []*RoomType) []*CombinedRoomType {

	combinedRoomTypes := []*CombinedRoomType{&CombinedRoomType{}}

	for i := 0; i < requestedRooms; i++ {

		combinedRoomTypesNew := []*CombinedRoomType{}

		for _, roomType := range roomTypes {

			if roomType.Seq == i+1 {

				for _, combination := range combinedRoomTypes {

					combinationNew := &CombinedRoomType{}

					for _, roomTypeC := range combination.RoomTypes {

						combinationNew.RoomTypes = append(combinationNew.RoomTypes, roomTypeC)
					}

					combinationNew.RoomTypes = append(combinationNew.RoomTypes, roomType)

					combinedRoomTypesNew = append(combinedRoomTypesNew, combinationNew)
				}

			}
		}

		combinedRoomTypes = combinedRoomTypesNew
	}

	return combinedRoomTypes
}

func FilterRoomTypeCombinations(combinations []*CombinedRoomType) []*CombinedRoomType {

	filtered := []*CombinedRoomType{}

	for _, combination := range combinations {

		var roomTypeRef string
		same := true
		for i, roomType := range combination.RoomTypes {

			if i == 0 {
				roomTypeRef = roomType.PropertyRoomTypeId
				continue
			}

			if roomTypeRef != roomType.PropertyRoomTypeId {
				same = false
				break
			}
		}

		if same {
			filtered = append(filtered, combination)
		}
	}

	return filtered
}

func GenerateCombinedRoomRef(propertyId, propertyReferenceId, arrivalDate string, duration int, roomTypes []*RoomType) *CombinedRoomRef {

	combinedRoomRef := &CombinedRoomRef{
		PropertyId:          propertyId,
		PropertyReferenceId: propertyReferenceId,
		ArrivalDate:         arrivalDate,
		Duration:            duration,
	}

	for _, roomType := range roomTypes {

		combinedRoomRef.RoomTypeRefs = append(
			combinedRoomRef.RoomTypeRefs,
			&RoomTypeRef{
				RoomTypeId:    roomType.PropertyRoomTypeId,
				BookingToken:  roomType.BookingToken,
				MealBasisId:   roomType.MealBasisId,
				Seq:           roomType.Seq,
				RequestedRoom: roomType.RequestedRoom,
				RoomType:      roomType,
			})
	}

	return combinedRoomRef
}

func EncodeCombinedRoomRef(combinedRoomRef *CombinedRoomRef) string {

	outputBuffer := new(bytes.Buffer)

	json.NewEncoder(outputBuffer).Encode(combinedRoomRef)

	return outputBuffer.String()
}

func DecodeCombinedRoomRef(value string, combinedRoomTypeRef *CombinedRoomRef) {

	err := json.NewDecoder(bytes.NewBuffer([]byte(value))).Decode(combinedRoomTypeRef)
	if err != nil {
		panic(err.Error())
	}
}

func GenerateCombinedRoomTypeName(roomTypes []*RoomType) string {

	var combinedRoomTypeName bytes.Buffer

	for _, roomType := range roomTypes {

		if combinedRoomTypeName.Len() > 0 {
			combinedRoomTypeName.WriteString(" + ")
		}

		combinedRoomTypeName.WriteString(roomType.RoomType)
	}

	return combinedRoomTypeName.String()
}

func IsBB(roomTypes []*RoomType) bool {

	if len(roomTypes) == 0 {
		return false
	}

	var isBB bool = true

	for _, roomType := range roomTypes {

		isBB = isBB && roomType.MealBasisId != 1 /* room only */

		if !isBB {
			break
		}
	}

	return isBB
}

var IsoLayout string = "2006-01-02"
var DDMMYYYYLayout string = "02/01/2006"
var CxlEndDate string = "2099-12-31T00:00:00"

func ConvertToDate(value string) (time.Time, error) {

	items := strings.Split(value, "T")

	return time.Parse(IsoLayout, items[0])
}

func CombinedSpecialRequests(hbeBookingRequest *hbecommon.BookingRequest) string {

	var specialRequests bytes.Buffer

	for _, hbeRoomType := range hbeBookingRequest.Hotel.Rooms {

		if specialRequests.Len() > 0 {
			specialRequests.WriteString("<br />")
		}

		specialRequests.WriteString(fmt.Sprintf("%s", hbeRoomType.SpecialRequest))
	}

	return specialRequests.String()
}

func GetGuestType(age int) string {
	return "Adult"
}

func GetPropertyDescriptionItem(propertyDescription, itemHeader string) string {

	descriptionItem := ""

	a := strings.Index(propertyDescription, itemHeader)

	if a >= 0 {

		b := strings.Index(propertyDescription[a+len(itemHeader):], "#")

		if b >= 0 {
			descriptionItem = propertyDescription[a+len(itemHeader) : b]
		} else {
			descriptionItem = propertyDescription[a+len(itemHeader):]
		}
	}

	return descriptionItem
}

type Serializer struct {
	DebugMode bool
}

func NewSerializer(debugMode bool) *Serializer {
	return &Serializer{DebugMode: debugMode}
}

func (m *Serializer) Serialize(object interface{}) ([]byte, error) {

	var objectBinary []byte
	var err error

	if m.DebugMode {
		objectBinary, err = xml.MarshalIndent(object, "", "\t")
	} else {
		objectBinary, err = xml.Marshal(object)
	}

	return objectBinary, err
}

func CurrentDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func ExtractSupplierHotelId(prefix, hotelId string) string {

	if strings.Index(strings.ToLower(hotelId), strings.ToLower(prefix)) == 0 {

		return hotelId[len(prefix):]
	} else {
		return hotelId
	}
}

func GenerateSupplierHotelId(prefix, supplierHotelId string) string {

	return fmt.Sprintf("%s%s", prefix, supplierHotelId)
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

type IntegrationHttpRequest struct {
	Method               string
	UrlParameters        string
	BodyParameters       map[string]string
	RequestBody          []byte
	RequestBodySpecified bool
	Timeout              time.Duration
}

type IntegrationHttpResponse struct {
	ResponseBody []byte
	Err          error
}

type IntegrationHttp struct {
	EndPoint string
	Headers  map[string]string
}

func NewIntegrationHttp(endPoint string, headers map[string]string) *IntegrationHttp {
	return &IntegrationHttp{EndPoint: endPoint, Headers: headers}
}

func (m *IntegrationHttp) Send(body []byte) ([]byte, error) {

	client := http.Client{}

	req, reqErr := http.NewRequest("POST", m.EndPoint, bytes.NewBuffer(body))

	if reqErr != nil {
		return nil, reqErr
	}

	for headerName, headerValue := range m.Headers {
		req.Header.Set(headerName, headerValue)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	respData, errResp := ioutil.ReadAll(resp.Body)

	return respData, errResp
}

func GenerateBodyUrl(bodyParams map[string]string) string {

	var bodyUrl bytes.Buffer

	for key, value := range bodyParams {

		if bodyUrl.Len() > 0 {
			bodyUrl.WriteString("&")
		}

		bodyUrl.WriteString(url.QueryEscape(key))
		bodyUrl.WriteString("=")
		bodyUrl.WriteString(url.QueryEscape(value))
	}

	return bodyUrl.String()
}

func (m *IntegrationHttp) SendRequest(request *IntegrationHttpRequest) (response *IntegrationHttpResponse) {

	client := http.Client{Timeout: request.Timeout}

	var req *http.Request
	var reqErr error

	//fmt.Printf("endpoint method, url : %s, %s\n", request.Method, m.EndPoint+request.UrlParameters)

	if request.RequestBodySpecified {

		if request.BodyParameters == nil {
			req, reqErr = http.NewRequest(request.Method,
				m.EndPoint+request.UrlParameters, bytes.NewBuffer(request.RequestBody))
		} else {
			req, reqErr = http.NewRequest(request.Method,
				m.EndPoint+request.UrlParameters, bytes.NewBuffer([]byte(GenerateBodyUrl(request.BodyParameters))))
		}

	} else {

		req, reqErr = http.NewRequest(request.Method,
			m.EndPoint+request.UrlParameters, nil)
	}

	if reqErr != nil {
		return &IntegrationHttpResponse{Err: reqErr}
	}

	if m.Headers != nil {
		for headerName, headerValue := range m.Headers {
			req.Header.Set(headerName, headerValue)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return &IntegrationHttpResponse{Err: err}
	}

	response = &IntegrationHttpResponse{}

	response.ResponseBody, response.Err = ioutil.ReadAll(resp.Body)

	return response

}

func (m *IntegrationHttp) SendRequestWithCompression(request *IntegrationHttpRequest) (response *IntegrationHttpResponse) {

	client := http.Client{Timeout: request.Timeout}

	var req *http.Request
	var reqErr error

	if request.RequestBodySpecified {

		/*
			var buf bytes.Buffer
			g := gzip.NewWriter(&buf)
			if _, err := g.Write(request.RequestBody); err != nil {
				return &IntegrationHttpResponse{Err: err}
			}
			if err := g.Close(); err != nil {
				return &IntegrationHttpResponse{Err: err}
			}

			req, reqErr = http.NewRequest(request.Method,
				m.EndPoint+request.UrlParameters, &buf)
		*/
		req, reqErr = http.NewRequest(request.Method,
			m.EndPoint+request.UrlParameters, bytes.NewBuffer(request.RequestBody))

	} else {
		req, reqErr = http.NewRequest(request.Method,
			m.EndPoint+request.UrlParameters, nil)
	}

	if reqErr != nil {
		return &IntegrationHttpResponse{Err: reqErr}
	}

	if m.Headers != nil {
		for headerName, headerValue := range m.Headers {
			req.Header.Set(headerName, headerValue)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return &IntegrationHttpResponse{Err: err}
	}

	response = &IntegrationHttpResponse{}

	response.ResponseBody, response.Err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return &IntegrationHttpResponse{Err: err}
	}

	if strings.ToLower(resp.Header.Get("Content-Encoding")) == "gzip" {
		gz, err := gzip.NewReader(bytes.NewBuffer(response.ResponseBody))
		if err != nil {
			return &IntegrationHttpResponse{Err: err}
		}
		defer gz.Close()

		var tmp []byte

		_, err = gz.Read(tmp)

		if err != nil {
			return &IntegrationHttpResponse{Err: err}
		}

		response.ResponseBody = tmp
	}

	return response

}

func SerializeArrayInt(array []int) string {
	var output bytes.Buffer
	for _, value := range array {
		if output.Len() > 0 {
			output.WriteString(",")
		}
		output.WriteString(strconv.Itoa(value))
	}
	return output.String()
}

func EncryptText(text string) (string, error) {

	encryptionKey := []byte(os.Getenv("ROOMRES_HBE_ENCRYPTION_KEY"))

	encryptedBytes, err := Encrypt(encryptionKey, []byte(text))

	encryptedText := ""

	if err == nil {
		encryptedText = fmt.Sprintf("%0x", encryptedBytes)
	}

	return encryptedText, err
}

func DecryptText(text string) (string, error) {

	encryptionKey := []byte(os.Getenv("ROOMRES_HBE_ENCRYPTION_KEY"))

	originalBytes, hexError := hex.DecodeString(text)

	if hexError != nil {

		return "", hexError
	}

	decryptedBytes, err := Decrypt(encryptionKey, originalBytes)

	decryptedText := ""

	if err == nil {

		decryptedText = fmt.Sprintf("%s", decryptedBytes)
	}

	return decryptedText, err
}

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetStringArrayHash(texts []string) string {
	hasher := md5.New()
	hasher.Write([]byte(strings.Join(texts, "")))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetSHA256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
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

func GetRawHotelId(hotelId string) string {
	if strings.Index(hotelId, "_") >= 0 {

		values := strings.Split(hotelId, "_")

		return values[1]
	}

	return hotelId
}

func GenerateCancellationPolicyShort(cancellations []*CancelPolicy, checkInDate string, nightlyPrice float32, total float32) string {
	From(cancellations).OrderByDescending(func(cxlPolicy interface{}) interface{} {
		return cxlPolicy.(*CancelPolicy).NumberDaysPrior
	}).ToSlice(&cancellations)

	arrivalDate, _ := time.Parse("2006-01-02", checkInDate)

	var cxlPolicy bytes.Buffer
	var lastDate = ""
	for i, cancelPolicy := range cancellations {
		if cxlPolicy.Len() != 0 {
			cxlPolicy.WriteString("<br />")
		}
		//fmt.Printf("cancelPolicy = %+v\n", cancelPolicy)
		if i == 0 {
			calcDate := GetCalculatedDate(cancelPolicy.NumberDaysPrior+1, arrivalDate)
			if calcDate != "" {
				cxlPolicy.WriteString("If you cancel this booking ")
				cxlPolicy.WriteString("from now ")
				cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
					calcDate,
					0.0))
			}
		} else {
			calcDate := GetCalculatedDate(cancelPolicy.NumberDaysPrior+1, arrivalDate)
			if calcDate != "" {
				cxlPolicy.WriteString("If you cancel this booking ")

				if lastDate == "" {
					cxlPolicy.WriteString("If you cancel this booking ")
					cxlPolicy.WriteString("from now ")
					cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
						calcDate,
						GetCalculatedPrice(cancellations[i-1], nightlyPrice, total)))
				} else {
					cxlPolicy.WriteString(fmt.Sprintf("from %[1]s until %[2]s, cancellation charge=$%.2f.",
						lastDate,
						calcDate,
						GetCalculatedPrice(cancellations[i-1], nightlyPrice, total)))
				}
			}
		}
		lastDate = GetCalculatedDate(cancelPolicy.NumberDaysPrior, arrivalDate)
		if i == len(cancellations)-1 {
			if lastDate != "" {
				if cxlPolicy.Len() != 0 {
					cxlPolicy.WriteString("<br />")
				}
				cxlPolicy.WriteString("If you cancel this booking ")

				cxlPolicy.WriteString(fmt.Sprintf("from %s onwards, cancellation charge $%.2f.",
					lastDate,
					GetCalculatedPrice(cancelPolicy, nightlyPrice, total)))
			}
		}
	}

	return cxlPolicy.String()
}

func GenerateCancellationPolicy(cancellations []*CancelPolicy) string {

	return ""
}

func GetFreeCancelDate(cancellations []*CancelPolicy, arrivalDate string) string {
	message := ""
	From(cancellations).OrderByDescending(func(cxlPolicy interface{}) interface{} {
		return cxlPolicy.(*CancelPolicy).NumberDaysPrior
	}).ToSlice(&cancellations)

	if cancellations != nil {
		freeCondition := cancellations[0]
		checkInDate, _ := time.Parse("2006-01-02", arrivalDate)
		freeCancellationDate := checkInDate.AddDate(0, 0, freeCondition.NumberDaysPrior*(-1)-1)
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

func GetCalculatedPrice(cancelPolicy *CancelPolicy, nightlyPrice float32, total float32) float32 {
	if cancelPolicy.NumberOfNights > 0 {
		return float32(cancelPolicy.NumberOfNights) * nightlyPrice
	} else if cancelPolicy.Percentage > 0.0 {
		return cancelPolicy.Percentage * total / 100.0
	}

	return 0.0
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
