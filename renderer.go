package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type TemplateRenderer struct {
	cache       map[string]*template.Template
	mu          sync.RWMutex // to protect the cache map from concurrent access
	tmplDirPath string
}

// Render (templateRenderer) retrieves a template if available from tr
// If a template is not found, it will return with an error
// If a template is found, it will execute that template using the base template
func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) {

	tmpl, err := tr.getTemplate(name)
	if err != nil {
		fmt.Println("Error getting template:", err)
		http.Error(w, "Error retrieving template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Rendering template:", name)
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Template rendered successfully:", tmpl.Name())
}

// getTemplate will check cache for the template
// if not found, it will call parseTemplate to retrieve the template
// and store it in the cache before returning it
func (tr *TemplateRenderer) getTemplate(name string) (*template.Template, error) {

	// if template in cache, return it
	tr.mu.RLock()
	if tmpl, exists := tr.cache[name]; exists {
		fmt.Println("Found template in cache")
		tr.mu.RUnlock()
		return tmpl, nil
	}
	tr.mu.RUnlock()

	// if template not in cache, retrieve and return
	tr.mu.Lock()
	tmpl, err := tr.parseTemplate(name)
	if err != nil {
		tr.mu.Unlock()
		return nil, err
	}
	tr.cache[name] = tmpl
	tr.mu.Unlock()

	fmt.Println("Template added to cache:", name)
	return tmpl, nil
}

// parseTemplate will get the file path for given name,
// the layouts, the partials, and return a template with these files
func (tr *TemplateRenderer) parseTemplate(name string) (*template.Template, error) {

	// Add original path to template files, then add layout and partial files
	p := filepath.Join(tr.tmplDirPath, name)

	files := []string{p}

	// Layouts - add related files to files slice
	p = filepath.Join(tr.tmplDirPath, "layouts", "*.html")
	pFiles, err := filepath.Glob(p)
	if err == nil {
		files = append(files, pFiles...)
	} else {
		fmt.Println("Error globbing layout templates:", err)
	}

	// Partials - add related files to files slice
	p = filepath.Join(tr.tmplDirPath, "partials", "*.html")
	pFiles, err = filepath.Glob(p)
	if err == nil {
		files = append(files, pFiles...)
	} else {
		fmt.Println("Error globbing partial templates:", err)
	}

	// create template with all required files
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		return nil, err
	}

	fmt.Println("Template is being parsed with files:", files)
	return tmpl, nil
}

func NewTemplateRenderer(dirPath string) *TemplateRenderer {
	return &TemplateRenderer{
		cache:       make(map[string]*template.Template),
		tmplDirPath: dirPath,
	}
}
