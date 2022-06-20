package configs

import (
	"Yaratam/internal/infra/bot"
	"Yaratam/internal/infra/http"
	"Yaratam/internal/infra/httpreq"
	"Yaratam/internal/infra/postgres"
	"Yaratam/pkg/logging"
	"github.com/jessevdk/go-flags"
	"os"
)

type Config struct {
	Logger   *logging.Config  `group:"Logger args" namespace:"logger" env-namespace:"YARATAM_LOGGER"`
	Postgres *postgres.Config `group:"Postgres args" namespace:"postgres" env-namespace:"YARATAM_POSTGRES"`
	HTTP     *http.Config     `group:"HTTP args" namespace:"http" env-namespace:"YARATAM_HTTP"`
	HTTPReq  *httpreq.Config  `group:"HTTP args" namespace:"httpreq" env-namespace:"YARATAM_HTTPREQ"`
	Telegram *bot.Config      `group:" Telegram args" namespace:"tgbot" env-namespace:"YRATAM_TG_BOT"`
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
