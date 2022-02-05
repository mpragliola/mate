package mate

import (
	"html/template"
	"time"
)

// Page is the logical representation of a post for the application.
type Page struct {
	CreatedOn time.Time
	Title     string
	Tags      []string
	Source    string
	Body      template.HTML
	Author    string
	Layout    string
	Path      string
	File      string
	Slug      string
	Project   *Project
}
