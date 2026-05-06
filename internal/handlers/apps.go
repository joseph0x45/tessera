package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/tessera/internal/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (h *Handler) processAppCreation(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Error while parsing form:", err.Error())
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	appName := r.FormValue("name")
	if h.conn.AppNameIsTaken(appName) {
		errorMsg := fmt.Sprintf("App+'%s'+already+exists", appName)
		redirectURL := fmt.Sprintf("/dashboard?error=%s", errorMsg)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	newApp := &models.App{
		ID:   gonanoid.Must(),
		Name: appName,
	}
	err := h.conn.InsertApp(newApp)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/dashboard?error=Something+went+wrong.+Check+logs", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) processAppDeletion(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	if err := h.conn.DeleteApp(appID); err != nil {
		http.Redirect(w, r, "/dashboard?error=Something+went+wrong.+Check+logs", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *Handler) renderAppPage(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	urlError := r.URL.Query().Get("error")
	app, err := h.conn.GetAppByID(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users, err := h.conn.GetUsersByAppID(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.cachedUsers = users
	h.render(w, "app", map[string]any{
		"Users": h.cachedUsers,
		"App":   app,
		"Error": urlError,
	})
}
