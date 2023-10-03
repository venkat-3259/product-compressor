package controllers

import (
	"log"

	"zocket/app/models"
	"zocket/pkg/response"
	"zocket/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// @Summary Add Prodcut
// @Id add-product
// @Tags products
// @Description add a new product
// @Accept json
// @Produce      json
// @Param body body models.CreateProduct true "body parameter"
// @Success 200 string ok "returns ok"
// @Failure 404 {object} response.Response "bad request: validate your input params"
// @Failure 500 {object} response.Response
// @Security     ApiKeyAuth
// @Router /products [post]
func (h *Handler) CreateProduct(c *fiber.Ctx) error {

	arg := new(models.CreateProduct)

	// Parsing data from body
	if err := c.BodyParser(arg); err != nil {
		return response.CodeMessage(c, fiber.StatusInternalServerError, err.Error())
	}

	// Json validator, which validates body params
	validate := utils.NewValidator()
	if err := validate.Struct(arg); err != nil {
		return response.CodeMessage(c, fiber.StatusBadRequest, err.Error())
	}

	// Adding product details to the database
	productLinks, err := h.Queries.ProductQueries.CreateProduct(c.Context(), *arg)
	if err != nil {
		log.Println("Failed to insert data in the database, Reason: ", err)
		return response.CodeMessage(c, fiber.StatusInternalServerError, err.Error())
	}

	// Sends link to the downloader
	go h.Process.ImageDownload(productLinks)

	return response.Data(c, "Product Added successfully")
}
