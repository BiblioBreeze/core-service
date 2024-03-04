package core_service

import "embed"

//go:embed migrations/*.sql
var embedMigrations embed.FS

func GetMigrationsFS() embed.FS {
	return embedMigrations
}
