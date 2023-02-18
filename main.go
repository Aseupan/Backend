package main

import (
	"fmt"
	"gsc/Config"
	"gsc/Handlers"
	"gsc/Middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	// Connect Database
	db := Config.Connect()
	if db != nil {
		println("Nice, DB Connected")
	}

	r := gin.Default()
	Handlers.Init()
	r.Use(Middleware.CORS())
	r.GET("/", Handlers.HandleMain)
	r.GET("/login", Handlers.HandleGoogleLogin)
	r.GET("/callback", Handlers.HandleGoogleCallback)
	r.GET("/status", Handlers.HandleStatus)
	r.GET("/TLI", Handlers.HandleTestLoggedIn)
	r.GET("/logout", Handlers.HandleLogout)

	r.Group("/api")
	Handlers.Register(db, r)

	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
