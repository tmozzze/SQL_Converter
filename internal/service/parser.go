package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"

	"github.com/tmozzze/SQL_Converter/internal/domain"
	"github.com/xuri/excelize/v2"
)

type fileParserService struct {
	log *slog.Logger
}

func newFileParserService(log *slog.Logger) domain.FileParserService {
	return &fileParserService{log: log}
}

// Parse - parsing file to Table from io.Reader with extension(.csv, .xlsx)
func (s *fileParserService) Parse(ctx context.Context, r io.Reader, extension string) ([][]string, error) {
	const op = "service.parser.Parse"
	log := s.log.With("op", op)

	// context checking
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: context canceled: %w", op, ctx.Err())
	default:
	}

	// choose extension
	switch extension {
	case domain.ExtCSV:
		log.Debug("parsing .CSV")
		return s.parseCSV(ctx, r)
	case domain.ExtXLSX:
		log.Debug("parsing .XLSX")
		return s.parseXLSX(ctx, r)
	default:
		return nil, fmt.Errorf("%s: failed to read file: %s: %w", op, extension, domain.ErrUnsupportedExtension)
	}
}

func (s *fileParserService) parseCSV(ctx context.Context, r io.Reader) ([][]string, error) {
	const op = "service.parser.parseCSV"
	log := s.log.With("op", op)

	reader := csv.NewReader(r)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	log.Debug("rows parsed", "count", len(rows))

	return rows, nil
}

func (s *fileParserService) parseXLSX(ctx context.Context, r io.Reader) ([][]string, error) {
	const op = "service.parser.parseXLSX"
	log := s.log.With("op", op)

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open XLSX: %w", op, err)
	}
	defer f.Close()

	if f.SheetCount == 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrEmptyData)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read XLSX: %w", op, err)
	}

	log.Debug("rows parsed", "count", len(rows))

	return rows, nil
}
