package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserProfile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user")
	// user profile
	r.GET("/profile", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		var user model.User
		if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "User profile", user)
	})

	// list of user addresses
	r.GET("/addresses", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		var user model.User
		if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var addresses []model.Address
		if err := db.Where("user_id = ?", ID).Find(&addresses).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "User addresses", addresses)
	})

	// add new address
	r.POST("/address", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var input model.AddressInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var maxIndex int
		err := db.Model(&model.Address{}).Select("MAX(`index`)").Where("user_id = ?", ID).Scan(&maxIndex).Error
		if err != nil {
			log.Printf("error max index: %v", maxIndex)
			maxIndex = 0
		}

		var primaryAddress bool
		if maxIndex == 0 {
			primaryAddress = true
		} else {
			primaryAddress = false
		}

		newAddres := model.Address{
			UserID:          ID.(uuid.UUID),
			Name:            input.Name,
			Phone:           input.Phone,
			Address:         input.Address,
			City:            input.City,
			Disctrict:       input.Disctrict,
			ZipCode:         input.ZipCode,
			DetailedAddress: input.DetailedAddress,
			PrimaryAddress:  primaryAddress,
			CreatedAt:       time.Now(),
		}

		if err := db.Create(&newAddres).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "New address added", newAddres)
	})

	// edit one of user addresses by index
	r.PATCH("/address/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		addressID := c.Param("id")

		var address model.Address
		if err := db.Where("user_id = ?", ID).Where("id = ?", addressID).First(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var input model.AddressInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		updatedAddres := model.Address{
			Name:            input.Name,
			Phone:           input.Phone,
			Address:         input.Address,
			City:            input.City,
			Disctrict:       input.Disctrict,
			ZipCode:         input.ZipCode,
			DetailedAddress: input.DetailedAddress,
			PrimaryAddress:  input.PrimaryAddress,
			UpdatedAt:       time.Now(),
		}

		if err := db.Where("user_id = ?", ID).Where("id = ?", addressID).Model(&address).Updates(updatedAddres).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address updated", updatedAddres)
	})
}
