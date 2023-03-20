package controller

import (
	"gsc/middleware"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserHome(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user")
	// user home
	r.GET("/home", middleware.Authorization(), func(c *gin.Context) {
		utils.HttpRespSuccess(c, http.StatusOK, "only validated user can see this!!", nil)
	})
}
