package postgres_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
)

var testDB *pgxpool.Pool

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
	dbURL    string
)

func setupDB() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}

	resource, err = pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=testuser",
		"POSTGRES_PASSWORD=testpass",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		log.Fatalf("Could not start PostgreSQL container: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pool.Retry(func() error {
		dbURL = fmt.Sprintf("postgres://testuser:testpass@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp"))
		testDB, err = pgxpool.Connect(ctx, dbURL)
		if err != nil {
			return err
		}

		err = migrate(testDB)
		if err != nil {
			log.Fatalf("migration PostgreSQL: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL container: %v", err)
	}
}

func teardownDB() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %v", err)
	}
	testDB.Close()
}

func TestMain(m *testing.M) {
	setupDB()
	code := m.Run()
	teardownDB()
	os.Exit(code)
}
