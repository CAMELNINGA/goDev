package bot

type Config struct {
	Token             string `env:"TOKEN,required"`
	APIEndPoint       string `env:"API_ENDPOINT"`
	RunMode           string `env:"RUN_WITH" envDefault:"direct"`
	Debug             bool   `env:"DEBUG" envDefault:"true"`
	InstructionFileID string `env:"INSTRUCTION_FILE"`
}
