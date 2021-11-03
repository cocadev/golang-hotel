package gta

import (
	"strings"
)

func CreateDestinationGroups(hotelIds []string) []*DestinationGroup {

	groupedCodes := map[string][]string{}

	for _, hotelId := range hotelIds {

		codes := strings.Split(hotelId, ";")

		if _, ok := groupedCodes[codes[0]]; !ok {

			groupedCodes[codes[0]] = []string{}
		}

		groupedCodes[codes[0]] = append(groupedCodes[codes[0]], codes[1])
	}

	destinationGroups := []*DestinationGroup{}

	for destinationCode, hotelIds := range groupedCodes {

		destinationGroup := &DestinationGroup{DestinationCode: destinationCode, HotelIds: hotelIds}

		destinationGroups = append(destinationGroups, destinationGroup)
	}

	return destinationGroups
}

func SplitDestinationGroups(hotelsPerBatch int, destinationGroups []*DestinationGroup) []*DestinationGroup {

	adjustedDestinationGroups := []*DestinationGroup{}

	for _, destinationGroup := range destinationGroups {

		if len(destinationGroup.HotelIds) > hotelsPerBatch {

			adjustedDestinationGroup := &DestinationGroup{DestinationCode: destinationGroup.DestinationCode, HotelIds: []string{}}

			for _, hotelId := range destinationGroup.HotelIds {

				if len(adjustedDestinationGroup.HotelIds) >= hotelsPerBatch {

					adjustedDestinationGroups = append(adjustedDestinationGroups, adjustedDestinationGroup)
					adjustedDestinationGroup = &DestinationGroup{DestinationCode: destinationGroup.DestinationCode, HotelIds: []string{}}
				}

				adjustedDestinationGroup.HotelIds = append(adjustedDestinationGroup.HotelIds, hotelId)
			}

			if len(adjustedDestinationGroup.HotelIds) > 0 {
				adjustedDestinationGroups = append(adjustedDestinationGroups, adjustedDestinationGroup)
			}

		} else {
			adjustedDestinationGroups = append(adjustedDestinationGroups, destinationGroup)
		}
	}

	return adjustedDestinationGroups
}
