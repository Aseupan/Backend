package main

import (
	"fmt"
	"gsc/config"
	"gsc/controller"
	"gsc/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// supabase
	databaseConf, err := config.NewDatabase()
	if err != nil {
		panic(err.Error())
	}
	db, err := config.MakeSupaBaseConnectionDatabase(databaseConf)
	if err != nil {
		panic(err.Error())
	}

	// localhost
	// databaseConf, err := config.NewDBLocal()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// db, err := config.MakeLocalhostConnectionDatabase(databaseConf)
	// if err != nil {
	// 	panic(err.Error())
	// }

	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
			"env":     os.Getenv("ENV"),
		})
	})

	controller.CompanyRegister(db, r)
	controller.UserRegister(db, r)
	controller.ResetPassword(db, r)
	controller.Login(db, r)
	controller.Profile(db, r)
	controller.Address(db, r)
	controller.CreditStore(db, r)
	controller.Rewards(db, r)
	controller.Info(db, r)
	controller.Campaign(db, r)
	controller.History(db, r)
	controller.LoginOTP(db, r)

	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err.Error())
	}
}
