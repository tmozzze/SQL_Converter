package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tmozzze/SQL_Converter/internal/domain"
)

// Handler - struct for handler
type Handler struct {
	service domain.Service
	log     *slog.Logger
}

// NewHandler - constructor for handler
func NewHandler(service domain.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// Response - struct for response
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// UploadFile godoc
// @Summary Upload a file and create a table
// @Description Accepts .csv or .xlsx, analyzes structure, creates a table in PG and inserts data.
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV or XLSX file"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 422 {object} Response
// @Failure 500 {object} Response
// @Router /upload [post]
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.http.UploadFile"
	log := h.log.With(slog.String("op", op))

	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, errors.New("only POST method is allowed"))
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		log.Error("failed to parse form", slog.Any("err", err))
		h.sendError(w, http.StatusBadRequest, errors.New("file is too large or invalid form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.sendError(w, http.StatusBadRequest, errors.New("field 'file' is required"))
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	err = h.service.Processor().UploadFile(r.Context(), header.Filename, file, ext)
	if err != nil {
		log.Error("failed to process file", slog.Any("err", err))

		h.handleServiceError(w, err)
		return
	}

	h.sendJSON(w, http.StatusOK, Response{Status: "OK"})
}

func (h *Handler) sendError(w http.ResponseWriter, code int, err error) {
	h.sendJSON(w, code, Response{
		Status: "Error",
		Error:  err.Error(),
	})
}

func (h *Handler) sendJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUnsupportedExtension):
		h.sendError(w, http.StatusUnprocessableEntity, domain.ErrUnsupportedExtension)

	case errors.Is(err, domain.ErrEmptyData), errors.Is(err, domain.ErrNoColumns):
		h.sendError(w, http.StatusBadRequest, domain.ErrNoColumns)

	case errors.Is(err, http.ErrAbortHandler):
		return

	default:
		h.sendJSON(w, http.StatusInternalServerError, Response{
			Status: "Error",
			Error:  "internal server error",
		})
	}
}
