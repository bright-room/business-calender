package _configuration

import (
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type GormGenConfiguration struct {
	DB *gorm.DB
}

func NewGormGenConfiguration() *GormGenConfiguration {
	i := injector

	var db *gorm.DB
	if err := i.Invoke(func(instance *gorm.DB) {
		db = instance
	}); err != nil {
		panic(xerrors.Errorf("failed to resolving dependencies a db session: %w", err))
	}

	return &GormGenConfiguration{DB: db}
}
