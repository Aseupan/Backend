package Handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"gsc/Entities"
	"gsc/Middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// using google
// func Register(db *gorm.DB, q *gin.Engine) {
// 	r := q.Group("/api")
// 	r.POST("/register", func(c *gin.Context) {
// 		if UserData == nil {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"message": "You are not logged in",
// 				"success": false,
// 			})
// 			return
// 		}

// 		var user GoogleProfile

// 		err := json.Unmarshal(UserData, &user)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			c.Redirect(http.StatusTemporaryRedirect, "/")
// 			return
// 		}

// 		type input struct {
// 			Name string `json:"name"`
// 		}

// 		var inputUser input

// 		if err := c.BindJSON(&inputUser); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"success": false,
// 				"message": "Error when binding JSON",
// 				"error":   err.Error(),
// 			})
// 			return
// 		}

// 		convertID, err := strconv.ParseUint(user.ID, 10, 32)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		parsedID := uint(convertID)
// 		fmt.Println(parsedID)

// 		register := Entities.User{
// 			ID:        parsedID,
// 			Name:      inputUser.Name,
// 			Email:     user.Email,
// 			Points:    0,
// 			CreatedAt: time.Now(),
// 		}

// 		if err := db.Create(&register).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"success": false,
// 				"message": "Error when creating user",
// 				"error":   err.Error(),
// 			})
// 		}

// 		c.JSON(http.StatusOK, gin.H{
// 			"success": true,
// 			"message": "User created",
// 			"ID":      register.ID,
// 			"Name":    register.Name,
// 			"Email":   register.Email,
// 			"Points":  register.Points,
// 			"Time":    register.CreatedAt,
// 		})
// 	})
// }

func Register(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.POST("/register", func(c *gin.Context) {
		type body struct {
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var input Entities.User
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error when binding JSON",
				"error":   err.Error(),
			})
			return
		}

		input.Points = 0
		input.CreatedAt = time.Now()

		if err := db.Create(&input); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with student creation",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Account created successfully",
			"error":   nil,
			"data": gin.H{
				"nama":     input.Name,
				"username": input.Username,
			},
		})
	})
}

func Login(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	// user login
	r.POST("/login", func(c *gin.Context) {
		type body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var input body
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Error when binding JSON",
				"error":   err.Error(),
			})
			return
		}
		login := Entities.User{}
		if err := db.Where("username = ?", input.Username).First(&login).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username does not exist",
				"error":   err.Error(),
			})
			return
		}
		if login.Password == input.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})
			godotenv.Load("../.env")
			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_G")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when loading token",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Welcome, take your token",
				"data": gin.H{
					"username": login.Username,
					"token":    strToken,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Wrong password",
			})
			return
		}
	})
}

func Profile(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api")
	r.GET("/profile", Middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var user Entities.User
		if err := db.Where("id = ?", ID).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong",
				"error":   err.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	})
}

func Hash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	pw := hex.EncodeToString(hash.Sum(nil))
	return pw
}
