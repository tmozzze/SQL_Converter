package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tmozzze/SQL_Converter/internal/domain"
)

type processorService struct {
	repo     domain.Repository
	parser   domain.FileParserService
	analyzer domain.SchemaAnalyzerService
	log      *slog.Logger
}

func newProcessorService(
	repo domain.Repository,
	parser domain.FileParserService,
	analyzer domain.SchemaAnalyzerService,
	log *slog.Logger,
) domain.ProcessorService {
	return &processorService{
		repo:     repo,
		parser:   parser,
		analyzer: analyzer,
		log:      log,
	}
}

func (s *processorService) UploadFile(ctx context.Context, tableName string, file io.Reader, extension string) error {
	const op = "service.processor.UploadFile"
	log := s.log.With("op", op)

	// table name
	cleanTableName := sanitizeTableName(tableName)

	// parsing
	rawData, err := s.parser.Parse(ctx, file, extension)
	if err != nil {
		return fmt.Errorf("%s: parsing failed: %w", op, err)
	}

	// analyzing
	table, err := s.analyzer.Analyze(ctx, cleanTableName, rawData)
	if err != nil {
		return fmt.Errorf("%s: analysis failed: %w", op, err)
	}

	// go to DB (create table)
	if err := s.repo.Table().Create(ctx, table); err != nil {
		return fmt.Errorf("%s: create failed: %w", op, err)
	}

	// go to DB (insert data)
	if len(rawData) > 1 {
		if err := s.repo.Table().SaveData(ctx, table, rawData[1:]); err != nil {
			return fmt.Errorf("%s: repo save failed: %w", op, err)
		}
	}

	log.Debug("file processed successfully", "table", cleanTableName, "rows", len(rawData)-1)

	return nil
}

func sanitizeTableName(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	name = reg.ReplaceAllString(name, "_")
	name = strings.ToLower(strings.Trim(name, "_"))
	if name == "" {
		return "imported_table"
	}
	return name
}
