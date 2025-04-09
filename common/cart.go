package common

import (
	"antara-api/internal/models"
	"github.com/google/uuid"
	"math"
)

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func CalcPrice(items []models.CartItemModel) (itemsPrice, shippingPrice, taxPrice, totalPrice float64) {
	itemsPrice = round2(
		func() float64 {
			var total float64
			for _, item := range items {
				total += item.Price * float64(item.Qty)
			}
			return total
		}(),
	)

	shippingPrice = round2(0.0)

	taxPrice = round2(0.0 * itemsPrice)

	totalPrice = round2(itemsPrice + shippingPrice + taxPrice)

	return
}

func FindCartItem(items []models.CartItemModel, productID uint) *models.CartItemModel {
	for i := range items {
		if items[i].ProductID == productID {
			return &items[i]
		}
	}
	return nil
}

func GenerateCartID() string {
	return uuid.New().String()
}
