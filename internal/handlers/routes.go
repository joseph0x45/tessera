package handlers

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/admin/login", h.renderAdminLogin)

	r.With(h.requireAdmin).Get("/admin/dashboard", h.renderAdminDashboard)
	r.Post("/admin/login", h.processAdminLogin)

}
