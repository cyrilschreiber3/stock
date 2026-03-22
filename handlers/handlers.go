package handlers

import (
	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/database/repository"
)

var db *repository.Queries

func Init() {
	db = repository.New(database.Pool)
}
