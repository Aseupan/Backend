package Handlers

import (
	"encoding/json"
	"fmt"
	"gsc/Entities"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/register", func(c *gin.Context) {
		if UserData == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You are not logged in",
				"success": false,
			})
			return
		}

		var user GoogleProfile

		err := json.Unmarshal(UserData, &user)
		if err != nil {
			fmt.Println(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		type input struct {
			Name string `json:"name"`
		}

		var inputUser input

		if err := c.BindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error when binding JSON",
				"error":   err.Error(),
			})
			return
		}

		convertID, err := strconv.ParseUint(user.ID, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		parsedID := uint(convertID)
		fmt.Println(parsedID)

		register := Entities.User{
			ID:        parsedID,
			Name:      inputUser.Name,
			Email:     user.Email,
			Points:    0,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&register).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when creating user",
				"error":   err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User created",
			"ID":      register.ID,
			"Name":    register.Name,
			"Email":   register.Email,
			"Points":  register.Points,
			"Time":    register.CreatedAt,
		})
	})
}
