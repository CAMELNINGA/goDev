package http

type Config struct {
	Address string `short:"a" long:"address" env:"ADDRESS" description:"Service address" required:"yes"`
}
