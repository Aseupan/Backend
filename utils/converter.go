package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	geo "github.com/kellydunn/golang-geo"
)

func LocationToKM(userLatitude float64, userLongitude float64, campaignLatitude float64, campaignLongitude float64) float64 {
	user := geo.NewPoint(userLatitude, userLongitude)

	campaign := geo.NewPoint(campaignLatitude, campaignLongitude)

	distance := user.GreatCircleDistance(campaign)

	return distance
}

func StringToInteger(input string, c *gin.Context) int {
	converted, err := strconv.Atoi(input)
	if err != nil {
		HttpRespFailed(c, http.StatusBadRequest, "Invalid input")
		return 0
	}
	return converted
}

func StringToFloat(input string, c *gin.Context) float64 {
	converted, err := strconv.ParseFloat(input, 64)
	if err != nil {
		HttpRespFailed(c, http.StatusBadRequest, "Invalid input")
		return 0
	}
	return converted
}

func StringToUint(input string, c *gin.Context) uint {
	converted, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		HttpRespFailed(c, http.StatusBadRequest, "Invalid input")
		return 0
	}
	return uint(converted)
}
