package service_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmozzze/SQL_Converter/internal/domain"
	"github.com/tmozzze/SQL_Converter/internal/repository/postgres"

	"github.com/tmozzze/SQL_Converter/internal/service"
)

func TestProcessorService_UploadFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repo := postgres.NewRepository(db, log)
	svc := service.NewService(repo, log)
	processor := svc.Processor()

	ctx := context.Background()

	t.Run("successful CSV upload with mixed types", func(t *testing.T) {
		csvData := `name,age,salary,is_active
Sasha,25,50000.50,true
Masha,30,60000.75,false
Petr,35,55000.00,true`

		reader := strings.NewReader(csvData)

		// Waiting CREATE TABLE with true types
		createQuery := `DROP TABLE IF EXISTS "users" CASCADE;CREATE TABLE IF NOT EXISTS "users" \("name" TEXT, "age" BIGINT, "salary" NUMERIC, "is_active" BOOLEAN\);`
		mock.ExpectExec(createQuery).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectBegin()
		mock.ExpectPrepare(`INSERT INTO "users" \("name", "age", "salary", "is_active"\) VALUES \(\$1, \$2, \$3, \$4\);`)

		mock.ExpectExec(`INSERT INTO "users" .*`).
			WithArgs("Sasha", "25", "50000.50", "true").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`INSERT INTO "users" .*`).
			WithArgs("Masha", "30", "60000.75", "false").
			WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectExec(`INSERT INTO "users" .*`).
			WithArgs("Petr", "35", "55000.00", "true").
			WillReturnResult(sqlmock.NewResult(3, 1))

		mock.ExpectCommit()

		err := processor.UploadFile(ctx, "users", reader, domain.ExtCSV)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
