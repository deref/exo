package background

import "github.com/jmoiron/sqlx"

type TaskStore struct {
	DB *sqlx.DB
}
