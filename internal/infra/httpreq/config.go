package httpreq

type Config struct {
	Address       string `long:"address" env:"ADDRESS" description:"Service address" `
	JWTPrivateKey string `long:"jwt-private-key" env:"JWT_PRIVATE_KEY" description:"Path to JWT private key" `
	UploadURL     string `long:"upload-URL" env:"UPLOAD_URL" description:"Path to upload url" `
}
