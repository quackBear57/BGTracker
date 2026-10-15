package main

import(
	"html/template"
	"path/filepath"
	"time"

	"bgtracker.ndr54.com/internal/models"
)

type templateData struct {
	CurrentYear		int
	Game			models.Game
	Games			[]models.Game
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate" : humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./assets/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./assets/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./assets/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}