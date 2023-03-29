package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Campaign(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/campaign")

	// big party / company
	r.POST("/create", middleware.Authorization(), func(c *gin.Context) {
		// ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// name := c.PostForm("name")

		// description := c.PostForm("description")

		// target := c.PostForm("target")

		// area := c.PostForm("area")

		// startDate := c.PostForm("startdate")

		// endDate := c.PostForm("enddate")

		// urgent := c.PostForm("urgent")

		// foodType := c.PostFormArray("type")

		// thumbnail1, err := c.FormFile("thumbnail1")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Thumbnail1 is required")
		// 	return
		// }

		// thumbnail2, err := c.FormFile("thumbnail2")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Thumbnail1 is required")
		// 	return
		// }

		// thumbnail3, err := c.FormFile("thumbnail3")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Thumbnail1 is required")
		// 	return
		// }

		// thumbnail4, err := c.FormFile("thumbnail4")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Thumbnail1 is required")
		// 	return
		// }

		// thumbnail5, err := c.FormFile("thumbnail5")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Thumbnail1 is required")
		// 	return
		// }

		// var newCampaign model.Campaign
		// newCampaign.Name = name
		// newCampaign.Description = description
		// newCampaign.Target =
		// newCampaign.Area =
		// newCampaign.StartDate =
		// newCampaign.EndDate =
		// newCampaign.Urgent =
		// newCampaign.Type =

	})

	// user
	r.GET("/all", middleware.Authorization(), func(c *gin.Context) {
		var campaigns []model.Campaign
		if res := db.Find(&campaigns); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", nil)
	})

	r.GET("/detail/:id", middleware.Authorization(), func(c *gin.Context) {
		// strType, _ := c.Get("type")

		// if strType != "user" {
		// 	utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
		// 	return
		// }

		id := c.Param("id")

		var campaign model.Campaign
		if res := db.Where("id = ?", id).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", campaign)
	})
}
