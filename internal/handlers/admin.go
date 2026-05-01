package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/tessera/internal/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var sessions = map[string]bool{}

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("Error while getting session cookie:", err.Error())
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		sessionID := cookie.Value
		if !sessions[sessionID] {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) renderAdminLogin(w http.ResponseWriter, r *http.Request) {
	h.render(w, "login", nil)
}

func (h *Handler) renderAdminDashboard(w http.ResponseWriter, r *http.Request) {
	apps, err := h.conn.GetAllApps()
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	h.render(w, "dashboard", models.DashboardData{
		Apps: apps,
	})
}

func (h *Handler) processAdminLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	adminUsername := os.Getenv("DASHBOARD_USER")
	adminPasswordHash := os.Getenv("DASHBOARD_PASSWORD_HASH")
	if username != adminUsername {
		h.render(w, "login", map[string]string{
			"Error": "Username not found",
		})
		return
	}
	if !goutils.HashMatchesPassword(adminPasswordHash, password) {
		h.render(w, "login", map[string]string{
			"Error": "Invalid credentials",
		})
		return
	}
	sessionID := gonanoid.Must()
	sessions[sessionID] = true
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.version != "debug",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(100 * 365 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
