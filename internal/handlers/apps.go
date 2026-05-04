package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
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
		"Error": "",
	})
}

func (h *Handler) processUserCreation(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		AppID            string `json:"app_id"`
		UserName         string `json:"username"`
		UserPasswordHash string `json:"password_hash"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Println("Error while decoding request body:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.AppID == "" || payload.UserName == "" || payload.UserName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app, err := h.conn.GetAppByID(payload.AppID)
	if err != nil {
		if errors.Is(err, shared.ErrAppNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &models.User{
		ID:       gonanoid.Must(),
		AppID:    app.ID,
		Name:     payload.UserName,
		Password: payload.UserPasswordHash,
	}
	err = h.conn.InsertUser(newUser)
	if err != nil {
		if errors.Is(err, shared.ErrUserExistsInApp) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
