package controller

import (
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CompanyProfile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/company")
	// company profile
	r.GET("/profile", func(c *gin.Context) {
		ID, _ := c.Get("id")
		var company model.Company
		if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Company profile", company)
	})
}
