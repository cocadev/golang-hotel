package jac

import (
	"bytes"
	"encoding/json"
	"fmt"

	hbecommon "../roomres/hbe/common"

	//roomresutils "roomres/utils"
	"strings"
	"time"
)

const (
	JacHotelPrefix string = "jac_"
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

func (m *HotelIdTransformer) ExtractJacHotelId(hotelId string) string {

	supplierHotelId := roomresutils.ExtractSupplierHotelId(JacHotelPrefix, hotelId)

	items := strings.Split(supplierHotelId, "-")

	if len(items) > 1 {
		supplierHotelId = items[0]
		region := items[1]

		m.HotelIdInfoItems[supplierHotelId] = &HotelIdInfo{SupplierHotelId: supplierHotelId, Region: region}
	} else {

		m.HotelIdInfoItems[supplierHotelId] = &HotelIdInfo{SupplierHotelId: supplierHotelId, Region: ""}
	}

	return supplierHotelId
}

func (m *HotelIdTransformer) GenerateJacHotelId(supplierHotelId string) string {

	info, ok := m.HotelIdInfoItems[supplierHotelId]

	region := ""

	if ok {
		region = info.Region
	}

	if region == "" {
		return roomresutils.GenerateSupplierHotelId(JacHotelPrefix, supplierHotelId)
	} else {
		return fmt.Sprintf("%s-%s", roomresutils.GenerateSupplierHotelId(JacHotelPrefix, supplierHotelId), region)
	}
}

func CompareRoomRefs(roomRefA, roomRefB string) bool {

	return roomRefA == roomRefB
}

func AssignRequestedRoomsToSearchResult(searchRequest *SearchRequest, searchResponse *SearchResponse) {

	for _, propertyResult := range searchResponse.PropertyResults {

		for _, roomType := range propertyResult.RoomTypes {
			roomType.RequestedRoom = searchRequest.SearchDetails.RoomRequests[roomType.Seq-1]
		}

	}

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

func GenerateRoomTypeNotes(roomType *RoomType) string {

	var notes bytes.Buffer

	for _, erratum := range roomType.Errata {

		if notes.Len() > 0 {
			notes.WriteString("<br />")
		}

		notes.WriteString("<b>")
		notes.WriteString(erratum.Subject)
		notes.WriteString("</b><br />")
		notes.WriteString(erratum.Description)
	}

	return notes.String()
}

func IsNonRefundableCancellation(cancellations []*Cancellation) bool {

	return false //return GenerateCancellationPolicyShort(checkIn, cancellations) == ""
}

func IsNonRefundableRoom(roomTypeName string) bool {
	return strings.Index(strings.ToLower(roomTypeName), "non-refundable") >= 0
}

func AdjustCancellationPolicyStartDate(cxlPolicyStartDate time.Time) time.Time {
	return cxlPolicyStartDate.AddDate(0, 0, -1)
	//return cxlPolicyStartDate
}

func GenerateCancellationPolicy(commissionProvider hbecommon.ICommissionProvider, cancellations []*Cancellation, mandatoryFee float32) string {

	var cxlPolicy bytes.Buffer

	currentDate := roomresutils.CurrentDate()

	for i, cancellation := range cancellations {

		if cxlPolicy.Len() != 0 {

			cxlPolicy.WriteString("<br />")
		}

		var startDate, endDate time.Time
		var err error

		startDate, err = ConvertToDate(cancellation.StartDate)
		if err != nil {
			panic(err.Error())
		}

		endDate, err = ConvertToDate(cancellation.EndDate)
		if err != nil {
			panic(err.Error())
		}

		cxlPolicy.WriteString("If you cancel this booking ")

		if startDate.Equal(currentDate) || startDate.Before(currentDate) {

			cxlPolicy.WriteString("from now ")
			cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
				AdjustCancellationPolicyStartDate(endDate).Format(DDMMYYYYLayout),
				CalculateCancellationAmount(commissionProvider, cancellation.Penalty, mandatoryFee)))

		} else if cancellation.EndDate == CxlEndDate {

			cxlPolicy.WriteString(fmt.Sprintf("from %s onwards, cancellation charge $%.2f.",
				AdjustCancellationPolicyStartDate(startDate).Format(DDMMYYYYLayout),
				CalculateCancellationAmount(commissionProvider, cancellation.Penalty, mandatoryFee)))
		} else {

			cxlPolicy.WriteString(fmt.Sprintf("from %[1]s until %[2]s, cancellation charge=$%.2f.",
				AdjustCancellationPolicyStartDate(startDate).Format(DDMMYYYYLayout),
				AdjustCancellationPolicyStartDate(endDate).Format(DDMMYYYYLayout),
				CalculateCancellationAmount(commissionProvider, cancellation.Penalty, mandatoryFee)))

		}

		if i == len(cancellations)-1 {

			cxlPolicy.WriteString(fmt.Sprintf("<br />No show charge=$%.2f.",
				CalculateCancellationAmount(commissionProvider, cancellation.Penalty, mandatoryFee)))
		}

	}

	return cxlPolicy.String()
}

func CalculateCancellationAmount(commissionProvider hbecommon.ICommissionProvider, amount, mandatoryFee float32) float32 {

	if amount > 0 {
		return commissionProvider.CalculateSellRate(amount + mandatoryFee)
	} else {
		return 0
	}
}

func GenerateCancellationPolicyShort(cancellations []*Cancellation) string {

	var cxlPolicy bytes.Buffer

	var lastFreeCxl *Cancellation

	for _, cancellation := range cancellations {

		startDate, err := ConvertToDate(cancellation.StartDate)

		if err != nil {
			panic(err.Error())
		}

		if cancellation.Penalty == 0 {
			lastFreeCxl = cancellation
		}

		if cancellation.Penalty > 0 {

			var date time.Time

			if lastFreeCxl == nil {
				date = startDate.AddDate(0, 0, -1)
			} else {
				var err error
				date, err = ConvertToDate(lastFreeCxl.EndDate)

				if err != nil {
					panic(err.Error())
				}
			}

			cxlPolicy.WriteString(
				fmt.Sprintf(
					"FREE Cancellation until %s",
					AdjustCancellationPolicyStartDate(date).Format(DDMMYYYYLayout)))

			break
		}
	}

	return cxlPolicy.String()
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
