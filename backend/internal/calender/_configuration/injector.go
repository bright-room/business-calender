package _configuration

import (
	"go.uber.org/dig"
	"gorm.io/gorm"
)

var injector *dig.Container

func init() {
	injector = dig.New()
	opts := createOption()

	if err := injector.Provide(func() *gorm.DB { return opts.DB }); err != nil {
		panic(err)
	}
}
