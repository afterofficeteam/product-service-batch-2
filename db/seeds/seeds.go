package seeds

import (
	"codebase-app/internal/adapter"
	"context"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {

	switch table {
	case "roles":
		s.rolesSeed()
	case "users":
		s.usersSeed(total)
	case "shops":
		s.shopsSeed(total)
	case "product_categories":
		s.productCategoriesSeed(total)
	case "products":
		s.productsSeed(total)
	case "all":
		s.rolesSeed()
		s.usersSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("Error committing transaction")
			}
		}
	}()

	_, err = tx.Exec(`DELETE FROM users`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users")
		return
	}
	log.Info().Msg("users table deleted successfully")

	_, err = tx.Exec(`DELETE FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting roles")
		return
	}
	log.Info().Msg("roles table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

// rolesSeed seeds the roles table.
func (s *Seed) rolesSeed() {
	roleMaps := []map[string]any{
		{"name": "admin"},
		{"name": "end_user"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(`
		INSERT INTO roles (name)
		VALUES (:name)
	`, roleMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating roles")
		return
	}

	log.Info().Msg("roles table seeded successfully")
}

func (s *Seed) productCategoriesSeed(total int) {
	var (
		args  = make([]map[string]any, 0)
		query = "INSERT INTO product_categories (name) VALUES (:name)"
	)

	for i := 0; i < total; i++ {
		var (
			name = gofakeit.ProductCategory()
			arg  = make(map[string]any)
		)

		arg["name"] = name
		args = append(args, arg)
	}

	_, err := s.db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Error creating product categories")
	}

	log.Info().Msg("product_categories table seeded successfully")
}

func (s *Seed) shopsSeed(total int) {
	var (
		query = "INSERT INTO shops (name) VALUES (:name)"
	)

	query = "INSERT INTO shops (name, description, terms, user_id) VALUES (?, ?, ?, ?)"

	for i := 0; i < total; i++ {
		_, err := s.db.Exec(s.db.Rebind(query), gofakeit.Company(), gofakeit.HackerPhrase(), gofakeit.HackerPhrase(), gofakeit.UUID())
		if err != nil {
			log.Error().Err(err).Msg("Error creating shops")
			return
		}
	}

	log.Info().Msg("shops table seeded successfully")
}

func (s *Seed) productsSeed(total int) {
	var (
		query = `
			INSERT INTO products(
				name,
				description,
				price,
				stock,
				shop_id,
				category_id
			) VALUES (
			 	?, ?, ?,
				20000,
				(SELECT id FROM shops ORDER BY RANDOM() LIMIT 1),
				(SELECT id FROM product_categories ORDER BY RANDOM() LIMIT 1)
			)`
	)

	for i := 0; i < total; i++ {
		_, err := s.db.Exec(s.db.Rebind(query),
			gofakeit.Product().Name,
			gofakeit.HackerPhrase(),
			gofakeit.Price(1000, 10000),
		)
		if err != nil {
			log.Error().Err(err).Msg("Error creating products")
			return
		}
	}

	log.Info().Msg("products table seeded successfully")
}
