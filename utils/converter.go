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
