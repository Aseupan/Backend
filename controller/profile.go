package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/join")

	// get user profile
	r.GET("/profile", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")
		ID, _ := c.Get("id")

		if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Company profile", company)
		} else if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "User profile", user)
		}
	})

	// update user profile
	r.PATCH("/profile", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")
		ID, _ := c.Get("id")

		if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			var input model.CompanyUpdateProfileInput
			if err := c.BindJSON(&input); err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			if err := db.Model(&company).Updates(input).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Company profile updated", company)

		} else if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			var input model.UserUpdateProfileInput
			if err := c.BindJSON(&input); err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			if err := db.Model(&user).Updates(input).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "User profile updated", user)
		}
	})

	// list of user addresses
	r.GET("/addresses", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")
		ID, _ := c.Get("id")

		var addresses []model.Address

		if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("company_id = ?", ID).Find(&addresses).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("user_id = ?", ID).Find(&addresses).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "addresses", addresses)
	})
}
