package databases_test

import (
	"strings"
	"testing"
	"time"

	databases "stpNew/backend/database"
	"stpNew/backend/testutil"
)

// Step 6: Empty-database migration test
func TestMigrations_EmptyDatabaseUp(t *testing.T) {
	_, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)

	// 1. Initial status on empty database must be 0
	status, err := databases.GetMigrationStatus(sqlDB)
	if err != nil {
		t.Fatalf("failed to get initial migration status: %v", err)
	}
	if status != 0 {
		t.Fatalf("expected initial migration version 0, got %d", status)
	}

	// 2. Apply migration up
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	// 3. Status must now be 1
	status, err = databases.GetMigrationStatus(sqlDB)
	if err != nil {
		t.Fatalf("failed to query migration status after up: %v", err)
	}
	if status != 1 {
		t.Fatalf("expected migration version 1, got %d", status)
	}

	// 4. Verify all tables exist
	tables := []string{"customers", "admins", "rooms", "orders", "schema_migrations"}
	for _, tbl := range tables {
		var exists bool
		query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1 AND table_schema = 'public')"
		if err := sqlDB.QueryRow(query, tbl).Scan(&exists); err != nil || !exists {
			t.Errorf("expected table '%s' to exist in public schema, exists=%v, err=%v", tbl, exists, err)
		}
	}

	// 5. Verify column data types for orders.date_start and date_finish
	colTypes := map[string]string{
		"date_start":  "timestamp with time zone",
		"date_finish": "timestamp with time zone",
	}
	for col, expectedType := range colTypes {
		var actualType string
		q := "SELECT data_type FROM information_schema.columns WHERE table_name = 'orders' AND column_name = $1"
		if err := sqlDB.QueryRow(q, col).Scan(&actualType); err != nil || actualType != expectedType {
			t.Errorf("expected column '%s' to have type '%s', got '%s', err=%v", col, expectedType, actualType, err)
		}
	}

	// 6. Verify constraints exist
	constraints := []string{
		"chk_orders_date_order",
		"chk_orders_status",
		"fk_orders_customer",
		"fk_orders_room",
		"exclude_overlapping_room_bookings",
	}
	for _, con := range constraints {
		var exists bool
		q := "SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = $1)"
		if err := sqlDB.QueryRow(q, con).Scan(&exists); err != nil || !exists {
			t.Errorf("expected constraint '%s' to exist, exists=%v, err=%v", con, exists, err)
		}
	}

	// 7. Verify indexes exist
	indexes := []string{
		"idx_orders_user_id",
		"idx_orders_room_id",
		"idx_orders_status",
		"idx_orders_date_start",
		"idx_orders_date_finish",
	}
	for _, idx := range indexes {
		var exists bool
		q := "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = $1)"
		if err := sqlDB.QueryRow(q, idx).Scan(&exists); err != nil || !exists {
			t.Errorf("expected index '%s' to exist, exists=%v, err=%v", idx, exists, err)
		}
	}
}

// Step 7: Migration rollback and reapplication
func TestMigrations_RollbackAndReapplication(t *testing.T) {
	_, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)

	// Apply migration 1
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	// Insert test data
	now := time.Now().UTC()
	start := now.Add(2 * time.Hour).Truncate(time.Second)
	finish := now.Add(4 * time.Hour).Truncate(time.Second)

	_, err := sqlDB.Exec("INSERT INTO customers (id, nama, email) VALUES ('c1', 'Customer One', 'c1@example.com')")
	if err != nil {
		t.Fatalf("failed to insert customer: %v", err)
	}
	_, err = sqlDB.Exec("INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room One', 100000)")
	if err != nil {
		t.Fatalf("failed to insert room: %v", err)
	}
	_, err = sqlDB.Exec("INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price) VALUES ('o1', 'c1', 'r1', $1, $2, 'Menunggu Pembayaran', 100000)", start, finish)
	if err != nil {
		t.Fatalf("failed to insert order: %v", err)
	}

	// Apply rollback
	if err := databases.ApplyMigrationsDown(sqlDB); err != nil {
		t.Fatalf("migration down failed: %v", err)
	}

	// Verify status is 0
	status, err := databases.GetMigrationStatus(sqlDB)
	if err != nil || status != 0 {
		t.Fatalf("expected version 0 after rollback, got %d (err: %v)", status, err)
	}

	// Verify exclusion constraint was dropped
	var conExists bool
	q := "SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'exclude_overlapping_room_bookings')"
	if err := sqlDB.QueryRow(q).Scan(&conExists); err != nil || conExists {
		t.Errorf("expected exclusion constraint to be dropped, exists=%v", conExists)
	}

	// Verify order survived rollback and date columns converted to string preserving UTC representation
	var dateStartStr, dateFinishStr string
	err = sqlDB.QueryRow("SELECT date_start, date_finish FROM orders WHERE id = 'o1'").Scan(&dateStartStr, &dateFinishStr)
	if err != nil {
		t.Fatalf("failed to query order after rollback: %v", err)
	}
	if !strings.Contains(dateStartStr, "Z") {
		t.Errorf("expected UTC ISO string with Z after rollback, got %s", dateStartStr)
	}

	// Reapply migration up
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("reapplying migration up failed: %v", err)
	}

	// Retrieve a fresh connection to verify state without client-side prepared statement cache interference
	_, freshSQLDB := testutil.GetTestDB(t)
	defer freshSQLDB.Close()

	// Verify status is 1 again
	status, err = databases.GetMigrationStatus(freshSQLDB)
	if err != nil || status != 1 {
		t.Fatalf("expected version 1 after reapplication, got %d (err: %v)", status, err)
	}

	// Verify order survived and is TIMESTAMPTZ
	var startTZ, finishTZ time.Time
	err = freshSQLDB.QueryRow("SELECT date_start, date_finish FROM orders WHERE id = 'o1'").Scan(&startTZ, &finishTZ)
	if err != nil {
		t.Fatalf("failed to query order after reapplication: %v", err)
	}
	if !startTZ.Equal(start) {
		t.Errorf("expected start %v, got %v", start, startTZ)
	}
}

// Step 8: Valid legacy-data migration test
func TestMigrations_ValidLegacyDataConversion(t *testing.T) {
	_, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)

	// Create legacy schema with VARCHAR date columns
	legacySchema := `
	CREATE TABLE customers (id VARCHAR(255) PRIMARY KEY, nama VARCHAR(255), email VARCHAR(255) UNIQUE);
	CREATE TABLE rooms (id VARCHAR(255) PRIMARY KEY, name VARCHAR(255), price NUMERIC(15,2));
	CREATE TABLE orders (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255),
		room_id VARCHAR(255),
		date_start VARCHAR(255),
		date_finish VARCHAR(255),
		status VARCHAR(50),
		total_price NUMERIC(15,2)
	);
	INSERT INTO customers (id, nama, email) VALUES ('c100', 'Customer Legacy', 'c100@example.com');
	INSERT INTO rooms (id, name, price) VALUES ('r100', 'Room 100', 50000), ('r200', 'Room 200', 75000);
	INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price) VALUES
	  ('ord-1', 'c100', 'r100', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 50000),
	  ('ord-2', 'c100', 'r100', '2026-11-01T10:00:00Z', '2026-11-01T12:00:00Z', 'Selesai', 50000),
	  ('ord-3', 'c100', 'r200', '2026-11-01T08:00:00Z', '2026-11-01T11:00:00Z', 'Dibatalkan', 75000);
	`
	if _, err := sqlDB.Exec(legacySchema); err != nil {
		t.Fatalf("failed to create legacy schema: %v", err)
	}

	// Apply migration up
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up on legacy data failed: %v", err)
	}

	// Verify all 3 orders survived
	var orderCount int
	if err := sqlDB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount); err != nil || orderCount != 3 {
		t.Fatalf("expected 3 orders after migration, got %d (err: %v)", orderCount, err)
	}

	// Verify converted values for ord-1
	var ordID, userID, roomID, status string
	var startTZ, finishTZ time.Time
	var totalPrice float64
	err := sqlDB.QueryRow("SELECT id, user_id, room_id, date_start, date_finish, status, total_price FROM orders WHERE id = 'ord-1'").
		Scan(&ordID, &userID, &roomID, &startTZ, &finishTZ, &status, &totalPrice)
	if err != nil {
		t.Fatalf("failed to query migrated ord-1: %v", err)
	}

	if ordID != "ord-1" || userID != "c100" || roomID != "r100" || status != "Menunggu Pembayaran" || totalPrice != 50000 {
		t.Errorf("order fields corrupted: id=%s user=%s room=%s status=%s price=%f", ordID, userID, roomID, status, totalPrice)
	}

	expectedStart := time.Date(2026, 11, 1, 8, 0, 0, 0, time.UTC)
	if !startTZ.Equal(expectedStart) {
		t.Errorf("expected start %v, got %v", expectedStart, startTZ)
	}
}

// Step 9: Invalid legacy-data migration test matrix
func TestMigrations_InvalidLegacyDataRollback(t *testing.T) {
	_, sqlDB := testutil.GetTestDB(t)

	cases := []struct {
		name        string
		setupSQL    string
		expectError string
	}{
		{
			name: "empty date string",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-1', 'c1', 'r1', '', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000);
			`,
			expectError: "unconverted date rows",
		},
		{
			name: "invalid date text",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-2', 'c1', 'r1', 'not-a-date', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000);
			`,
			expectError: "invalid input syntax for type timestamp",
		},
		{
			name: "finish earlier than start",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-3', 'c1', 'r1', '2026-11-01T12:00:00Z', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000);
			`,
			expectError: "chk_orders_date_order",
		},
		{
			name: "missing referenced customer",
			setupSQL: `
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-4', 'nonexistent-cust', 'r1', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000);
			`,
			expectError: "fk_orders_customer",
		},
		{
			name: "missing referenced room",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-5', 'c1', 'nonexistent-room', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000);
			`,
			expectError: "fk_orders_room",
		},
		{
			name: "unknown status",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price)
			VALUES ('bad-6', 'c1', 'r1', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z', 'StatusAneh', 10000);
			`,
			expectError: "chk_orders_status",
		},
		{
			name: "existing overlapping blocking reservations",
			setupSQL: `
			INSERT INTO customers (id, nama) VALUES ('c1', 'Cust 1');
			INSERT INTO rooms (id, name, price) VALUES ('r1', 'Room 1', 10000);
			INSERT INTO orders (id, user_id, room_id, date_start, date_finish, status, total_price) VALUES
			  ('conf-1', 'c1', 'r1', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z', 'Menunggu Pembayaran', 10000),
			  ('conf-2', 'c1', 'r1', '2026-11-01T09:00:00Z', '2026-11-01T11:00:00Z', 'Selesai', 10000);
			`,
			expectError: "exclude_overlapping_room_bookings",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutil.ResetTestDBSchema(t, sqlDB)

			// Setup base legacy schema
			baseLegacy := `
			CREATE TABLE customers (id VARCHAR(255) PRIMARY KEY, nama VARCHAR(255), email VARCHAR(255) UNIQUE);
			CREATE TABLE rooms (id VARCHAR(255) PRIMARY KEY, name VARCHAR(255), price NUMERIC(15,2));
			CREATE TABLE orders (
				id VARCHAR(255) PRIMARY KEY,
				user_id VARCHAR(255),
				room_id VARCHAR(255),
				date_start VARCHAR(255),
				date_finish VARCHAR(255),
				status VARCHAR(50),
				total_price NUMERIC(15,2)
			);
			`
			if _, err := sqlDB.Exec(baseLegacy); err != nil {
				t.Fatalf("failed to create base legacy schema: %v", err)
			}
			if _, err := sqlDB.Exec(tc.setupSQL); err != nil {
				t.Fatalf("failed to insert test case records: %v", err)
			}

			// Capture row count before migration attempt
			var countBefore int
			if err := sqlDB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&countBefore); err != nil {
				t.Fatalf("failed to count orders before: %v", err)
			}

			// Attempt migration up - EXPECT FAILURE
			err := databases.ApplyMigrationsUp(sqlDB)
			if err == nil {
				t.Fatalf("[%s] expected migration failure, but migration succeeded!", tc.name)
			}
			if !strings.Contains(err.Error(), tc.expectError) {
				t.Logf("[%s] migration failed as expected: %v", tc.name, err)
			}

			// Verify schema_migrations does not have version 1 recorded
			status, err := databases.GetMigrationStatus(sqlDB)
			if err != nil || status != 0 {
				t.Fatalf("[%s] expected migration status 0 after failed migration rollback, got %d (err: %v)", tc.name, status, err)
			}

			// Verify orders table row count is unchanged
			var countAfter int
			if err := sqlDB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&countAfter); err != nil {
				t.Fatalf("[%s] failed to count orders after: %v", tc.name, err)
			}
			if countBefore != countAfter {
				t.Fatalf("[%s] data loss detected! before=%d, after=%d", tc.name, countBefore, countAfter)
			}

			// Verify orders date columns are STILL VARCHAR (not mutated)
			var colType string
			q := "SELECT data_type FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'date_start'"
			if err := sqlDB.QueryRow(q).Scan(&colType); err != nil || colType != "character varying" {
				t.Fatalf("[%s] expected column date_start to remain character varying, got %s", tc.name, colType)
			}
		})
	}
}

// Step 17: Migration CLI behavior test
func TestMigrationStatusAndEdgeCases(t *testing.T) {
	_, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)

	// Status on fresh database
	status, err := databases.GetMigrationStatus(sqlDB)
	if err != nil || status != 0 {
		t.Fatalf("expected status 0, got %d", status)
	}

	// Down when status is 0 should be no-op (no error)
	if err := databases.ApplyMigrationsDown(sqlDB); err != nil {
		t.Fatalf("expected down on version 0 to succeed as no-op, got %v", err)
	}

	// Up to 1
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("up failed: %v", err)
	}

	// Up when already at 1 should be no-op
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("expected up on version 1 to succeed as no-op, got %v", err)
	}

	// Ensure closed DB returns error
	_, tempSQLDB := testutil.GetTestDB(t)
	tempSQLDB.Close()
	_, err = databases.GetMigrationStatus(tempSQLDB)
	if err == nil {
		t.Error("expected error on closed DB")
	}
}
