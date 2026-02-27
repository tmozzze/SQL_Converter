package domain

import (
	"context"

	"github.com/tmozzze/SQL_Converter/internal/domain/models"
)

// Repository - interface for data repositories
type Repository interface {
	Table() TableRepository
}

// TableRepository - interface for table operations
type TableRepository interface {
	// Create - create table
	Create(ctx context.Context, table models.Table) error
	// SaveData - save data in table
	SaveData(ctx context.Context, table models.Table, data [][]string) error
}
