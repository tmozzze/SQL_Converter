package postgres

import (
	"database/sql"
	"log/slog"

	"github.com/tmozzze/SQL_Converter/internal/domain"
)

// Repository - main repository struct
type repository struct {
	table domain.TableRepository
	log   *slog.Logger
}

// NewRepository - constructor for Repository
func NewRepository(db *sql.DB, log *slog.Logger) *repository {
	return &repository{
		table: newTableRepository(db, log),
		log:   log,
	}
}

// Table - return TableRepository
func (r *repository) Table() domain.TableRepository {
	return r.table
}
