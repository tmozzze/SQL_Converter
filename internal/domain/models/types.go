package models

// DataType - represent a data type
type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeInteger
	DataTypeFloat
	DataTypeBoolean
	DataTypeString
)

// String - return a DataType string
func (d DataType) String() string {
	switch d {
	case DataTypeInteger:
		return "Integer"
	case DataTypeFloat:
		return "Float"
	case DataTypeBoolean:
		return "Boolean"
	case DataTypeString:
		return "String"
	default:
		return "Unknown"
	}

}
