package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var templates map[string]*template.Template

// Load templates on program initialisation
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates_path := os.Getenv("TEMPLATES_PATH")
	if templates_path == "" {
		templates_path = "templates"
	}

	layouts, err := filepath.Glob(path.Join(templates_path, "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(path.Join(templates_path, "includes", "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our templates/ directory
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

// renderTemplate() is a wrapper around template.ExecuteTemplate
func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", data)
}

/*
 * Handlers
 */

// Handler for 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := renderTemplate(w, "404.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler for a pattern
func simpleHandler(pattern string, template string,
	w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != pattern {
		notFoundHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	err := renderTemplate(w, template, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	static_path := os.Getenv("STATIC_PATH")
	if static_path == "" {
		static_path = "static"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(static_path))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialize endpoints handlers
	endpoints := []struct {
		pattern string
		template string
	} {
		{"/", "index.html"},
		{"/presentations/", "presentations.html"},
	}

	for _, endpoint := range endpoints {
		/*
		 * We need to declare new variables here since the new clusure will
		 * reuse the last endpoint entry.
		 */
		pattern := endpoint.pattern
		template := endpoint.template
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			simpleHandler(pattern, template, w, r)
		})
	}

	log.Println("Listening...")
	http.ListenAndServe(":" + port, mux)
}
