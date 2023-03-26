package controller

import (
	"gsc/model"
	"gsc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Info(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/user/info")
	r.GET("/all-info", func(c *gin.Context) {
		var infos []model.Info
		if res := db.Find(&infos); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Queried all info", infos)
	})

	r.GET("/detailed/:id", func(c *gin.Context) {
		var info model.Info
		id := c.Param("id")
		if res := db.Where("id = ?", id).First(&info); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Queried info", info)
	})
}
