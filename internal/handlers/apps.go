package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/joseph0x45/tessera/internal/models"
	"github.com/joseph0x45/tessera/internal/shared"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (h *Handler) processAppCreation(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Error while parsing form:", err.Error())
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	appName := r.FormValue("name")
	if h.conn.AppNameIsTaken(appName) {
		h.render(w, "/admin/dashboard")
	}
	_, err := h.conn.GetAppByName(appName)
	if !errors.Is(err, shared.ErrAppNotFound) {
		// log.Println("Error while getting app by name:", err.Error())
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	newApp := &models.App{
		ID:   gonanoid.Must(),
		Name: appName,
	}
	err = h.conn.InsertApp(newApp)
	if err != nil {
		log.Println(err.Error())
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	return
}
