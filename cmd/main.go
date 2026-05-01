package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/tessera/internal/buildinfo"
	"github.com/joseph0x45/tessera/internal/cli"
	"github.com/joseph0x45/tessera/internal/db"
	"github.com/joseph0x45/tessera/internal/handlers"
)

//go:embed tailwind.css
var tailwindCSS string

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

func init() {
	var err error
	funcMap := template.FuncMap{}
	templates, err = template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

func main() {
	goutils.SetAppName(buildinfo.AppName)
	goutils.Setup()

	cli.DispatchCommands(os.Args)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	conn := db.GetConn(false)
	defer conn.Close()
	handler := handlers.NewHandler(conn, templates, buildinfo.Version)

	r := chi.NewRouter()
	server := http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
	}

	handler.RegisterRoutes(r)

	goutils.ServeJS(r, "/tailwindcss", tailwindCSS)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf(
			"Starting %s %s on http://0.0.0.0:%s",
			buildinfo.AppName, buildinfo.Version, port,
		)
		log.Printf("Admin dashboard http://0.0.0.0:%s/admin/dashboard", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server startup failed (addr=%s): %v", server.Addr, err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown failed: %s\n", err)
	}
}
