package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func UserRegister(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/user-register", func(c *gin.Context) {
		var input model.UserRegisterInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if input.Password != input.ConfirmPassword {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Password and confirm password does not match")
			return
		}

		if !utils.IsEmailValid(input.Email) {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Email is not valid")
			return
		}

		if !utils.IsPasswordValid(input.Password) {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Password is not valid")
			return
		}

		var existingUser model.User
		if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "Email is already registered")
			return
		}

		var existingCompany model.Company
		if err := db.Where("company_email = ?", input.Email).First(&existingCompany).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "Email is already registered")
			return
		}

		newUser := model.User{
			ID:        uuid.New(),
			Name:      input.Name,
			Email:     input.Email,
			Password:  utils.Hash(input.Password),
			CreatedAt: time.Now(),
		}

		if err := db.Create(&newUser).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusCreated, "Account created", input)
	})
}

func CompanyRegister(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/company-register", func(c *gin.Context) {
		var input model.CompanyRegisterInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		if input.Password != input.ConfirmPassword {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Password and confirm password does not match")
			return
		}

		if !utils.IsEmailValid(input.CompanyEmail) {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Email is not valid")
			return
		}

		if !utils.IsPasswordValid(input.Password) {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, "Password is not valid")
			return
		}

		var existingCompany model.Company
		if err := db.Where("company_email = ?", input.CompanyEmail).First(&existingCompany).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "Email is already registered")
			return
		}

		var existingUser model.User
		if err := db.Where("email = ?", input.CompanyEmail).First(&existingUser).Error; err == nil {
			utils.HttpRespFailed(c, http.StatusConflict, "Email is already registered")
			return
		}

		newCompany := model.Company{
			ID:             uuid.New(),
			CompanyName:    input.CompanyName,
			CompanyAddress: input.CompanyAddress,
			CompanyEmail:   input.CompanyEmail,
			Password:       utils.Hash(input.Password),
			CreatedAt:      time.Now(),
		}

		if err := db.Create(&newCompany).Error; err != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
			return
		}
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// user login
	r.POST("/login", func(c *gin.Context) {
		var input model.LoginInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		var user model.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			// utils.HttpRespFailed(c, http.StatusNotFound, "User not found")
		}

		var company model.Company
		if err := db.Where("company_email = ?", input.Email).First(&company).Error; err != nil {
			// utils.HttpRespFailed(c, http.StatusNotFound, "Company not found")
		}

		var accountType string

		if user.ID != uuid.Nil && utils.CompareHash(input.Password, user.Password) {
			accountType = "user"
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":   user.ID,
				"type": accountType,
				"exp":  time.Now().Add(time.Hour).Unix(),
			})

			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN")))
			if err != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Parsed token", gin.H{
				"name":  user.Name,
				"token": strToken,
				"type":  accountType,
			})

		} else if company.ID != uuid.Nil && utils.CompareHash(input.Password, company.Password) {
			accountType = "company"
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":   company.ID,
				"type": accountType,
				"exp":  time.Now().Add(time.Hour).Unix(),
			})

			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN")))
			if err != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "Parsed token", gin.H{
				"name":  company.CompanyName,
				"token": strToken,
				"type":  accountType,
			})

		} else {
			utils.HttpRespFailed(c, http.StatusForbidden, "Wrong email or password")
			return
		}
	})

}

func ResetPassword(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")

	r.POST("/reset-password", middleware.Authorization(), func(c *gin.Context) {
		var input model.UserResetPasswordInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType == "company" {
			var company model.Company
			if err := db.Where("id = ?", ID).First(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			company.Password = utils.Hash(input.Password)

			if err := db.Save(&company).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}

		} else if strType == "user" {
			var user model.User
			if err := db.Where("id = ?", ID).First(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusNotFound, err.Error())
				return
			}

			user.Password = utils.Hash(input.Password)

			if err := db.Save(&user).Error; err != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Password reset", nil)
	})
}
