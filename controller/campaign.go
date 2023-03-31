package controller

import (
	"gsc/middleware"
	"gsc/model"
	"gsc/utils"
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	r.POST("company/create", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")

		if strType != "company" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var company model.Company
		if res := db.Where("id = ?", ID).First(&company); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, res.Error.Error())
			return
		}

		if !company.Verified {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Company not verified")
			return
		}

		if company.CompanyPhone == "" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Company phone number not set")
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
	r.GET("user/all", middleware.Authorization(), func(c *gin.Context) {
		var campaigns []model.Campaign
		if res := db.Find(&campaigns); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", campaigns)
	})

	r.GET("user/detail/:id", middleware.Authorization(), func(c *gin.Context) {
		strType, _ := c.Get("type")

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		id := c.Param("id")

		var campaign model.Campaign
		if res := db.Where("id = ?", id).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "Campaign", campaign)
	})

	r.POST("user/donate/personal/:campaignID", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		strType, _ := c.Get("type")
		campaignID := utils.StringToUint(c.Param("campaignID"), c)

		if strType != "user" {
			utils.HttpRespFailed(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var input model.UserPersonalDonationInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		newDonation := model.UserPersonalDonation{
			UserID:      ID.(uuid.UUID),
			CampaignID:  campaignID,
			Description: input.Description,
			FoodType:    input.FoodType,
			Quantity:    input.Quantity,
			Weight:      input.Weight,
			ExpiredDate: input.ExpiredDate,
		}

		var campaign model.Campaign
		if res := db.Where("id = ?", campaignID).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		if campaign.Progress < campaign.Target {
			if campaign.Progress+input.Quantity > campaign.Target {
				utils.HttpRespFailed(c, http.StatusBadRequest, "target exceeded")
				return
			}

			if res := db.Create(&newDonation); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			utils.HttpRespSuccess(c, http.StatusOK, "created new donation", newDonation)
		}

		utils.HttpRespFailed(c, http.StatusNotFound, "is finished")

	})

	// get user primary address
	r.GET("user/user-primary-address", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var primaryAddress model.Address

		if res := db.Where("user_id = ?", ID).Where("primary_address = ?", true).First(&primaryAddress); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusNotFound, "user doesnt have primary address / user doesnt have an address")
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "primary address", primaryAddress)
	})

	// get campaign address
	r.GET("user/:campaignID/reciever", middleware.Authorization(), func(c *gin.Context) {
		campaignID := utils.StringToUint(c.Param("campaignID"), c)

		var campaign model.Campaign
		if res := db.Where("id = ?", campaignID).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		var company model.Company
		if res := db.Where("id = ?", campaign.CompanyID).First(&company); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		reciever := model.CampaignCompanyReciever{
			Name:    company.CompanyName,
			Phone:   company.CompanyPhone,
			Address: company.CompanyAddress,
		}

		utils.HttpRespSuccess(c, http.StatusOK, "company reciever", reciever)
	})

	// confirm user personal donation
	r.POST("user/donate/:campaignID/confirm", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		campaignID := utils.StringToUint(c.Param("campaignID"), c)

		var input model.UserPersonalDonationConfirmationInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		var donation model.UserPersonalDonation
		if res := db.Where("user_id = ?", ID).Where("campaign_id = ?", campaignID).First(&donation); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		var user model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		var additionalChip int
		for _, value := range input.AdditionalChips {
			if value == 1 {
				additionalChip += 120
			} else if value == 2 {
				additionalChip += 75
			} else if value == 3 {
				additionalChip += 50
			}
		}

		if additionalChip > user.Point {
			utils.HttpRespFailed(c, http.StatusBadRequest, "not enough chips")
			return
		}

		user.Point -= additionalChip
		user.UpdatedAt = time.Now()
		chipAcquisition := utils.GetFoodPoints(donation.FoodType) * donation.Quantity
		user.Point += chipAcquisition
		if err := db.Save(&user); err.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error.Error())
			return
		}

		donation.PickUp = input.PickUp
		donation.AdditionalChips = input.AdditionalChips
		donation.ChipAcquisition = chipAcquisition
		donation.IsDone = true
		donation.UpdatedAt = time.Now()

		if res := db.Save(&donation); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		var campaign model.Campaign
		if res := db.Where("id = ?", campaignID).First(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		campaign.Progress += donation.Quantity

		if res := db.Save(&campaign); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		go func() {
			time.Sleep(30 * time.Minute) // wait for 30 minutes before executing the task

			var campaign model.Campaign
			if res := db.Where("id = ?", campaignID).First(&campaign); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			// execute the desired task (e.g., move donation data to history table)
			newHistory := model.History{
				UserID:    user.ID,
				Title:     "Donate for " + campaign.Name,
				Category:  1,
				CreatedAt: time.Now(),
			}

			if res := db.Create(&newHistory); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			// delete the new donation data
			if res := db.Delete(&donation); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		}()

		if campaign.Progress == campaign.Target {
			if res := db.Delete(&campaign); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}
		}

		utils.HttpRespSuccess(c, http.StatusOK, "donation confirmed", donation)
	})

	// get all user catering
	r.GET("user/catering", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")

		var catering []model.Catering
		if res := db.Where("user_id = ?", ID).Find(&catering); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "user catering", catering)
	})

	// create new catering
	r.POST("user/catering", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		var input model.NewCateringInput

		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		newCatering := model.Catering{
			UserID:          ID.(uuid.UUID),
			Name:            input.Name,
			Phone:           input.Phone,
			Address:         input.Address,
			AddressDetailed: input.AddressDetailed,
			IsSaved:         input.IsSaved,
			CreatedAt:       time.Now(),
		}

		if res := db.Create(&newCatering); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		utils.HttpRespSuccess(c, http.StatusOK, "new catering created", newCatering)

	})

	// user donate through catering
	r.POST("user/donate/catering/:campaignID/:cateringID", middleware.Authorization(), func(c *gin.Context) {
		ID, _ := c.Get("id")
		campaignID := utils.StringToUint(c.Param("campaignID"), c)
		cateringID := utils.StringToUint(c.Param("cateringID"), c)

		var input model.UserPersonalDonationConfirmationInput
		if err := c.BindJSON(&input); err != nil {
			utils.HttpRespFailed(c, http.StatusBadRequest, err.Error())
			return
		}

		donation := model.UserCateringDonation{
			UserID:          ID.(uuid.UUID),
			CampaignID:      campaignID,
			CateringID:      cateringID,
			PickUp:          input.PickUp,
			AdditionalChips: input.AdditionalChips,
			IsDone:          true,
			CreatedAt:       time.Now(),
		}

		var user model.User
		if res := db.Where("id = ?", ID).First(&user); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		var additionalChip int
		for _, value := range input.AdditionalChips {
			if value == 1 {
				additionalChip += 120
			} else if value == 2 {
				additionalChip += 75
			} else if value == 3 {
				additionalChip += 50
			}
		}

		if additionalChip > user.Point {
			utils.HttpRespFailed(c, http.StatusBadRequest, "not enough chips")
			return
		}

		user.Point -= additionalChip
		user.UpdatedAt = time.Now()
		if err := db.Save(&user); err.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, err.Error.Error())
			return
		}

		if res := db.Create(&donation); res.Error != nil {
			utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
			return
		}

		go func() {
			time.Sleep(30 * time.Minute) // wait for 30 minutes before executing the task

			var campaign model.Campaign
			if res := db.Where("id = ?", campaignID).First(&campaign); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			// execute the desired task (e.g., move donation data to history table)
			newHistory := model.History{
				UserID:    user.ID,
				Title:     "Donate for " + campaign.Name,
				Category:  1,
				CreatedAt: time.Now(),
			}

			if res := db.Create(&newHistory); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			// delete the new donation data
			if res := db.Delete(&donation); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			var catering model.Catering
			if res := db.Where("id = ?", cateringID).Where("user_id = ?", user.ID).First(&catering); res.Error != nil {
				utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
				return
			}

			if !catering.IsSaved {
				if res := db.Delete(&catering); res.Error != nil {
					utils.HttpRespFailed(c, http.StatusInternalServerError, res.Error.Error())
					return
				}
			}
		}()

		utils.HttpRespSuccess(c, http.StatusOK, "donation confirmed", donation)
	})
}
