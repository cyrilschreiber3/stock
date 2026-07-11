package seed

import "embed"

//go:embed *.sql
var SeedMigrations embed.FS
