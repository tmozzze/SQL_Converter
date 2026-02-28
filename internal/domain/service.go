package domain

import (
	"io"

	"github.com/tmozzze/SQL_Converter/internal/domain/models"
)

// Service - interface for buisness logic
type Service interface {
	Parser() FileParserService
	SchemaAnalyzer() SchemaAnalyzerService
}

// FileParserService - interface for file parser buisness logic
type FileParserService interface {
	Parse(r io.Reader, extension string) ([][]string, error)
}

// SchemaAnalyzerService - interface for schema analyzer buisness logic
type SchemaAnalyzerService interface {
	Analyze(tableName string, data [][]string) (models.Table, error)
}
