package models

type CreateProduct struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Images      []string `json:"images" validate:"required,min=1"`
	Price       float64  `json:"price" validate:"required"`
} //@name CreateProduct

type ProductLinks struct {
	ID    int
	Links []string
} //@ProductLinks

type ProductImages struct {
	ID    int
	Paths []string
}
