package Config

import (
	"fmt"
	"gsc/Entities"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf(
			"user=%s password=%s host=%s TimeZone=%s port=%s dbname=%s",
			os.Getenv("SB_User"),
			os.Getenv("SB_Password"),
			os.Getenv("SB_Host"),
			os.Getenv("SB_TimeZone"),
			os.Getenv("SB_Port"),
			os.Getenv("SB_DB"),
		),
	), &gorm.Config{})

	if err != nil {
		return nil
	}

	if err := db.AutoMigrate(&Entities.User{}); err != nil {
		return nil
	}

	return db
}
