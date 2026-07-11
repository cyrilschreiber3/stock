package migrations

import (
	"embed"
)

//go:embed *.sql
var SchemaMigrations embed.FS
