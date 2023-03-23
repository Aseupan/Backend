package middleware

import (
	"gsc/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]

		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN")), nil
		})
		if err != nil {
			log.Println("pas bagian ngambil token error")
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			idStr := claims["id"].(string)
			id, err := uuid.Parse(idStr)
			if err != nil {
				log.Println("pas bagian ngambil id error")
				utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
				c.Abort()
			}
			c.Set("id", id)
			c.Next()
			return
		} else {
			utils.HttpRespFailed(c, http.StatusForbidden, err.Error())
			c.Abort()
			return
		}
	}
}
