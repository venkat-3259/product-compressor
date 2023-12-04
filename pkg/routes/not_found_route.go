package routes

import "github.com/gofiber/fiber/v2"

// RegisterNotFoundRoute func for describe 404 Error route.
func RegisterNotFoundRoute(a *fiber.App) {
	// Register new special route.
	a.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    404,
				"message": "sorry, endpoint is not found",
			})
		},
	)
}
