package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

func CreditStore(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/credit-store")
	// get all
	r.GET("/all", middleware.Authorization(), func(c *gin.Context) {
		var scores []model.CreditStore
		if err := db.Find(&scores).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "All credit store", scores)
	})

	// view cart
	r.GET("/view-cart", middleware.Authorization(), func(c *gin.Context) {
		var total int
		var totalPoints int

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		var cart []model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			if err := db.Where("company_id = ?", companyID).Find(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			if err := db.Where("user_id = ?", userID).Find(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		for _, v := range cart {
			total += v.Price
			totalPoints += v.Points
		}

		utils.HttpRespSuccess(c, http.StatusOK, "View cart", gin.H{
			"total":       total,
			"totalPoints": totalPoints,
			"cart":        cart,
		})
	})

	// add to cart
	r.POST("/add-to-cart", middleware.Authorization(), func(c *gin.Context) {
		var input model.CreditStoreCartInput
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
		strType, _ := c.Get("type")

		var addToCart model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			parsedID, err := strconv.ParseUint(strconv.Itoa(input.ID), 10, 0)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			addToCart = model.CreditStoreCart{
				CompanyID:     companyID,
				CreditStoreID: uint(parsedID),
				Points:        credit.Points,
				Price:         credit.Price,
				Quantity:      1,
			}

			var isExist model.CreditStoreCart
			if err := db.Where("company_id = ? ", companyID).Where("credit_store_id = ?", input.ID).First(&isExist).Error; err == nil {
				isExist.Points += credit.Points
				isExist.Price += credit.Price
				isExist.Quantity += 1
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", input.ID).Save(&isExist).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusBadGateway, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Update cart", isExist)
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			parsedID, err := strconv.ParseUint(strconv.Itoa(input.ID), 10, 0)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			addToCart = model.CreditStoreCart{
				UserID:        userID,
				CreditStoreID: uint(parsedID),
				Points:        credit.Points,
				Price:         credit.Price,
				Quantity:      1,
			}

			var isExist model.CreditStoreCart
			if err := db.Where("user_id = ? ", userID).Where("credit_store_id = ?", input.ID).First(&isExist).Error; err == nil {
				isExist.Points += credit.Points
				isExist.Price += credit.Price
				isExist.Quantity += 1
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", input.ID).Save(&isExist).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusBadGateway, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "added to cart", isExist)
				return
			}
		}

		if err := db.Create(&addToCart).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Add to cart", addToCart)
	})

	// add 1 amount by id
	r.POST("/add-amount/:itemID", middleware.Authorization(), func(c *gin.Context) {
		itemID := c.Param("itemID")

		var credit model.CreditStore
		if err := db.Where("id = ?", itemID).First(&credit).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType := c.GetString("type")

		var updated model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("company_id = ?", companyID).First(&updated).Error; err != nil {
				addToCart := model.CreditStoreCart{
					CompanyID:     companyID,
					CreditStoreID: credit.ID,
					Points:        credit.Points,
					Price:         credit.Price,
					Quantity:      1,
				}

				if err := db.Create(&addToCart).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "added to cart", addToCart)
				return
			}

			updated.Points += credit.Points
			updated.Price += credit.Price
			updated.Quantity += 1

			if err := db.Where("credit_store_id = ?", itemID).Where("company_id = ?", companyID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Update cart", updated)
			return

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("user_id = ?", userID).First(&updated).Error; err != nil {
				addToCart := model.CreditStoreCart{
					UserID:        userID,
					CreditStoreID: credit.ID,
					Points:        credit.Points,
					Price:         credit.Price,
					Quantity:      1,
				}

				if err := db.Create(&addToCart).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "added new item", addToCart)
				return
			}

			updated.Points += credit.Points
			updated.Price += credit.Price
			updated.Quantity += 1

			var isExist model.CreditStoreCart
			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Added 1 amount", updated)
	})

	// remove 1 amount by id
	r.POST("/remove-amount/:itemID", middleware.Authorization(), func(c *gin.Context) {
		itemID := c.Param("itemID")

		var credit model.CreditStore
		if err := db.Where("id = ?", itemID).First(&credit).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType := c.GetString("type")

		var updated model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("company_id = ?", companyID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points -= credit.Points
			updated.Price -= credit.Price
			updated.Quantity -= 1

			var isExist model.CreditStoreCart
			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if updated.Quantity == 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			var updated model.CreditStoreCart
			if err := db.Where("credit_store_id = ?", itemID).Where("user_id = ?", userID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points -= credit.Points
			updated.Price -= credit.Price
			updated.Quantity -= 1

			var isExist model.CreditStoreCart
			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if updated.Quantity == 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Removed 1 amount", updated)
	})

	// add a custom input if item already in cart
	r.POST("/add-custom-input/:itemID", middleware.Authorization(), func(c *gin.Context) {
		itemID := c.Param("itemID")

		var credit model.CreditStore
		if err := db.Where("id = ?", itemID).First(&credit).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var input model.CreditStoreUpdateQuantityInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		var updated model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			var isExist model.CreditStoreCart
			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("company_id = ?", companyID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points += credit.Points * input.Quantity
			updated.Price += credit.Price * input.Quantity
			updated.Quantity += input.Quantity

			if updated.Quantity == 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			var isExist model.CreditStoreCart
			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("user_id = ?", userID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points += credit.Points * input.Quantity
			updated.Price += credit.Price * input.Quantity
			updated.Quantity += input.Quantity

			if updated.Quantity == 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Added custom amount", updated)
	})

	// remove a custom input if item already in cart
	r.POST("/remove-custom-input/:itemID", middleware.Authorization(), func(c *gin.Context) {
		itemID := c.Param("itemID")

		var credit model.CreditStore
		if err := db.Where("id = ?", itemID).First(&credit).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
			return
		}

		var input model.CreditStoreUpdateQuantityInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		var updated model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			var isExist model.CreditStoreCart
			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("company_id = ?", companyID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points -= credit.Points * input.Quantity
			updated.Price -= credit.Price * input.Quantity
			updated.Quantity -= input.Quantity

			if updated.Quantity == 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("company_id = ?", companyID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			var isExist model.CreditStoreCart
			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).First(&isExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", itemID).Where("user_id = ?", userID).First(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			updated.Points -= credit.Points * input.Quantity
			updated.Price -= credit.Price * input.Quantity
			updated.Quantity -= input.Quantity

			if updated.Quantity == 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			} else if updated.Quantity < 0 {
				if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Delete(&updated).Error; err != nil {
					utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
					return
				}

				utils.HttpRespSuccess(c, http.StatusOK, "Removed item from cart", updated)
				return
			}

			if err := db.Where("user_id = ?", userID).Where("credit_store_id = ?", credit.ID).Save(&updated).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Removed custom amount", updated)
	})

	// remove item from cart
	r.DELETE("/remove-from-cart", middleware.Authorization(), func(c *gin.Context) {
		var input model.CreditStoreCartInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		var cart model.CreditStoreCart
		var isItemExist model.CreditStoreCart

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			if err := db.Where("credit_store_id = ?", input.ID).Where("company_id = ?", companyID).First(&isItemExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", input.ID).Where("company_id = ?", companyID).Delete(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			if err := db.Where("credit_store_id = ?", input.ID).Where("user_id = ?", userID).First(&isItemExist).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("credit_store_id = ?", input.ID).Where("user_id = ?", userID).Delete(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

		}

		utils.HttpRespSuccess(c, http.StatusOK, "Item removed from cart", nil)
	})

	// payment gateway
	r.POST("/payment", middleware.Authorization(), func(c *gin.Context) {
		var total int
		var totalPoints int

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType == "company" {
			companyID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
				return
			}

			var company model.Company
			if err := db.Where("id = ?", companyID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			var cart []model.CreditStoreCart
			if err := db.Where("company_id = ?", companyID).Find(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			for _, v := range cart {
				total += v.Price
				totalPoints += v.Points
			}

			var input model.CreditStorePaymentInput
			if err := c.BindJSON(&input); err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			midtransClient := coreapi.Client{}
			midtransClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
			orderID := utils.RandomOrderID()

			req := &coreapi.ChargeReq{}

			if input.PaymentMethod == 1 {
				req = &coreapi.ChargeReq{
					PaymentType: "gopay",
					TransactionDetails: midtrans.TransactionDetails{
						OrderID:  orderID,
						GrossAmt: int64(total),
					},
					Gopay: &coreapi.GopayDetails{
						EnableCallback: true,
						CallbackUrl:    "https://example.com/callback",
					},
					CustomerDetails: &midtrans.CustomerDetails{
						FName: company.CompanyName,
						Email: company.CompanyEmail,
						Phone: company.CompanyPhone,
					},
				}
			} else if input.PaymentMethod == 2 {
				req = &coreapi.ChargeReq{
					PaymentType: "shopeepay",
					TransactionDetails: midtrans.TransactionDetails{
						OrderID:  orderID,
						GrossAmt: int64(total),
					},
					ShopeePay: &coreapi.ShopeePayDetails{
						CallbackUrl: "https://example.com/callback",
					},
					CustomerDetails: &midtrans.CustomerDetails{
						FName: company.CompanyName,
						Email: company.CompanyEmail,
						Phone: company.CompanyPhone,
					},
				}
			}

			resp, err := midtransClient.ChargeTransaction(req)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			company.Point += totalPoints
			if err := db.Save(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("company_id = ?", companyID).Delete(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			inputTransactionHistory := model.TransactionHistory{
				CompanyID: companyID,
				OrderID:   orderID,
				Price:     total,
				Points:    totalPoints,
				CreatedAt: time.Now(),
			}

			if err := db.Create(&inputTransactionHistory).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			newHistory := model.History{
				CompanyID: companyID,
				Title:     "Buy " + strconv.Itoa(totalPoints) + " Chips",
				Category:  2,
				CreatedAt: time.Now(),
			}

			if err := db.Create(&newHistory).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Payment success", resp)

		} else if strType == "user" {
			userID, ok := ID.(uuid.UUID)
			if !ok {
				utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
				return
			}

			var user model.User
			if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			var cart []model.CreditStoreCart
			if err := db.Where("user_id = ?", userID).Find(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			for _, v := range cart {
				total += v.Price
				totalPoints += v.Points
			}

			var input model.CreditStorePaymentInput
			if err := c.BindJSON(&input); err != nil {
				utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
				return
			}

			midtransClient := coreapi.Client{}
			midtransClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
			orderID := utils.RandomOrderID()

			req := &coreapi.ChargeReq{}

			if input.PaymentMethod == 1 {
				req = &coreapi.ChargeReq{
					PaymentType: "gopay",
					TransactionDetails: midtrans.TransactionDetails{
						OrderID:  orderID,
						GrossAmt: int64(total),
					},
					Gopay: &coreapi.GopayDetails{
						EnableCallback: true,
						CallbackUrl:    "https://example.com/callback",
					},
					CustomerDetails: &midtrans.CustomerDetails{
						FName: user.Name,
						Email: user.Email,
						Phone: user.Phone,
					},
				}
			} else if input.PaymentMethod == 2 {
				req = &coreapi.ChargeReq{
					PaymentType: "shopeepay",
					TransactionDetails: midtrans.TransactionDetails{
						OrderID:  orderID,
						GrossAmt: int64(total),
					},
					ShopeePay: &coreapi.ShopeePayDetails{
						CallbackUrl: "https://example.com/callback",
					},
					CustomerDetails: &midtrans.CustomerDetails{
						FName: user.Name,
						Email: user.Email,
						Phone: user.Phone,
					},
				}
			}

			resp, err := midtransClient.ChargeTransaction(req)
			if err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			var tempUser model.User
			if err := db.Where("id = ?", userID).First(&tempUser).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			tempUser.Point += totalPoints

			if err := db.Model(&user).Updates(&tempUser).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			if err := db.Where("user_id = ?", userID).Delete(&cart).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			inputTransactionHistory := model.TransactionHistory{
				UserID:    userID,
				OrderID:   orderID,
				Price:     total,
				Points:    totalPoints,
				CreatedAt: time.Now(),
			}

			if err := db.Create(&inputTransactionHistory).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			newHistory := model.History{
				UserID:    userID,
				Title:     "Buy " + strconv.Itoa(totalPoints) + " Chips",
				Category:  2,
				CreatedAt: time.Now(),
			}

			if err := db.Create(&newHistory).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Payment success", resp)
		}
	})

}
