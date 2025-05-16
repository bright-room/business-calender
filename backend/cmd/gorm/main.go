package main

import (
	"gorm.io/gen"
	"net.bright-room.dev/calender-api/internal/calender/_configuration"
)

func main() {
	cfg := _configuration.NewGormGenConfiguration()

	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/calender/infrastructure/datasource/db/query",
		ModelPkgPath:  "internal/calender/infrastructure/datasource/db/entity",
		Mode:          gen.WithDefaultQuery,
		FieldNullable: false,
	})

	g.UseDB(cfg.DB)

	tables := g.GenerateAllTable()
	g.ApplyBasic(tables...)

	g.Execute()
}
