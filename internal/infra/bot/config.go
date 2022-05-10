package bot

type Config struct {
	Token             string `env:"GAME_TG_BOT_TOKEN,required"`
	APIEndPoint       string `env:"GAME_TG_BOT_API_ENDPOINT"`
	ProxyURL          string `env:"GAME_TG_BOT_PROXY_URL"`
	ProxyLogin        string `env:"GAME_TG_BOT_PROXY_LOGIN"`
	ProxyPass         string `env:"GAME_TG_BOT_PROXY_PASS"`
	RunMode           string `env:"GAME_TG_BOT_RUN_WITH" envDefault:"direct"`
	Debug             bool   `env:"GAME_TG_BOT_DEBUG" envDefault:"true"`
	InstructionFileID string `env:"GAME_TG_BOT_INSTRUCTION_FILE"`
}
