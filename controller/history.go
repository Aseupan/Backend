package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func History(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/history")
	// user home
	r.GET("/ongoing", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Not authorized")
		}

		userID, _ := c.Get("id")
		var ongoing []model.UserPersonalDonation
		if err := db.Where("user_id = ?", userID).Preload("Campaign").Find(&ongoing).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "get all ongoing", ongoing)
	})

	r.GET("/completed", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Not authorized")
		}

		userID, _ := c.Get("id")
		var completed []model.History
		if err := db.Where("user_id = ?", userID).Find(&completed).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "get all completed", completed)
	})
}
