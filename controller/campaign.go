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
		strType, _ := c.Get("type")

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// name, err := c.FormFile("name")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Name is required")
		// 	return
		// }

		// description, err := c.FormFile("description")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Description is required")
		// 	return
		// }

		// target, err := c.FormFile("target")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Target is required")
		// 	return
		// }

		// area, err := c.FormFile("area")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Area is required")
		// 	return
		// }

		// startdate, err := c.FormFile("startdate")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Startdate is required")
		// 	return
		// }

		// enddate, err := c.FormFile("enddate")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Enddate is required")
		// 	return
		// }

		// urgent, err := c.FormFile("urgent")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Urgent is required")
		// 	return
		// }

		// foodType, err := c.FormFile("type")
		// if err != nil {
		// 	utils.HttpRespFailed(c, http.StatusBadRequest, "Type is required")
		// 	return
		// }

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
