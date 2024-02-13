package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ApplyMigration(ctx context.Context, db *pgxpool.Pool, sql string) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := db.Exec(rCtx, sql)
	if err != nil {
		return fmt.Errorf("failed apply sql '%s': %w", sql, err)
	}

	return nil
}
