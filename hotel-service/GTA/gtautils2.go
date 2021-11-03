package gta

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/ahmetb/go-linq"
)

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

func makeOccupancyString(Adults int, Children int) string {
	Occupancy := "Single"
	if Children == 0 {
		if Adults == 2 {
			Occupancy = "Double"
		} else if Adults == 3 {
			Occupancy = "Triple"
		} else if Adults == 4 {
			Occupancy = "Quad"
		}
		//Todo validation
	} else {
		if Adults == 1 && Children == 1 {
			Occupancy = "SGL+1CH"
		} else if Adults == 1 && Children == 2 {
			Occupancy = "SGL+2CH"
		} else if Adults == 1 && Children == 3 {
			Occupancy = "SGL+3CH"
		} else if Adults == 2 && Children == 1 {
			Occupancy = "DBL+1CH"
		} else if Adults == 2 && Children == 2 {
			Occupancy = "DBL+2CH"
		} else if Adults == 3 && Children == 1 {
			Occupancy = "TPL+1CH"
		}
	}

	return Occupancy
}

func GetRawHotelId(hotelId string) string {
	if strings.Index(hotelId, "_") >= 0 {

		values := strings.Split(hotelId, "_")

		return values[1]
	}

	return hotelId
}

func ConvertDateFormat(dateStr string, sourceFormat string, destFormat string) string {
	t, _ := time.Parse(sourceFormat, dateStr)
	return t.Format(destFormat)
}

func GenerateCancellationPolicyShort(chargeConditions []*ChargeCondition) string {
	var chargeCondition *ChargeCondition = nil
	for _, c := range chargeConditions {
		if c.Type == "cancellation" {
			chargeCondition = c
			break
		}
	}
	if chargeCondition != nil {
		var cxlPolicy bytes.Buffer
		From(chargeCondition.Conditions).OrderBy(func(t interface{}) interface{} {
			return t.(*Condition).FromDate
		}).ToSlice(&chargeCondition.Conditions)
		for i, cancelPolicy := range chargeCondition.Conditions {
			if cxlPolicy.Len() != 0 {

				cxlPolicy.WriteString("<br />")
			}
			//fmt.Printf("cancelPolicy = %+v\n", cancelPolicy)
			cxlPolicy.WriteString("If you cancel this booking ")
			if cancelPolicy.Charge == false {
				cxlPolicy.WriteString("from now ")
				cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
					ConvertDateFormat(cancelPolicy.FromDate, "2006-01-02", "02/01/2006"),
					cancelPolicy.ChargeAmount))
			} else {
				cxlPolicy.WriteString(fmt.Sprintf("from %[1]s until %[2]s, cancellation charge=$%.2f.",
					ConvertDateFormat(cancelPolicy.ToDate, "2006-01-02", "02/01/2006"),
					ConvertDateFormat(cancelPolicy.FromDate, "2006-01-02", "02/01/2006"),
					cancelPolicy.ChargeAmount))
			}

			if i == len(chargeCondition.Conditions)-1 {

				cxlPolicy.WriteString(fmt.Sprintf("<br />No show charge=$%.2f.",
					cancelPolicy.ChargeAmount))
			}
		}

		return cxlPolicy.String()
	}

	return ""
}

func GetFreeCancelDate(chargeConditions []*ChargeCondition) string {
	var chargeCondition *ChargeCondition = nil
	var message = ""
	for _, c := range chargeConditions {
		if c.Type == "cancellation" {
			chargeCondition = c
			break
		}
	}
	if chargeCondition != nil {
		var freeCondition *Condition = nil
		for _, c := range chargeCondition.Conditions {
			if !c.Charge {
				freeCondition = c
				break
			}
		}

		if freeCondition != nil {
			message = fmt.Sprintf("FREE Cancellation until %s", ConvertDateFormat(freeCondition.FromDate, "2006-01-02", "02/01/2006"))
		}
	}

	return message
}

func checkElementInArray(list []string, ele string) bool {
	for _, val := range list {
		if val == ele {
			return true
		}
	}

	return false
}
