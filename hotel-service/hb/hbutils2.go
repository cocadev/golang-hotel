package hb

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	. "github.com/ahmetb/go-linq"
)

func GenerateCancellationPolicyShort(cancelPolicies []*CancelationPolicy) string {
	var cxlPolicy bytes.Buffer
	From(cancelPolicies).OrderBy(func(t interface{}) interface{} {
		return t.(*CancelationPolicy).From
	}).ToSlice(&cancelPolicies)
	lastDate := ""
	for i, cancelPolicy := range cancelPolicies {
		if cxlPolicy.Len() != 0 {

			cxlPolicy.WriteString("<br />")
		}

		cxlPolicy.WriteString("If you cancel this booking ")

		if lastDate == "" {
			fd, _ := time.Parse("2006-01-02T15:04:05-07:00", cancelPolicy.From)
			td := fd.AddDate(0, 0, -1)
			lastDate = fd.Format("02/01/2006")
			cxlPolicy.WriteString("from now ")
			cxlPolicy.WriteString(fmt.Sprintf(" until %s, cancellation charge=$%.2f.",
				td.Format("02/01/2006"),
				0.0))
		} else {
			fd, _ := time.Parse("02/01/2006", lastDate)
			fd2, _ := time.Parse("2006-01-02T15:04:05-07:00", cancelPolicy.From)
			td := fd2.AddDate(0, 0, -1)
			lastDate = fd2.Format("02/01/2006")

			cxlPolicy.WriteString(fmt.Sprintf("from %[1]s until %[2]s, cancellation charge=$%.2f.",
				fd.Format("02/01/2006"),
				td.Format("02/01/2006"),
				ConvertStringToFloat32(cancelPolicy.Amount)))
		}

		if i == len(cancelPolicies)-1 {
			if cxlPolicy.Len() != 0 {
				cxlPolicy.WriteString("<br />")
			}
			cxlPolicy.WriteString(fmt.Sprintf("from %s onwards, cancellation charge $%.2f.",
				lastDate,
				ConvertStringToFloat32(cancelPolicy.Amount)))
		}
	}

	return cxlPolicy.String()
}

func GetFreeCancelDate(cancelPolicies []*CancelationPolicy) string {
	message := ""
	From(cancelPolicies).OrderBy(func(t interface{}) interface{} {
		return t.(*CancelationPolicy).From
	}).ToSlice(&cancelPolicies)

	if cancelPolicies != nil {
		freeCondition := cancelPolicies[0]
		fd, _ := time.Parse("2006-01-02T15:04:05-07:00", freeCondition.From)
		fd = fd.AddDate(0, 0, -1)
		lastDate := fd.Format("02/01/2006")

		if freeCondition != nil {
			message = fmt.Sprintf("FREE Cancellation until %s", lastDate)
		}
	}

	return message
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
