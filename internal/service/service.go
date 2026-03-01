package service

import (
	"log/slog"

	"github.com/tmozzze/SQL_Converter/internal/domain"
)

type service struct {
	fileParser     domain.FileParserService
	schemaAnalyzer domain.SchemaAnalyzerService
	processor      domain.ProcessorService
	log            *slog.Logger
}

// NewService - constructor for main service
func NewService(
	repo domain.Repository,
	log *slog.Logger,
) domain.Service {
	parser := newFileParserService(log)
	analyzer := newSchemaAnalyzerService(log)
	processor := newProcessorService(repo, parser, analyzer, log)
	return &service{
		fileParser:     parser,
		schemaAnalyzer: analyzer,
		processor:      processor,
		log:            log,
	}
}

// Parser - return FileParserService
func (s *service) Parser() domain.FileParserService {
	return s.fileParser
}

// SchemaAnalyzer - return SchemaAnalyzerService
func (s *service) SchemaAnalyzer() domain.SchemaAnalyzerService {
	return s.schemaAnalyzer
}

// return ProcessorService
func (s *service) Processor() domain.ProcessorService {
	return s.processor
}
