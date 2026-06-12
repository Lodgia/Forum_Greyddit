package handlers

import (
	"greyddit/database"
	"/static/template"
	"golang.org/x/crypto/bcrypt"
)

var funcMap = template.FuncMap{
	"sub": func(a, b int) int { return a - b },
}