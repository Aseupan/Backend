package Config

import (
	"fmt"
	entities "gsc/Entities"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	var db *gorm.DB
	var err error

	db, err = gorm.Open(
		mysql.Open(
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASS"),
				os.Getenv("DB_HOST"),
				os.Getenv("DB_NAME"),
			),
		),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.AutoMigrate(
		// &Model.Admin{},
		// &Model.Class{},
		// &Model.Course{},
		// &Model.Student{},
		&entities.User{},
	); err != nil {
		log.Fatal(err.Error())
	}

	return db
}
