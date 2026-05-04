package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (h *Handler) processRegistration(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		AppID        string `json:"app_id"`
		UserName     string `json:"username"`
		UserPassword string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Println("Error while decoding request body:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.AppID == "" || payload.UserName == "" || payload.UserName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": shared.ErrRequiredFieldMissing.Error(),
		})
		return
	}
	app, err := h.conn.GetAppByID(payload.AppID)
	if err != nil {
		if errors.Is(err, shared.ErrAppNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hash, err := goutils.HashPassword(payload.UserPassword)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &models.User{
		ID:       gonanoid.Must(),
		AppID:    app.ID,
		Name:     payload.UserName,
		Password: hash,
	}
	newSession := &models.Session{
		ID:            gonanoid.Must(),
		SessionUserID: newUser.ID,
	}
	err = h.conn.InsertUserAndSession(newUser, newSession)
	if err != nil {
		if errors.Is(err, shared.ErrUserExistsInApp) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": newSession.ID,
	})
}

func (h *Handler) processLogin(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		AppID        string `json:"app_id"`
		UserName     string `json:"username"`
		UserPassword string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Println("Error while decoding request body:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.AppID == "" || payload.UserName == "" || payload.UserName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": shared.ErrRequiredFieldMissing.Error(),
		})
		return
	}
	app, err := h.conn.GetAppByID(payload.AppID)
	if err != nil {
		if errors.Is(err, shared.ErrAppNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := h.conn.GetUserByNameAndAppID(payload.UserName, app.ID)
	if err != nil {
		if errors.Is(err, shared.ErrUserNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !goutils.HashMatchesPassword(user.Password, payload.UserPassword) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": shared.ErrInvalidPassword.Error(),
		})
		return
	}
	newSession := &models.Session{
		ID:            gonanoid.Must(),
		SessionUserID: user.ID,
	}
	if err := h.conn.InsertSession(newSession, nil); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": newSession.ID,
	})
}

func (h *Handler) processUserDeletion(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if err := h.conn.DeleteUser(userID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
