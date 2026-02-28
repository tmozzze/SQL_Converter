package postgres

import (
	"database/sql"
	"log/slog"

	"github.com/tmozzze/SQL_Converter/internal/domain"
)

type Repository struct {
	table domain.TableRepository
	log   *slog.Logger
}

func NewRepository(db *sql.DB, log *slog.Logger) *Repository {
	return &Repository{
		table: newTableRepository(db, log),
		log:   log,
	}
}

func (r *Repository) Table() domain.TableRepository {
	return r.table
}
