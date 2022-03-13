package configs

import (
	"Yaratam/internal/infra/http"
	"github.com/jessevdk/go-flags"
	"os"
)

type Config struct {
	HTTP *http.Config `group:"HTTP args" namespace:"http" env-namespace:"YARATAM_HTTP"`
}

func Parse() (*Config, error) {
	var config Config
	p := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)

	_, err := p.ParseArgs(os.Args)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
