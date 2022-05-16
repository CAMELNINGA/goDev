package http

type Config struct {
	Address       string `short:"a" long:"address" env:"ADDRESS" description:"Service address" `
	JWTPrivateKey string `long:"jwt-private-key" env:"JWT_PRIVATE_KEY" description:"Path to JWT private key" `
	UploadURL     string
}
