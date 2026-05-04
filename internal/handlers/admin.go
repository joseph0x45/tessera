package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("Error while getting session cookie:", err.Error())
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		sessionID := cookie.Value
		adminSessionID, err := h.conn.GetMetadata("admin_session")
		if err != nil {
			if !errors.Is(err, shared.ErrValueNotFound) {
				log.Println(err)
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if sessionID != *adminSessionID {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.cachedApps = apps
	urlError := r.URL.Query().Get("error")
	h.render(w, "dashboard", models.DashboardData{
		Apps:  apps,
		Error: urlError,
	})
}

func (h *Handler) processAdminLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	password := r.FormValue("password")
	adminPasswordHash, err := h.conn.GetMetadata("admin_password")
	if err != nil {
		if !errors.Is(err, shared.ErrValueNotFound) {
			log.Println(err.Error())
		}
		h.render(w, "login", map[string]string{
			"Error": "Something went wrong",
		})
		return
	}
	if !goutils.HashMatchesPassword(*adminPasswordHash, password) {
		h.render(w, "login", map[string]string{
			"Error": "Invalid password",
		})
		return
	}
	sessionID := gonanoid.Must()
	err = h.conn.SetMetadata(&models.MetaData{
		Key:   "admin_session",
		Value: sessionID,
	})
	if err != nil {
		log.Println(err)
		h.render(w, "login", map[string]string{
			"Error": "Something went wrong",
		})
		return
	}
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
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
