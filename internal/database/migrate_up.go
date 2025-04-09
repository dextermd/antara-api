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

	err = db.AutoMigrate(
		&models.RoleModel{},
		&models.CartModel{},
		&models.SocialProviderModel{},
		&models.ShippingAddressModel{},
		&models.UserModel{},
		&models.VerificationTokenModel{},
		&models.CategoryModel{},
		&models.ProductModel{},
		&models.OptionModel{},
		&models.MCategoryModel{},
		&models.SessionModel{},
		&models.CartItemModel{},
		&models.PageModel{},
	)

	if err != nil {
		panic(err.Error())
	}

	log.Println("Migration completed")
}

// go run ./internal/database/migrate_up.go
