package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreditStore(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/credit-store")
	// get all
	r.GET("/all", middleware.Authorization(), func(c *gin.Context) {
		var scores []model.CreditStore
		if err := db.Find(&scores).Error; err != nil {
			log.Println("di sini error mencari data")
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "All credit store", scores)
	})

	// add to cart
	r.POST("/add-to-cart", middleware.Authorization(), func(c *gin.Context) {
		var input model.CreditStoreWalletInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var credit model.CreditStore
		if err := db.Where("id = ?", input.ID).First(&credit).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		ID, _ := c.Get("id")
		userID, ok := ID.(uuid.UUID)
		if !ok {
			utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
			return
		}

		addToCart := model.CreditStoreWallet{
			UserID: userID,
			Points: credit.Points,
			Price:  credit.Price,
		}

		if err := db.Create(&addToCart).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Add to cart", addToCart)
	})

	// remove item from cart
	r.DELETE("/remove-from-cart", middleware.Authorization(), func(c *gin.Context) {
		var input model.CreditStoreWalletInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		userID, ok := ID.(uuid.UUID)
		if !ok {
			utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
			return
		}

		var cart model.CreditStoreWallet
		if err := db.Where("id = ?", input.ID).Where("user_id = ?", userID).First(&cart).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		if err := db.Delete(&cart).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Remove from cart", cart)
	})
}
