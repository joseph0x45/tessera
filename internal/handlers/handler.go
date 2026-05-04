package handlers

import (
	"html/template"
	"net/http"

	"github.com/joseph0x45/tessera/internal/db"
	"github.com/joseph0x45/tessera/internal/models"
)

type Handler struct {
	conn        *db.Conn
	templates   *template.Template
	version     string
	cachedApps  []models.App
	cachedUsers []models.User
}

func NewHandler(
	conn *db.Conn,
	templates *template.Template,
	version string,
) *Handler {
	return &Handler{
		conn:       conn,
		templates:  templates,
		version:    version,
		cachedApps: nil,
	}
}

func (h *Handler) render(w http.ResponseWriter, templateName string, data any) {
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
