package hb

import (
	"fmt"
	"strings"
	"time"

	hbecommon "../roomres/hbe/common"
	roomresutils "../roomres/utils"
)

func GenerateSignature(apiKey, secretKey string) string {

	return roomresutils.GetSHA256Hash(fmt.Sprintf("%s%s%d", apiKey, secretKey, time.Now().Unix()))
}

type RoomRate struct {
	Room            *Room
	Rate            *Rate
	OccupanciesFits int
}

func FindBestRoomRate(roomRates []*RoomRate) *RoomRate {

	var bestRoomRate *RoomRate

	for i, roomRate := range roomRates {

		if i == 0 {
			bestRoomRate = roomRate
		} else {
			if bestRoomRate.Rate.NetValue > roomRate.Rate.NetValue {
				bestRoomRate = roomRate
			}
		}
	}

	return bestRoomRate
}

func FindBestRoomRateCombinations(roomRequests []*hbecommon.RoomRequest, rooms []*Room) []*RoomRate {

	roomRates := map[string]*RoomRate{}

	for _, room := range rooms {

		for _, rate := range room.Rates {

			if strings.ToUpper(rate.RateType) != "BOOKABLE" && strings.ToUpper(rate.RateType) != "RECHECK" {
				continue
			}

			if strings.ToUpper(rate.RateType) == "RECHECK" && len(roomRequests) > 1 {
				continue
			}

			rate.NetValue = rate.GetNet()

			if roomRate, ok := roomRates[rate.GetOccupancyHash()]; !ok {

				roomRates[rate.GetOccupancyHash()] = &RoomRate{Room: room, Rate: rate, OccupanciesFits: 0}
			} else {

				if roomRate.Rate.NetValue > rate.NetValue {
					roomRate.Room = room
					roomRate.Rate = rate
				}
			}
		}
	}

	for _, roomRequest := range roomRequests {
		for _, roomRate := range roomRates {

			if roomRate.Rate.Adults >= roomRequest.Adults && roomRate.Rate.Children >= roomRequest.Children {
				roomRate.OccupanciesFits++
			}

		}
	}

	validCombinations := []*RoomRate{}

	for _, roomRate := range roomRates {
		if roomRate.OccupanciesFits == len(roomRequests) {
			validCombinations = append(validCombinations, roomRate)
		}
	}

	return validCombinations
}
