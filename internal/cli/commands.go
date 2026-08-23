package cli

import (
	"flag"
	"fmt"
)

type Config struct {
	Address string
	DBPath  string
	Seed    bool
}

func DefaultConfig() Config {
	return Config{Address: ":8080", DBPath: "venue.db", Seed: true}
}

func Parse(args []string) (Config, error) {
	config := DefaultConfig()
	flags := flag.NewFlagSet("venue-server", flag.ContinueOnError)
	flags.StringVar(&config.Address, "address", config.Address, "HTTP listen address")
	flags.StringVar(&config.DBPath, "db", config.DBPath, "bbolt database path")
	flags.BoolVar(&config.Seed, "seed", config.Seed, "load deterministic catalog data")
	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	if config.Address == "" {
		return Config{}, fmt.Errorf("address is required")
	}
	if config.DBPath == "" {
		return Config{}, fmt.Errorf("db path is required")
	}
	return config, nil
}

func Usage() string {
	return "venue-server -address :8080 -db venue.db -seed=true"
}

func IsSeedEnabled(config Config) bool {
	return config.Seed
}
