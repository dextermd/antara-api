package handlers

import (
	"antara-api/cmd/api/dtos/requests"
	"antara-api/cmd/api/services"
	"antara-api/common"
	"antara-api/internal/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) AddItemToCart(c echo.Context) error {
	var sessionCartID string
	var sessionID string

	cartIDCookie, err := c.Cookie("session_cart_id")
	if err != nil || cartIDCookie == nil {
		return common.SendBadRequestResponse(c, "Cart session not found")
	}
	sessionCartID = cartIDCookie.Value

	payload := new(requests.AddCartItemRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	sessionIDCookie, err := c.Cookie("session_id")
	if err == nil || sessionIDCookie != nil {
		sessionID = sessionIDCookie.Value
	}

	sessionService := services.NewSessionService(h.DB)
	user, _ := sessionService.GetUserFromSession(sessionID)

	var userID *uint
	if user != nil {
		userID = &user.ID
	}

	productService := services.NewProductService(h.DB)
	product, err := productService.GetByID(h.DB, payload.ProductID)
	if err != nil {
		return common.SendBadRequestResponse(c, "Product not found")
	}

	cartItem := models.CartItemModel{
		ProductID: payload.ProductID,
		Qty:       payload.Qty,
		Price:     product.Price,
		Slug:      product.Slug,
		Name:      product.Name,
		Image:     product.Images[0],
	}

	cartService := services.NewCartService(h.DB)
	cart, err := cartService.GetCart(sessionCartID, userID)
	if err != nil {
		return common.SendBadRequestResponse(c, "Error fetching cart")
	}

	if cart == nil {
		cart = &models.CartModel{
			SessionCartID: sessionCartID,
			UserID:        userID,
			Items:         []models.CartItemModel{cartItem},
		}

		itemsPrice, shippingPrice, taxPrice, totalPrice := common.CalcPrice(cart.Items)
		cart.ItemsPrice = itemsPrice
		cart.ShippingPrice = shippingPrice
		cart.TaxPrice = taxPrice
		cart.TotalPrice = totalPrice

		if err := cartService.Create(cart); err != nil {
			return common.SendBadRequestResponse(c, "Error creating cart")
		}

	} else {
		existingItem := common.FindCartItem(cart.Items, payload.ProductID)
		if existingItem != nil {
			if product.Stock < (existingItem.Qty + cartItem.Qty) {
				return common.SendBadRequestResponse(c, "Product out of stock")
			}
			existingItem.Qty += cartItem.Qty
		} else {
			if product.Stock < cartItem.Qty {
				return common.SendBadRequestResponse(c, "Product out of stock")
			}
			cart.Items = append(cart.Items, cartItem)
		}

		itemsPrice, shippingPrice, taxPrice, totalPrice := common.CalcPrice(cart.Items)
		cart.ItemsPrice = itemsPrice
		cart.ShippingPrice = shippingPrice
		cart.TaxPrice = taxPrice
		cart.TotalPrice = totalPrice

		if err := cartService.Update(cart); err != nil {
			return common.SendBadRequestResponse(c, "Error updating cart")
		}
	}
	return common.SendSuccessResponse(c, "Item added to cart successfully", cart)
}

func (h *Handler) RemoveItemFromCart(c echo.Context) error {
	payload := new(requests.AddCartItemRequest)
	if err := c.Bind(payload); err != nil {
		return common.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedValidationResponse(c, validationErrors)
	}

	var sessionCartID string
	var sessionID string

	cartIDCookie, err := c.Cookie("session_cart_id")
	if err != nil || cartIDCookie == nil {
		return common.SendBadRequestResponse(c, "Cart session not found in cookie")
	}
	sessionCartID = cartIDCookie.Value

	sessionIDCookie, err := c.Cookie("session_id")
	if err == nil || sessionIDCookie != nil {
		sessionID = sessionIDCookie.Value
	}

	sessionService := services.NewSessionService(h.DB)
	user, _ := sessionService.GetUserFromSession(sessionID)

	var userID *uint
	if user != nil {
		userID = &user.ID
	}

	cartService := services.NewCartService(h.DB)
	cart, err := cartService.GetCart(sessionCartID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.SendNotFoundResponse(c, "Cart not found")
		}
		return common.SendBadRequestResponse(c, "Error fetching cart")
	}

	existingItem := common.FindCartItem(cart.Items, payload.ProductID)
	if existingItem == nil {
		return common.SendBadRequestResponse(c, "Item not found in cart")
	}

	if existingItem.Qty > 1 {
		existingItem.Qty -= 1
	} else {
		for i, item := range cart.Items {
			if item.ProductID == payload.ProductID {
				cartItemService := services.NewCartItemService(h.DB)
				err := cartItemService.Delete(&item)
				if err != nil {
					return common.SendBadRequestResponse(c, "Error removing item from cart")
				}
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				break
			}
		}
	}

	itemsPrice, shippingPrice, taxPrice, totalPrice := common.CalcPrice(cart.Items)
	cart.ItemsPrice = itemsPrice
	cart.ShippingPrice = shippingPrice
	cart.TaxPrice = taxPrice
	cart.TotalPrice = totalPrice

	if err := cartService.Update(cart); err != nil {
		return common.SendBadRequestResponse(c, "Error updating cart")
	}

	return common.SendSuccessResponse(c, "Item removed from cart successfully", cart)
}

func (h *Handler) GetCartHandler(c echo.Context) error {
	var sessionCartID string
	var sessionID string

	cartIDCookie, err := c.Cookie("session_cart_id")
	if err != nil || cartIDCookie == nil {
		return common.SendBadRequestResponse(c, "Cart session not found in cookie")
	}
	sessionCartID = cartIDCookie.Value

	sessionIDCookie, err := c.Cookie("session_id")
	if err == nil || sessionIDCookie != nil {
		sessionID = sessionIDCookie.Value
	}

	sessionService := services.NewSessionService(h.DB)
	user, _ := sessionService.GetUserFromSession(sessionID)

	var userID *uint
	if user != nil {
		userID = &user.ID
	}

	cartService := services.NewCartService(h.DB)
	cart, err := cartService.GetCart(sessionCartID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return common.SendBadRequestResponse(c, "Error fetching cart")
	}

	return common.SendSuccessResponse(c, "Cart retrieved successfully", cart)
}
