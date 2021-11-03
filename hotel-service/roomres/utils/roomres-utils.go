package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"

	"bufio"
	"hash/crc32"
	"math"
	"net/url"

	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	//"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	//httperrors "roomres/errors"

	"github.com/go-errors/errors"
)

/**
var (
	encryptionKey []byte = []byte("RoomRes Encryption Key - Hn67ji6")
)
*/

const (
	LayoutYYYYMMDD      string = "2006-01-02"
	LayoutDDMMYYYY      string = "02-01-2006"
	LayoutDD_Month_YYYY string = "02 Jan 2006"
)

func ConvertDDMMYYYYToYYYYMMDD(value string) string {

	const (
		layoutFrom string = "02-Jan-2006"
		layoutTo   string = "2006-01-02"
	)

	dt, _ := time.Parse(layoutFrom, value)

	return dt.Format(layoutTo)
}

func ConvertDDMMYYYYToDate(value string) (time.Time, error) {

	const (
		layoutFrom string = "02-Jan-2006"
	)

	return time.Parse(layoutFrom, value)
}

func SwapKeyValue(mapKeys map[string]int, mapValues map[int]string) {

	for key, id := range mapKeys {
		mapValues[id] = key
	}
}

func CreateRootLogScopeUsingDefaults() ILog {
	return NewLog(
		LogSettings{
			AllowErrors:    true,
			AllowWarnings:  true,
			MaxAllowedInfo: GetMaxAllowedLogInfo(),
		},
		[]ILogProvider{NewConsoleLogProvider()}).StartLogScope(LogScopeRef{})
}

func GetMaxAllowedLogInfo() EventType {

	var eventType EventType = EventTypeInfo1

	logInfoLevel, _ := strconv.Atoi(os.Getenv("LOG_INFO_LEVEL"))

	if logInfoLevel == 1 {
		eventType = EventTypeInfo1
	} else if logInfoLevel == 2 {
		eventType = EventTypeInfo2
	} else if logInfoLevel == 3 {
		eventType = EventTypeInfo3
	}

	return eventType
}

func CompareDateYYMMDD(date1 time.Time, date2 time.Time) bool {
	return date1.Year() == date2.Year() &&
		date1.Month() == date2.Month() &&
		date1.Day() == date2.Day()
}

func SliceMaxLength(value string, maxLength int) string {
	if len(value) < maxLength {
		return value
	} else {
		return value[:maxLength]
	}
}

func CalculateLos(checkIn string, checkOut string) (int, error) {

	return CalculateLosLayout(checkIn, checkOut, LayoutYYYYMMDD)
}

func CalculateLosLayout(checkIn, checkOut, layout string) (int, error) {

	var checkInTime, checkOutTime time.Time
	var err error

	checkInTime, err = time.Parse(layout, checkIn)

	if err != nil {
		panic(err.Error())
		return 0, err
	}

	checkOutTime, err = time.Parse(layout, checkOut)

	if err != nil {
		panic(err.Error())
		return 0, err
	}

	return int(checkOutTime.Sub(checkInTime).Hours() / 24), nil
}

func Round(x, unit float32) float32 {
	return float32(int64(x/unit+0.5)) * unit
}

func Clone(a, b interface{}) {

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
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

func Index(s string, from int, a string) (p int) {
	p = strings.Index(s[from:], a)

	if p >= 0 /*found, set global pos*/ {
		p += from
	}

	return
}

func LastIndex(s string, from int, a string) (p int) {
	p = strings.LastIndex(s[:from+1], a)

	if p >= 0 /*found, set global pos*/ {
		p += len(a)
	}

	return
}

func IndexCI(s, a_0, a_1, a_2 string) (p int) {

	p = strings.Index(s, a_0)

	if p < 0 {

		p = strings.Index(s, a_1)

		if p < 0 {

			p = strings.Index(s, a_2)
		}
	}

	return p
}

func LastIndexCI(s, a_0, a_1, a_2 string) (p int) {

	p = strings.LastIndex(s, a_0)

	if p < 0 {

		p = strings.LastIndex(s, a_1)

		if p < 0 {

			p = strings.LastIndex(s, a_2)
		}
	}

	return p
}

func ReturnJSONP(element_name, callback_func_name string, json_data []byte) []byte {

	json_string := string(json_data)

	result_string := fmt.Sprintf("%s('%s', %s)", callback_func_name, element_name, json_string)

	return []byte(result_string)
}

func ReturnJSONPSimple(callback_func_name string, json_data []byte) []byte {

	json_string := string(json_data)

	result_string := fmt.Sprintf("%s(%s)", callback_func_name, json_string)

	return []byte(result_string)
}

func MultiReplaceAllUsingMap(s string, old_new_map map[string]string) string {

	var olds, news []string

	for old_, new_ := range old_new_map {

		olds = append(olds, old_)
		news = append(news, new_)
	}

	return MultiReplaceAll(s, olds, news)
}

func ReplaceRepeat(s string, old, new string) (result string) {

	result = s

	for {

		tmp := strings.Replace(result, old, new, -1)

		if tmp == result {
			break
		}

		result = tmp
	}

	return result
}

func MultiReplaceAll(s string, olds, news []string) (result string) {

	result = s

	for i, old := range olds {
		result = strings.Replace(result, old, news[i], -1)
	}

	return
}

func MultiReplaceAllExt(s string, olds, news []string, num []int) (result string) {

	result = s

	for i, old := range olds {
		result = strings.Replace(result, old, news[i], num[i])
	}

	return
}

func ReadLine(reader *bufio.Reader) (res bool, line string) {

	var (
		part   []byte
		err    error
		prefix bool
	)

	buffer := bytes.NewBuffer(make([]byte, 0))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {

			res = false
			line = ""

			break
		}

		buffer.Write(part)

		if !prefix {

			line = buffer.String()
			res = true

			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func ReadAll(reader *bufio.Reader) (text string) {

	text = ""

	for {

		res, line := ReadLine(reader)

		if !res {
			break
		}

		text = text + line
	}

	return
}

func AdjustDateYMD(t time.Time) (res time.Time) {
	y, mon, d := t.Date()

	res = time.Date(y, mon, d, 0, 0, 0, 0, time.UTC)

	return
}

func MakeWordsCapital(value string) string {

	newvalue := value

	for true {

		if strings.Index(newvalue, "  ") >= 0 {
			newvalue = strings.Replace(newvalue, "  ", " ", -1)
		} else {
			break
		}
	}

	items := strings.Split(strings.Trim(newvalue, " "), " ")

	res := ""

	for i, item := range items {

		if item == "" {
			continue
		}

		res += strings.ToUpper(item[0:1])

		if len(item) > 1 {
			res += item[1:len(item)]
		}

		if i < len(items)-1 {

			res += " "
		}
	}

	return res
}

func GetStringHash(value string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(value)))
}

func CreateKey(subkeys ...string) string {

	var subkeys_ []string = make([]string, len(subkeys))

	for i := 0; i < len(subkeys_); i++ {
		subkeys_[i] = strings.ToUpper(subkeys[i])
	}

	return strings.Join(subkeys_, "|")
}

func CreateQSParamsMap(qs string) (params map[string][]string, err_ error) {

	params = make(map[string][]string)

	qs_un, err := url.QueryUnescape(qs)

	if err != nil {
		err_ = err
		return
		panic(err.Error())
	}

	pairs := strings.Split(qs_un, "&")

	//fmt.Println(pairs)

	for _, pair := range pairs {

		items := strings.Split(pair, "=")

		if len(items) < 2 {
			continue
		}

		name := items[0]
		value := strings.Trim(items[1], " ")

		if _, ok := params[strings.ToUpper(name)]; ok {
			params[strings.ToUpper(name)] = append(params[strings.ToUpper(name)], value)
		} else {
			params[strings.ToUpper(name)] = []string{value}
		}
	}

	return
}

func IntToDate(value int) time.Time {

	a, b := value%10000, value%100

	return time.Date((value-a)/10000, time.Month((a-b)/100), b, 0, 0, 0, 0, time.UTC)
}

func DateToInt(value time.Time) int {
	return value.Year()*10000 + int(value.Month())*100 + value.Day()
}

func GetLos(checkin, checkout time.Time) int {
	return int(int64(checkout.Sub(checkin)) / (int64(time.Hour) * 24))
}

func MoneyToInt(value float64) int {

	return int(value * 100.0)
}

func IntToMoney(value int) float64 {

	return float64(value) / float64(100)
}

func FormatMoney(digits_after int, value float64) string {

	return fmt.Sprintf("%."+strconv.Itoa(digits_after)+"f", value)
}

func RoundOn(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func InterfaceToString(o interface{}) string {

	switch o.(type) {
	case string:
		return o.(string)
	case int32:
		return strconv.FormatInt(int64(o.(int32)), 32)
	case int64:
		return strconv.FormatInt(o.(int64), 64)
	case float64:
		return fmt.Sprintf("%.f", o.(float64))
	default:
		return "unknown"
	}
}

func CheckSignature(w http.ResponseWriter, r *http.Request, hbeSignature, method string, log ILog) bool {

	if hbeSignature == "" {
		return true
	}

	if r.Header.Get("Hbe-Signature") != hbeSignature {

		log.LogEvent(EventTypeInfo2,
			"StatusUnauthorized ",
			method,
		)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))

		return false
	}

	return true
}

func CreateRootLogScope() ILog {
	return NewLog(
		LogSettings{
			AllowErrors:    true,
			AllowWarnings:  true,
			MaxAllowedInfo: GetMaxAllowedLogInfo(),
		},
		[]ILogProvider{NewConsoleLogProvider()}).StartLogScope(LogScopeRef{})
}

// func HandleRequest(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	hbeSignature string,
// 	requestType interface{},
// 	actionName string,
// 	action func(request interface{}, logScope ILog) (response interface{}, err *httperrors.HttpError),
// ) {

// 	logScope := CreateRootLogScope()

// 	defer func(logScope ILog) {
// 		if err := recover(); err != nil {

// 			logScope.LogEvent(EventTypeError,
// 				actionName+" ",
// 				fmt.Sprintf("%s", errors.Wrap(err, 2).ErrorStack()),
// 			)

// 			http.Error(w, "Internal Server Error", 500)

// 			return

// 		}
// 	}(logScope)

// 	if !CheckSignature(w, r, hbeSignature, actionName, logScope) {
// 		return
// 	}

// 	if r.Body == nil {
// 		http.Error(w, "Please send a request body", 400)
// 		return
// 	}

// 	bodyBuf := &bytes.Buffer{}
// 	bodyBuf.ReadFrom(r.Body)

// 	{
// 		s := bodyBuf.String()

// 		logScope.LogEvent(EventTypeInfo2,
// 			actionName+" request",
// 			fmt.Sprintf("%s - %s - requestdata: >>%s<<", r.RemoteAddr, r.Header.Get("X-Forwarded-For"), s),
// 		)
// 	}

// 	if requestType != nil {

// 		err := json.NewDecoder(bytes.NewReader(bodyBuf.Bytes())).Decode(requestType)
// 		if err != nil {
// 			http.Error(w, err.Error(), 400)
// 			return
// 		}
// 	}

// 	response, httpErr := action(requestType, logScope)

// 	if httpErr != nil {
// 		http.Error(w, httpErr.String(), int(httpErr.Code))
// 	}

// 	w.Header().Add("Content-Type", "application/json")

// 	outputBuffer := new(bytes.Buffer)

// 	json.NewEncoder(outputBuffer).Encode(response)

// 	w.Write(outputBuffer.Bytes())
// }
