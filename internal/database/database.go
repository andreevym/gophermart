// database/database.go

package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Connect открывает соединение с базой данных
func Connect(uri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate применяет миграции к базе данных
func Migrate(db *sql.DB) error {
	// выполнение миграций базы данных
	return nil
}
