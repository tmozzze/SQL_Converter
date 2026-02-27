package models

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeInteger
	DataTypeFloat
	DataTypeBoolean
	DataTypeString
)

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
