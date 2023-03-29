package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"log"
	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func Campaign(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/api/campaign")

	SupaBaseClient := supabasestorageuploader.NewSupabaseClient(
		"https://flldkbhntqqaiflpxlhg.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZsbGRrYmhudHFxYWlmbHB4bGhnIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY3NzU4Njk4OCwiZXhwIjoxOTkzMTYyOTg4fQ.CezKv4eOdEOyPEnVCqp3i0rNRLpz4MJOgL2GvM74QtQ",
		"photo",
		"",
	)

	// big party / company
	r.POST("/create", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType != "company" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		name := c.PostForm("name")

		description := c.PostForm("description")

		target := utils.StringToInteger(c.PostForm("target"), c)

		area := c.PostForm("area")

		startDate := c.PostForm("startdate")

		endDate := c.PostForm("enddate")

		urgent := utils.StringToInteger(c.PostForm("urgent"), c)

		foodType := c.PostFormArray("type[]")
		log.Print(foodType)
		log.Println(pq.Array(foodType))

		thumbnail1, _ := c.FormFile("thumbnail1")

		thumbnail2, _ := c.FormFile("thumbnail2")

		thumbnail3, _ := c.FormFile("thumbnail3")

		thumbnail4, _ := c.FormFile("thumbnail4")

		thumbnail5, _ := c.FormFile("thumbnail5")

		link1, _ := SupaBaseClient.Upload(thumbnail1)

		link2, _ := SupaBaseClient.Upload(thumbnail2)

		link3, _ := SupaBaseClient.Upload(thumbnail3)

		link4, _ := SupaBaseClient.Upload(thumbnail4)

		link5, _ := SupaBaseClient.Upload(thumbnail5)

		var newCampaign model.Campaign
		newCampaign.CompanyID = ID.(uuid.UUID)
		newCampaign.Name = name
		newCampaign.Description = description
		newCampaign.Target = target
		newCampaign.Area = area
		newCampaign.StartDate = startDate
		newCampaign.EndDate = endDate
		newCampaign.Urgent = urgent
		newCampaign.Type = foodType
		newCampaign.Thumbnail1 = link1
		newCampaign.Thumbnail2 = link2
		newCampaign.Thumbnail3 = link3
		newCampaign.Thumbnail4 = link4
		newCampaign.Thumbnail5 = link5

		if res := db.Create(&newCampaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "created new campaign", newCampaign)

	})

	// user
	r.GET("/all", middleware.Authorization(), func(c *gin.Context) {
		var campaigns []model.Campaign
		if res := db.Find(&campaigns); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", nil)
	})

	r.GET("/detail/:id", middleware.Authorization(), func(c *gin.Context) {
		// strType, _ := c.Get("type")

		// if strType != "user" {
		// 	utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
		// 	return
		// }

		id := c.Param("id")

		var campaign model.Campaign
		if res := db.Where("id = ?", id).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", campaign)
	})
}
