package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joseph0x45/tessera/internal/models"
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
		errorMsg := fmt.Sprintf("App+'%s'+already+exists", appName)
		redirectURL := fmt.Sprintf("/admin/dashboard?error=%s", errorMsg)
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
		http.Redirect(w, r, "/admin/dashboard?error=Something+went+wrong.+Check+logs", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	return
}
