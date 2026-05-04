package handlers

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/login", h.renderAdminLogin)
	r.Post("/login", h.processAdminLogin)

	r.With(h.requireAdmin).Get("/dashboard", h.renderAdminDashboard)
	r.With(h.requireAdmin).Get("/dashboard/logout", h.processLogout)
	r.With(h.requireAdmin).Post("/apps", h.processAppCreation)
	r.With(h.requireAdmin).Post("/apps/{id}/delete", h.processAppDeletion)
	r.With(h.requireAdmin).Get("/apps/{id}", h.renderAppPage)

	r.Post("/api/users", h.processUserCreation)
}
