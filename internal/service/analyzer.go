package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/tmozzze/SQL_Converter/internal/domain"
	"github.com/tmozzze/SQL_Converter/internal/domain/models"
)

type schemaAnalyzerService struct {
	log *slog.Logger
}

func newSchemaAnalyzerService(log *slog.Logger) domain.SchemaAnalyzerService {
	return &schemaAnalyzerService{log: log}
}

// Analyze - analyzing data schema
func (s *schemaAnalyzerService) Analyze(ctx context.Context, tableName string, data [][]string) (models.Table, error) {
	const op = "service.analyzer.Analyze"
	log := s.log.With("op", op)

	if len(data) < 1 {
		return models.Table{}, fmt.Errorf("%s: %w", op, domain.ErrEmptyData)
	}

	headers := data[0]
	if len(headers) == 0 {
		return models.Table{}, fmt.Errorf("%s: %w", op, domain.ErrNoColumns)
	}

	table := models.Table{
		Name:    tableName,
		Columns: make([]models.Column, len(headers)),
	}

	usedNames := make(map[string]int)

	for i, h := range headers {
		name := strings.TrimSpace(h)
		if name == "" {
			name = fmt.Sprintf("col_%d", i+1)
		}

		if count, exists := usedNames[name]; exists {
			usedNames[name]++
			name = fmt.Sprintf("%s_%d", name, count)
		} else {
			usedNames[name] = 1
		}

		table.Columns[i] = models.Column{
			Name: name,
			Type: models.DataTypeUnknown,
		}
	}

	// only headers --> type of columns (string)
	if len(data) == 1 {
		for i := range table.Columns {
			table.Columns[i].Type = models.DataTypeString
		}
		log.Debug("only the headers arrived: all type is (string)")
		return table, nil
	}

	// analyze data
	rows := data[1:]
	for count, row := range rows {

		// checking context
		if count%15 == 0 {
			select {
			case <-ctx.Done():
				return models.Table{}, fmt.Errorf("%s: context canceled: %w", op, ctx.Err())
			default:
			}
		}

		for i, val := range row {
			if i >= len(table.Columns) {
				break
			}
			if table.Columns[i].Type == models.DataTypeString {
				continue
			}
			table.Columns[i].Type = s.detectType(val, table.Columns[i].Type)
		}
	}

	for i := range table.Columns {
		if table.Columns[i].Type == models.DataTypeUnknown {
			table.Columns[i].Type = models.DataTypeString
		}
	}

	return table, nil
}

func (s *schemaAnalyzerService) detectType(val string, currentType models.DataType) models.DataType {
	// string
	if currentType == models.DataTypeString {
		return models.DataTypeString
	}

	// trimming
	val = strings.TrimSpace(val)
	if val == "" {
		return currentType
	}

	// boolean
	lowVal := strings.ToLower(val)
	isBool := lowVal == "true" || lowVal == "false"

	// integer
	_, errInt := strconv.ParseInt(val, 10, 64)
	isInt := errInt == nil

	// float
	_, errFloat := strconv.ParseFloat(val, 64)
	isFloat := errFloat == nil

	switch currentType {
	// unknown
	case models.DataTypeUnknown:
		if isBool {
			return models.DataTypeBoolean
		}
		if isInt {
			return models.DataTypeInteger
		}
		if isFloat {
			return models.DataTypeFloat
		}
		return models.DataTypeString
	// boolean
	case models.DataTypeBoolean:
		if isBool {
			return models.DataTypeBoolean
		}
		return models.DataTypeString

	// integer
	case models.DataTypeInteger:
		if isInt {
			return models.DataTypeInteger
		}
		if isFloat {
			return models.DataTypeFloat
		}
		return models.DataTypeString

	// float
	case models.DataTypeFloat:
		if isFloat {
			return models.DataTypeFloat
		}
		return models.DataTypeString
	default:
		return models.DataTypeString
	}

}
