package main

import(
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
	"strings"

	"bgtracker.ndr54.com/internal/models"
	"bgtracker.ndr54.com/assets"
)

type templateData struct {
	CurrentYear			int
	Game				models.Game
	Games				[]models.Game
	Form				any
	Flash				string
	IsAuthenticated		bool
	CSRFToken			string
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate" : humanDate,
	"toTitle"	: strings.Title,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(assets.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(assets.Files,
			patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}