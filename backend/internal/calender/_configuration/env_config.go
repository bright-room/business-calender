package _configuration

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"golang.org/x/xerrors"
)

type envConfig struct {
	BusinessDataSource businessDataSource `envPrefix:"BUSINESS_DB_"`
	TimeZone           string             `env:"TZ,notEmpty"`
}

type businessDataSource struct {
	Host     string `env:"HOST,notEmpty"`
	Port     string `env:"PORT,notEmpty"`
	Username string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Database string `env:"NAME,notEmpty"`
	Schema   string `env:"SCHEMA"  envDefault:""`
	UseSSL   string `env:"USE_SSL" envDefault:"disable"`
}

func (e envConfig) dsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s search_path=%s sslmode=%s TimeZone=%s",
		e.BusinessDataSource.Host,
		e.BusinessDataSource.Username,
		e.BusinessDataSource.Password,
		e.BusinessDataSource.Database,
		e.BusinessDataSource.Port,
		e.BusinessDataSource.Schema,
		e.BusinessDataSource.UseSSL,
		e.TimeZone,
	)
}

func envParse() envConfig {
	var e envConfig
	if err := env.Parse(&e); err != nil {
		panic(xerrors.Errorf("failed to parse env: %w", err))
	}

	return e
}
