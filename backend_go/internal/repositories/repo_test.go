package repositories

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"log"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
)

func GetDBForTests() (testcontainers.Container, *gorm.DB, error) {
	c := context.Background()
	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(c,
		"docker.io/postgres",
		postgres.WithInitScripts(filepath.Join("../testdata", "init.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, nil, err
	}

	port, err := postgresContainer.MappedPort(c, "5432")
	if err != nil {
		return nil, nil, err
	}

	connURL := fmt.Sprintf("host=localhost port=%s user=user password=password dbname=testdb sslmode=disable", port.Port())
	db, err := gorm.Open(gormpostgres.Open(connURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, db, nil
}

// Test helper function to clean up the database between tests
//
//goland:noinspection ALL
func cleanUpDB(container testcontainers.Container, db *gorm.DB) {
	// Remove all users from the database
	db.Exec("DELETE FROM \"backend_diff\".\"tasks\"")
	db.Exec("DELETE FROM \"backend_diff\".\"users\"")
	if err := container.Terminate(context.Background()); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

// Test suite for userRepository
func setupTest(t *testing.T) (testcontainers.Container, *gorm.DB, func()) {
	postgresContainer, db, err := GetDBForTests()
	if err != nil {
		log.Println("Failed to set up test DB")
		t.FailNow()
	}
	cleanup := func() {
		cleanUpDB(postgresContainer, db) // Clean up after test run
	}
	return postgresContainer, db, cleanup
}
