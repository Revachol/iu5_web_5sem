package main

import (
	"github.com/Revachol/iu5_web_5sem/internal/app/ds"
	"github.com/Revachol/iu5_web_5sem/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Historical_service{},
		&ds.Historical_request{},
		&ds.User{},
		&ds.Historical_request_service{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
