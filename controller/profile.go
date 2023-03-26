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

	// add new address
	r.POST("/address", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		var input model.AddressInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var primaryAddress bool

		var address model.Address
		var newAddress model.Address

		if strType == "company" {
			err := db.Where("company_id = ? AND primary_address = ?", ID, true).First(&address).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					primaryAddress = true
				} else {
					utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				}
			} else {
				primaryAddress = false
			}

			newAddress = model.Address{
				CompanyID:       ID.(uuid.UUID),
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
		} else if strType == "user" {
			err := db.Where("user_id = ? AND primary_address = ?", ID, true).First(&address).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					primaryAddress = true
				} else {
					utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				}
			} else {
				primaryAddress = false
			}

			newAddress = model.Address{
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
		}

		if err := db.Create(&newAddress).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "New address added", newAddress)
	})

	// edit one of addresses by id
	r.PATCH("/address/:id", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")
		addressID := c.Param("id")

		var input model.AddressInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var address model.Address
		var updatedAddress model.Address

		if strType == "company" {
			if err := db.Where("company_id = ?", ID).Where("id = ?", addressID).First(&address).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updatedAddress = model.Address{
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

			if err := db.Where("company_id = ?", ID).Where("id = ?", addressID).Model(&address).Updates(updatedAddress).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
		} else if strType == "user" {
			if err := db.Where("user_id = ?", ID).Where("id = ?", addressID).First(&address).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updatedAddress = model.Address{
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

			if err := db.Where("user_id = ?", ID).Where("id = ?", addressID).Model(&address).Updates(updatedAddress).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address updated", updatedAddress)
	})

	r.PATCH("/address/:id/primary", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")
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

		if strType == "company" {
			if err := db.Model(&model.Address{}).Where("company_id = ?", address.UserID).Where("id != ?", address.ID).Update("primary_address", false).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
		} else if strType == "user" {
			if err := db.Model(&model.Address{}).Where("user_id = ?", address.UserID).Where("id != ?", address.ID).Update("primary_address", false).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address updated", address)
	})

	r.DELETE("/address/:id", middleware.Authorization(), func(c *gin.Context) {
		addressID := c.Param("id")

		var address model.Address

		if err := db.Where("id = ?", addressID).Delete(&address).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Address deleted", nil)
	})
}
