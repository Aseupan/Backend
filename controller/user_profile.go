package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
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

		var primaryAddress bool

		var address model.Address
		err := db.Where("user_id = ? AND primary_address = ?", ID, true).First(&address).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// No rows with primary_address = true for the given user_id
				primaryAddress = true
			} else {
				// basic error
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			}
		} else {
			// There is a primary address for the given user_id
			primaryAddress = false
		}

		newAddres := model.Address{
			UserID:          ID.(uuid.UUID),
			Name:            input.Name,
			Phone:           input.Phone,
			Address:         input.Address,
			City:            input.City,
			State:           input.State,
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

	// edit one of user addresses by id
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
			State:           input.State,
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

	r.PATCH("/address/:id/primary", middleware.Authorization(), func(c *gin.Context) {
		addressID := c.Param("id")

		var address model.Address
		if err := db.Where("id = ?", addressID).First(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		address.PrimaryAddress = true
		if err := db.Save(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if err := db.Model(&model.Address{}).Where("user_id = ?", address.UserID).Where("id != ?", address.ID).Update("primary_address", false).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address updated", address)
	})

	r.DELETE("/address/:id", middleware.Authorization(), func(c *gin.Context) {
		addressID := c.Param("id")

		var address model.Address
		if err := db.Where("id = ?", addressID).First(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if err := db.Where("id = ?", addressID).Delete(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address deleted", nil)
	})
}
