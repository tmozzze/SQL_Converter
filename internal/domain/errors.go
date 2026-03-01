package domain

import "errors"

var (
	ErrUnsupportedExtension = errors.New("unsupported extension")
	ErrEmptyData            = errors.New("file is empty or has no data rows")
	ErrNoColumns            = errors.New("no columns")
)
