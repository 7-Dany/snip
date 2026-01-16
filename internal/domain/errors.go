package domain

import "errors"

var (
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrEmptyTitle    = errors.New("title cannot be empty")
	ErrEmptyLanguage = errors.New("language cannot be empty")
	ErrEmptyCode     = errors.New("code cannot be empty")
)
