package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StringToInteger(input string, c *gin.Context) int {
	converted, err := strconv.Atoi(input)
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
