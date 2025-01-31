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
