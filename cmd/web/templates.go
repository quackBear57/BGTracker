package main

import(
	"html/template"
	"path/filepath"

	"bgtracker.ndr54.com/internal/models"
)

type templateData struct {
	Game			models.Game
	Games			[]models.Game
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Add more specific *.tmpl.html here?
	pages, err := filepath.Glob("./assets/html/pages/*tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./assets/html/base.tmpl.html",
			"./assets/html/partials/nav.tmpl.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}