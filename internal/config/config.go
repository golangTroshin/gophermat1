package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var (
	Options struct {
		ServerAddress        string
		DataBaseURI          string
		AccrualSystemAddress string
	}

	Config struct {
		ServerAddress        string `env:"RUN_ADDRESS"`
		DataBaseURI          string `env:"DATABASE_URI"`
		AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	}
)

func ParseFlags() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	if Config.ServerAddress != "" {
		Options.ServerAddress = Config.ServerAddress
	} else {
		flag.StringVar(&Options.ServerAddress, "a", ":8081", "address and port to run server")
	}

	if Config.DataBaseURI != "" {
		Options.DataBaseURI = Config.DataBaseURI
	} else {
		flag.StringVar(&Options.DataBaseURI, "d", "host=localhost user=gopheruser password=password dbname=gophermat port=5432 sslmode=disable", "database uri")
	}

	if Config.AccrualSystemAddress != "" {
		Options.AccrualSystemAddress = Config.AccrualSystemAddress
	} else {
		flag.StringVar(&Options.AccrualSystemAddress, "r", "http://localhost:8080", "")
	}

	flag.Parse()

	return nil
}
