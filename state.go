package main

import (
	"github.com/pbojar/dictextract/internal/config"
	"github.com/pbojar/dictextract/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
