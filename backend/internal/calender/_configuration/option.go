package _configuration

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type option struct {
	DB *gorm.DB
}

func createOption() *option {
	e := envParse()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  e.dsn(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &option{
		DB: db,
	}
}
