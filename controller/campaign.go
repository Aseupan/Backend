package controller

import (
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Campaign(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/campaign")
	r.GET("/all", func(c *gin.Context) {
		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", nil)
	})
}
