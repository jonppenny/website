package main

import (
	"jonppenny.co.uk/webapp/pkg/forms"
	"jonppenny.co.uk/webapp/pkg/models"
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
