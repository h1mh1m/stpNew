package databases

import (
	"strings"
	"testing"
)

func TestLoadMigrationSQL_EmbeddedAccess(t *testing.T) {
	// 1. Up migration for version 1
	upSQL, err := LoadMigrationSQL(1, "up")
	if err != nil {
		t.Fatalf("unexpected error loading migration 1 up: %v", err)
	}
	if len(strings.TrimSpace(upSQL)) == 0 {
		t.Fatal("expected non-empty up migration SQL")
	}
	if !strings.Contains(upSQL, "btree_gist") {
		t.Error("expected up migration to contain btree_gist extension")
	}
	if !strings.Contains(upSQL, "exclude_overlapping_room_bookings") {
		t.Error("expected up migration to contain exclude_overlapping_room_bookings constraint")
	}

	// 2. Down migration for version 1
	downSQL, err := LoadMigrationSQL(1, "down")
	if err != nil {
		t.Fatalf("unexpected error loading migration 1 down: %v", err)
	}
	if len(strings.TrimSpace(downSQL)) == 0 {
		t.Fatal("expected non-empty down migration SQL")
	}
	if !strings.Contains(downSQL, "DROP CONSTRAINT IF EXISTS exclude_overlapping_room_bookings") {
		t.Error("expected down migration to drop exclusion constraint")
	}

	// 3. Invalid version
	_, err = LoadMigrationSQL(999, "up")
	if err == nil {
		t.Error("expected error for non-existent migration version 999")
	}

	// 4. Invalid direction
	_, err = LoadMigrationSQL(1, "sideways")
	if err == nil {
		t.Error("expected error for invalid migration direction")
	}
}
