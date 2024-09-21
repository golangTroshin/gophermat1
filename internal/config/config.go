package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var Options struct {
	ServerAddress          string `env:"RUN_ADDRESS" envDefault:":8081"`
	DataBaseURI            string `env:"DATABASE_URI" envDefault:"host=localhost user=gopheruser password=password dbname=gophermat port=5432 sslmode=disable"`
	AccrualSystemAddress   string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	AccrualSystemWorkers   int    `env:"ACCRUAL_SYSTEM_WORKERS" envDefault:"2"`
	AccrualSystemInterval  int    `env:"ACCRUAL_SYSTEM_INTERVAL" envDefault:"5"`
	AccrualSystemRateLimit int    `env:"ACCRUAL_SYSTEM_RATE_LIMIT" envDefault:"5"`
}

func ParseFlags() error {
	if err := env.Parse(&Options); err != nil {
		return err
	}

	flag.StringVar(&Options.ServerAddress, "a", Options.ServerAddress, "address and port to run server")
	flag.StringVar(&Options.DataBaseURI, "b", Options.DataBaseURI, "base result url")
	flag.StringVar(&Options.AccrualSystemAddress, "f", Options.AccrualSystemAddress, "accrual system address")
	flag.IntVar(&Options.AccrualSystemWorkers, "w", Options.AccrualSystemWorkers, "accrual system number of workers")
	flag.IntVar(&Options.AccrualSystemInterval, "i", Options.AccrualSystemInterval, "accrual system interval of requests")
	flag.IntVar(&Options.AccrualSystemRateLimit, "r", Options.AccrualSystemRateLimit, "accrual system rate limit")

	flag.Parse()

	return nil
}
