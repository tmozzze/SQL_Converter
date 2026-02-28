package postgres

import "github.com/tmozzze/SQL_Converter/internal/domain/models"

func mapDataType(t models.DataType) string {
	switch t {
	case models.DataTypeInteger:
		return "BIGINT"
	case models.DataTypeFloat:
		return "NUMERIC"
	case models.DataTypeBoolean:
		return "BOOLEAN"
	case models.DataTypeString:
		return "TEXT"
	default:
		return "TEXT"
	}
}
