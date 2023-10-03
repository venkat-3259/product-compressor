package routes

import (
	"zocket/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// RegisterPublicRoutes func for describe group of public routes.
func RegisterPublicRoutes(a *fiber.App, h *controllers.Handler) {

	// Create routes group.
	route := a.Group("/api/v1")

	route.Post("/products", h.CreateProduct)
}
