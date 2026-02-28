package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tmozzze/SQL_Converter/internal/domain/models"
)

type tableRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func newTableRepository(db *sql.DB, log *slog.Logger) *tableRepository {
	return &tableRepository{db: db, log: log}
}

// Create - create table in DB
func (r *tableRepository) Create(ctx context.Context, table models.Table) error {
	const op = "postgres.table.Create"
	log := r.log.With("op", op)

	// Generate query
	var sb strings.Builder

	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(quoteIdentifier(table.Name))
	sb.WriteString(" (")

	for i, col := range table.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(quoteIdentifier(col.Name))
		sb.WriteString(" ")
		sb.WriteString(mapDataType(col.Type))
	}

	sb.WriteString(");")

	query := sb.String()

	// Go to DB

	log.Debug("CREATE query is ready", "query", query)

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: failed to create table %s: %w", op, table.Name, err)
	}

	return nil
}

// SaveData - save data in DB
func (r *tableRepository) SaveData(ctx context.Context, table models.Table, data [][]string) error {
	const op = "postgres.table.Create"
	log := r.log.With("op", op)

	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Debug("rollback failed", "err", err)
		}
	}()

	query := buildInsertQuery(table)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for i, row := range data {
		args := make([]any, len(row))
		for j, v := range row {
			args[j] = v
		}

		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return fmt.Errorf("%s: failed to insert row %d: %w", op, i+1, err)
		}
	}

	// Go to DB

	log.Debug("INSERT query is ready", "query", query)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func buildInsertQuery(table models.Table) string {
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(quoteIdentifier(table.Name))
	sb.WriteString(" (")

	// Columns
	for i, col := range table.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(quoteIdentifier(col.Name))
	}

	sb.WriteString(") VALUES (")

	for i := range table.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "$%d", i+1)
	}
	sb.WriteString(");")
	return sb.String()
}
