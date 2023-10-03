package queries

import (
	"context"
	"time"

	"zocket/app/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductQueries struct {
	DB *pgxpool.Pool
}

func (q *ProductQueries) CreateProduct(ctx context.Context, arg models.CreateProduct) (models.ProductLinks, error) {

	query := `INSERT INTO products (name, description, images, price, created_at, updated_at) VALUES
						($1, $2, $3, $4, $5, $5) RETURNING id, product_images;`

	var productDetail models.ProductLinks

	err := q.DB.QueryRow(ctx, query, arg.Name, arg.Description, arg.Images, arg.Price, time.Now()).Scan(
		&productDetail.ID, &productDetail.Links)

	return productDetail, err
}

func (q *ProductQueries) AddCompressedImages(ctx context.Context, arg models.ProductImages) error {

	query := `UPDATE products SET compressed_images = $1 WHERE id = $2;`

	_, err := q.DB.Exec(ctx, query, arg.Paths, arg.ID)

	return err
}
