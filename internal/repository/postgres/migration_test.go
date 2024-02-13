package postgres_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

func migrate(db *pgxpool.Pool) error {
	return filepath.Walk("../../../migrations", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}

			err = postgres.ApplyMigration(context.TODO(), db, string(bytes))
			if err != nil {
				return fmt.Errorf("apply migration: %w", err)
			}
		}

		return nil
	})
}
