package utils

import (
	"encoding/json"
	"fmt"
	"gsc/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LocationToKM(c *gin.Context, userLatitude, userLongitude, campaignLatitude, campaignLongitude string) string {
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json?"
	params := url.Values{}
	params.Set("origins", userLatitude+","+userLongitude)
	params.Set("destinations", campaignLatitude+","+campaignLongitude)
	params.Set("mode", "driving")
	params.Set("key", os.Getenv("GOOGLE_API_KEY"))

	requestURL := baseURL + params.Encode()

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "0"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "0"
	}

	var distanceMatrixResponse model.DistanceMatrixResponse
	err = json.Unmarshal(body, &distanceMatrixResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "0"
	}

	distanceStr := distanceMatrixResponse.Rows[0].Elements[0].Distance.Text

	return distanceStr
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
