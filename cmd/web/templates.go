package main

import (
	"html/template"
	"jonppenny.co.uk/website/pkg/forms"
	"jonppenny.co.uk/website/pkg/models"
	"path/filepath"
	"time"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Post            *models.Post
	Page            *models.Page
	User            *models.User
	MediaItem       *models.Media
	Posts           []*models.Post
	Pages           []*models.Page
	Users           []*models.User
	MediaItems      []*models.Media
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (
	map[string]*template.Template,
	error,
) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(
		filepath.Join(
			dir,
			"*.page.tmpl",
		),
	)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(
			filepath.Join(
				dir,
				"*.layout.tmpl",
			),
		)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(
			filepath.Join(
				dir,
				"*.partial.tmpl",
			),
		)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
