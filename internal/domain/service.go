package domain

import (
	"context"
	"io"

	"github.com/tmozzze/SQL_Converter/internal/domain/models"
)

const (
	ExtXLSX = ".xlsx"
	ExtCSV  = ".csv"
)

// Service - interface for buisness logic
type Service interface {
	Parser() FileParserService
	SchemaAnalyzer() SchemaAnalyzerService
	Processor() ProcessorService
}

// FileParserService - interface for file parser buisness logic
type FileParserService interface {
	Parse(ctx context.Context, r io.Reader, extension string) ([][]string, error)
}

// SchemaAnalyzerService - interface for schema analyzer buisness logic
type SchemaAnalyzerService interface {
	Analyze(ctx context.Context, tableName string, data [][]string) (models.Table, error)
}

// ProcessorService - interface for process manager
type ProcessorService interface {
	UploadFile(ctx context.Context, tableName string, file io.Reader, extension string) error
}
