// package utils

// import (
// 	"testing"

// 	_ "github.com/lib/pq"
// 	"github.com/stretchr/testify/assert"
// )

// func TestConnectDatabase(t *testing.T) {
// 	// Call ConnectDatabase to establish a connection
// 	db, err := ConnectDatabase()

// 	// Assert no error occurred and database connection is established
// 	assert.NoError(t, err)
// 	assert.NotNil(t, db)

// 	// Check the connection pool settings using db.Stats()
// 	stats := db.Stats()

// 	// Assert that the number of open connections is within expected range (as an example)
// 	assert.GreaterOrEqual(t, stats.OpenConnections, 0)
// }

// func TestEnsureTableExistsAndTruncate(t *testing.T) {
// 	// Establish a mock or real DB connection
// 	db, err := ConnectDatabase()
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.Close()

// 	// Call EnsureTableExistsAndTruncate to ensure table exists and is truncated
// 	err = EnsureTableExistsAndTruncate(db)

// 	// Assert no error occurred
// 	assert.NoError(t, err)

// 	// Verify the table 'devices' exists in the database
// 	var tableExists bool
// 	err = db.QueryRow(`SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'devices')`).Scan(&tableExists)
// 	if err != nil {
// 		t.Fatalf("Failed to check if table exists: %v", err)
// 	}

// 	// Assert the table exists
// 	assert.True(t, tableExists)

// 	// Verify the table was truncated by checking row count
// 	var rowCount int
// 	err = db.QueryRow(`SELECT COUNT(*) FROM devices`).Scan(&rowCount)
// 	if err != nil {
// 		t.Fatalf("Failed to count rows in the 'devices' table: %v", err)
// 	}

// 	// Assert no rows are present in the table after truncation
// 	assert.Equal(t, 0, rowCount)
// }

package utils

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEnsureTableExistsAndTruncate(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Expect the TRUNCATE query
	mock.ExpectExec("TRUNCATE TABLE devices RESTART IDENTITY CASCADE").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = EnsureTableExistsAndTruncate(db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestEnsureTableExistsAndCreate(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Simulate the TRUNCATE query failing because the table does not exist
	mock.ExpectExec("TRUNCATE TABLE devices RESTART IDENTITY CASCADE").
		WillReturnError(fmt.Errorf(`pq: relation "devices" does not exist`))

	// Expect the CREATE TABLE query as a fallback
	mock.ExpectExec("CREATE TABLE devices .*").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the function being tested
	err = EnsureTableExistsAndTruncate(db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

// package utils

// import (
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// )

// func TestEnsureTableExistsAndTruncate(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("failed to open mock db: %v", err)
// 	}
// 	defer db.Close()

// 	mock.ExpectExec("TRUNCATE TABLE devices RESTART IDENTITY CASCADE").
// 		WillReturnResult(sqlmock.NewResult(0, 1))

// 	err = EnsureTableExistsAndTruncate(db)
// 	if err != nil {
// 		t.Errorf("unexpected error: %v", err)
// 	}

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("unmet expectations: %v", err)
// 	}
// }
