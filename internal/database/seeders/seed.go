package main

import (
	"antara-api/common"
	"antara-api/internal/models"
	"github.com/lib/pq"
	"log"
)

func seedDatabase() error {
	db, err := common.DBConnect()
	if err != nil {
		return err
	}

	db.Exec("DELETE FROM user_has_roles")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM products")

	roles := []models.RoleModel{
		{Name: "Administrator"},
		{Name: "Moderator"},
		{Name: "User"},
	}

	password, err := common.HashPassword("password")
	if err != nil {
		panic(err)
	}

	users := []models.UserModel{
		{FirstName: ptr("Demo Admin"), Email: "admin@example.com", PasswordHash: password, Roles: []models.RoleModel{}},
		{FirstName: ptr("Demo Moderator"), Email: "moderator@example.com", PasswordHash: password, Roles: []models.RoleModel{}},
		{FirstName: ptr("Demo User"), Email: "user@example.com", PasswordHash: password, Roles: []models.RoleModel{}},
	}

	products := []models.ProductModel{
		{
			Name:        "Заправка Samsung MLT-D104S (+чип)",
			Slug:        "zaprawka-samsung-mlt-d104s-chip",
			Description: ptr("Заправка картриджа Samsung MLT-D104S с установкой чипа. Заправка производится на профессиональном оборудовании с использованием высококачественного тонера. После заправки картридж проходит тестирование на печать и проверку чипа."),
			Images:      pq.StringArray{"/products/p1-1.jpg", "/products/p1-1.jpg"},
			Price:       59.99,
			Brand:       "Samsung",
			Rating:      4.5,
			NumReviews:  10,
			Stock:       5,
			IsFeatured:  true,
			Banner:      ptr("banner-1.jpg"),
		},
		{
			Name:        "Заправка Brother TN-1085",
			Slug:        "zaprawka-brother-tm-1085",
			Description: ptr("Заправка картриджа Brother TN-1075. Заправка производится на профессиональном оборудовании с использованием высококачественного тонера. После заправки картридж проходит тестирование на печать."),
			Images:      pq.StringArray{"/products/p2-1.jpg", "/products/p2-1.jpg"},
			Price:       85.9,
			Brand:       "Brother",
			Rating:      4.2,
			NumReviews:  8,
			Stock:       10,
			IsFeatured:  true,
			Banner:      ptr("banner-2.jpg"),
		},
		{
			Name:        "Картридж РО MLT-D104S",
			Slug:        "kartirzh-ro-mlt-d104s",
			Description: ptr("Картридж РО MLT-D104S"),
			Images:      pq.StringArray{"/products/p3-1.jpg", "/products/p3-1.jpg"},
			Price:       99.95,
			Brand:       "Samsung",
			Rating:      4.9,
			NumReviews:  3,
			Stock:       0,
			IsFeatured:  false,
			Banner:      ptr(""),
		},
		{
			Name:        "Заправка HP CF226X",
			Slug:        "zaprawka-hp-cf226x",
			Description: ptr("Заправка картриджа HP CF226X"),
			Images:      pq.StringArray{"/products/p4-1.jpg", "/products/p4-1.jpg"},
			Price:       39.95,
			Brand:       "HP",
			Rating:      3.6,
			NumReviews:  5,
			Stock:       10,
			IsFeatured:  false,
			Banner:      ptr(""),
		},
		{
			Name:        "G&G NT-CF226X",
			Slug:        "g-g-nt-cf226x",
			Description: ptr("Картридж G&G NT-CF226X"),
			Images:      pq.StringArray{"/products/p5-1.png", "/products/p5-1.png"},
			Price:       79.99,
			Brand:       "Brother",
			Rating:      4.7,
			NumReviews:  18,
			Stock:       6,
			IsFeatured:  false,
			Banner:      ptr(""),
		},
		{
			Name:        "Картридж C-CLT-Y406S совместимый",
			Slug:        "kartirzh-c-clt-y406s-sovmestimiy",
			Description: ptr("Картридж C-CLT-Y406S совместимый"),
			Images:      pq.StringArray{"/products/p5-1.png", "/products/p5-1.png"},
			Price:       99.99,
			Brand:       "Samsung",
			Rating:      4.6,
			NumReviews:  12,
			Stock:       8,
			IsFeatured:  true,
			Banner:      ptr(""),
		},
	}

	if err := db.Create(&roles).Error; err != nil {
		return err
	}

	for i := range users {

		var role *models.RoleModel
		switch *users[i].FirstName {
		case "Demo Admin":
			role = &roles[0] // Администратор
		case "Demo Moderator":
			role = &roles[1] // Модератор
		case "Demo User":
			role = &roles[2] // Пользователь
		}

		if role != nil {
			users[i].Roles = append(users[i].Roles, *role)
		}

		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}

		for _, role := range users[i].Roles {
			var count int64
			err := db.Model(&models.RoleModel{}).
				Joins("JOIN user_has_roles ON roles.id = user_has_roles.role_id").
				Where("user_has_roles.user_id = ? AND user_has_roles.role_id = ?", users[i].ID, role.ID).
				Count(&count).Error
			if err != nil {
				return err
			}

			if count == 0 {
				if err := db.Exec("INSERT INTO user_has_roles (user_id, role_id) VALUES (?, ?)", users[i].ID, role.ID).Error; err != nil {
					return err
				}
			}
		}
	}

	if err := db.Create(&products).Error; err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func ptr(s string) *string {
	return &s
}

func main() {
	if err := seedDatabase(); err != nil {
		log.Fatal(err)
	}
}
