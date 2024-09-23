package seeds

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (s *Seed) usersSeed(total int) {
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

	type generalData struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var (
		roles    = make([]generalData, 0)
		userMaps = make([]map[string]any, 0)
	)

	err = s.db.Select(&roles, `SELECT id, name FROM roles`)
	if err != nil {
		log.Error().Err(err).Msg("Error selecting roles")
		return
	}

	for i := 0; i < total; i++ {
		selectedRole := roles[gofakeit.Number(0, len(roles)-1)]

		dataUserToInsert := make(map[string]any)
		dataUserToInsert["id"] = ulid.Make().String()
		dataUserToInsert["role_id"] = selectedRole.Id
		dataUserToInsert["name"] = gofakeit.Name()
		dataUserToInsert["email"] = gofakeit.Email()
		dataUserToInsert["whatsapp_number"] = gofakeit.Phone()
		dataUserToInsert["password"] = "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G" // password

		userMaps = append(userMaps, dataUserToInsert)
	}

	var (
		endUserId   string
		adminUserId string
	)

	// iterate over roles to get service advisor id
	for _, role := range roles {
		if role.Name == "admin" {
			adminUserId = role.Id
			continue
		}
		if role.Name == "end_user" {
			endUserId = role.Id
			continue
		}
	}

	EndUser := map[string]any{
		"id":              ulid.Make().String(),
		"role_id":         endUserId,
		"name":            "Irham",
		"email":           "irham@fake.com",
		"whatsapp_number": gofakeit.Phone(),
		"password":        "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G", // password
	}

	AdminUser := map[string]any{
		"id":              ulid.Make().String(),
		"role_id":         adminUserId,
		"name":            "Fathan",
		"email":           "fathan@fake.com",
		"whatsapp_number": gofakeit.Phone(),
		"password":        "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G", // password
	}

	userMaps = append(userMaps, EndUser)
	userMaps = append(userMaps, AdminUser)

	_, err = tx.NamedExec(`
		INSERT INTO users (id, role_id, name, email, whatsapp_number, password)
		VALUES (:id, :role_id, :name, :email, :whatsapp_number, :password)
	`, userMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users")
		return
	}

	log.Info().Msg("users table seeded successfully")
}
