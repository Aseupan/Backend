package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Rewards(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/rewards")
	// user home
	r.GET("/view-all", middleware.Authorization(), func(c *gin.Context) {
		var rewards []model.Rewards
		if err := db.Find(&rewards).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}
		utils.HttpRespSuccess(c, http.StatusOK, "get all rewards", rewards)
	})

	r.POST("/purchase-reward", middleware.Authorization(), func(c *gin.Context) {
		var input model.RewardsInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var reward model.Rewards
		if err := db.Where("id = ?", input.RewardID).First(&reward).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if user.Point < reward.Points {
				utils.HttpRespFailed(c, http.StatusNotFound, "Not enough points")
				return
			}

			user.Point -= reward.Points
			reward.Quantity -= 1

			if err := db.Save(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		} else if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if company.Point < reward.Points {
				utils.HttpRespFailed(c, http.StatusNotFound, "Not enough points")
				return
			}

			company.Point -= reward.Points
			reward.Quantity -= 1

			if err := db.Save(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		if err := db.Save(&reward).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Purchase reward success", nil)
	})
}
