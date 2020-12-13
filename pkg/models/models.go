package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Media struct {
	ID          int
	Title       string
	Filename    string
	Path        string
	Description string
	Width       int64
	Height      int64
	Size        int64
	Created     time.Time
	Updated     time.Time
}

type Page struct {
	ID      int
	Title   string
	Slug    string
	Status  string
	Content string
	Created time.Time
	Updated time.Time
}

type Post struct {
	ID      int
	Title   string
	Content string
	Status  string
	Created time.Time
	Updated time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	LastLogin      time.Time
	Active         bool
	Role           string
	Created        time.Time
	Updated        time.Time
}
