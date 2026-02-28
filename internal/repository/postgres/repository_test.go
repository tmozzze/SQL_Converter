package postgres_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tmozzze/SQL_Converter/internal/domain/models"
	"github.com/tmozzze/SQL_Converter/internal/repository/postgres"
)

func TestCreateTable(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := postgres.NewRepository(db, log)
	ctx := context.Background()

	table := models.Table{
		Name: "users",
		Columns: []models.Column{
			{Name: "id", Type: models.DataTypeInteger},
			{Name: "name", Type: models.DataTypeString},
		},
	}

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS "users" \("id" BIGINT, "name" TEXT\);`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Table().Create(ctx, table)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := postgres.NewRepository(db, log)
	ctx := context.Background()

	table := models.Table{
		Name: "users",
		Columns: []models.Column{
			{Name: "id", Type: models.DataTypeInteger},
			{Name: "name", Type: models.DataTypeString},
		},
	}

	data := [][]string{{"1", "John"}}

	mock.ExpectBegin()
	mock.ExpectPrepare(`INSERT INTO "users" \("id", "name"\) VALUES \(\$1, \$2\);`)
	mock.ExpectExec(`INSERT INTO "users" \("id", "name"\) VALUES \(\$1, \$2\);`).
		WithArgs("1", "John").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Table().SaveData(ctx, table, data)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
