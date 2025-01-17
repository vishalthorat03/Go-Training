package utils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEnsureTableExistsAndTruncate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("TRUNCATE TABLE devices RESTART IDENTITY CASCADE").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = EnsureTableExistsAndTruncate(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
