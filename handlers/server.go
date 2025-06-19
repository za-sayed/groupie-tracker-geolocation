package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Global variable to store templates
var tpl *template.Template

// InitializeTemplates loads all templates from the provided directory path.
func InitializeTemplates(basePath string) error {
	templatesPath := filepath.Join(basePath, "templates", "*")
	var err error
	tpl, err = template.ParseGlob(templatesPath)
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}
	return nil
}


// StartStaticFileServer sets up a static file server for serving files under the "/static/" route.
func StartStaticFileServer(basePath string) {
	staticFilesPath := filepath.Join(basePath, "static")
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(staticFilesPath)))
	http.Handle("/static/", staticHandler)
}

// StartServer initializes the templates, configures routes, and starts the HTTP server.
func StartServer() {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting the current working directory: %v", err)
	}

	// Initialize the templates
	if err := InitializeTemplates(basePath); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Serve static files
	StartStaticFileServer(basePath)

	// Set up routes
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/Artists", GalleryHandler)
	http.HandleFunc("/Artists/", ArtistHandler)

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
