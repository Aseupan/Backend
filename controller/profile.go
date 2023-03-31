package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

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

	// update user profile picture
	r.PATCH("/profile/picture", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		SupaBaseClient := supabasestorageuploader.NewSupabaseClient(
			"https://flldkbhntqqaiflpxlhg.supabase.co",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZsbGRrYmhudHFxYWlmbHB4bGhnIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3NzU4Njk4OCwiZXhwIjoxOTkzMTYyOTg4fQ.CezKv4eOdEOyPEnVCqp3i0rNRLpz4MJOgL2GvM74QtQ",
			"photo",
			"",
		)

		if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			photo, err := c.FormFile("pp")
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
			link, err := SupaBaseClient.Upload(photo)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			if err := db.Model(&company).Update("company_picture", link).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Company profile picture updated", company)

		} else if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			photo, err := c.FormFile("pp")
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
			link, err := SupaBaseClient.Upload(photo)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			if err := db.Model(&user).Update("profile_picture", link).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "User profile picture updated", user)
		}
	})
}
