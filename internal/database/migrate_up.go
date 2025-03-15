package main

import (
	"antara-api/common"
	"antara-api/internal/models"
	"log"
)

func main() {
	db, err := common.DBConnect()
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&models.RoleModel{}, &models.CartModel{}, &models.ShippingAddressModel{}, &models.UserModel{}, &models.AppTokenModel{})

	if err != nil {
		panic(err.Error())
	}

	log.Println("Migration completed")
}

// go run ./internal/database/migrate_up.go
