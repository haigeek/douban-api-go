package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Host      string
	Port      int
	Limit     int
	Cookie    string
	Debug     bool
	BasicUser string
	BasicPass string
}

func Load() Config {
	defaultLimit := 3
	if v := os.Getenv("DOUBAN_API_LIMIT_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			defaultLimit = n
		}
	}

	defaultCookie := os.Getenv("DOUBAN_COOKIE")

	cfg := Config{}
	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "Listen host")
	flag.IntVar(&cfg.Port, "port", 8080, "Listen port")
	flag.IntVar(&cfg.Limit, "limit", defaultLimit, "Search limit for Jellyfin requests")
	flag.StringVar(&cfg.Cookie, "cookie", defaultCookie, "Douban web cookie")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug mode")
	flag.StringVar(&cfg.BasicUser, "basic-user", "", "Basic auth username (enable when both basic-user and basic-pass are set)")
	flag.StringVar(&cfg.BasicPass, "basic-pass", "", "Basic auth password (enable when both basic-user and basic-pass are set)")
	flag.Parse()

	if cfg.Limit < 0 {
		cfg.Limit = 0
	}

	return cfg
}
