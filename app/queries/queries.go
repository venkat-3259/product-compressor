package queries

import "github.com/jackc/pgx/v4/pgxpool"

// Queries struct for collect all app queries.
type Queries struct {
	*ProductQueries
}

func NewQueries(db *pgxpool.Pool) *Queries {
	return &Queries{
		ProductQueries: &ProductQueries{DB: db},
	}
}
