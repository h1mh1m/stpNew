package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ExpectedTestDBName = "trackiq_integration_test"
	DefaultTestPort    = "55432"
	DefaultTestHost    = "127.0.0.1"
	DefaultTestUser    = "trackiq_test"
	DefaultTestPass    = "testpassword123"
)

// GetTestDB connects to the isolated PostgreSQL integration test database.
// It strictly verifies that the target database is trackiq_integration_test.
func GetTestDB(t *testing.T) (*gorm.DB, *sql.DB) {
	t.Helper()

	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = DefaultTestHost
	}
	port := os.Getenv("TEST_DB_PORT")
	if port == "" {
		port = DefaultTestPort
	}
	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = DefaultTestUser
	}
	pass := os.Getenv("TEST_DB_PASSWORD")
	if pass == "" {
		pass = DefaultTestPass
	}
	dbname := os.Getenv("TEST_DB_NAME")
	if dbname == "" {
		dbname = ExpectedTestDBName
	}

	if dbname != ExpectedTestDBName {
		t.Fatalf("SAFETY VIOLATION: Integration tests may only run against '%s', got '%s'", ExpectedTestDBName, dbname)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, dbname, port)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("PostgreSQL integration test database unavailable at %s:%s (%v). Skipping integration test.", host, port, err)
		return nil, nil
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("Failed to retrieve underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Skipf("PostgreSQL ping failed at %s:%s (%v). Skipping integration test.", host, port, err)
		return nil, nil
	}

	// Verify connected database name directly from the database server
	var currentDB string
	if err := sqlDB.QueryRow("SELECT current_database()").Scan(&currentDB); err != nil {
		t.Fatalf("Failed to query current database name: %v", err)
	}
	if currentDB != ExpectedTestDBName {
		t.Fatalf("SAFETY VIOLATION: Connected to unexpected database '%s', expected '%s'", currentDB, ExpectedTestDBName)
	}

	return gormDB, sqlDB
}

// ResetTestDBSchema drops all tables and extensions in trackiq_integration_test.
func ResetTestDBSchema(t *testing.T, sqlDB *sql.DB) {
	t.Helper()

	var currentDB string
	if err := sqlDB.QueryRow("SELECT current_database()").Scan(&currentDB); err != nil {
		t.Fatalf("Failed to query current database name: %v", err)
	}
	if currentDB != ExpectedTestDBName {
		t.Fatalf("SAFETY VIOLATION: Refusing to reset non-test database: %s", currentDB)
	}

	queries := []string{
		"DROP TABLE IF EXISTS orders CASCADE;",
		"DROP TABLE IF EXISTS rooms CASCADE;",
		"DROP TABLE IF EXISTS customers CASCADE;",
		"DROP TABLE IF EXISTS admins CASCADE;",
		"DROP TABLE IF EXISTS schema_migrations CASCADE;",
		"DROP EXTENSION IF EXISTS btree_gist CASCADE;",
	}

	for _, q := range queries {
		if _, err := sqlDB.Exec(q); err != nil {
			t.Fatalf("Failed to execute reset query '%s': %v", q, err)
		}
	}
}
