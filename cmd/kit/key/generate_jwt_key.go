package key

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/app/config"

	"crypto/ed25519"
	"encoding/base64"
)

func GenerateJwtKey(envPath string) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)

	if err != nil {
		panic(err)
	}
	publicKeyString := base64.RawURLEncoding.EncodeToString(publicKey)
	privateKeyString := base64.RawURLEncoding.EncodeToString(privateKey)
	writeKeys(privateKeyString, publicKeyString, envPath)
}

func writeKeys(privateKey string, publicKey string, envPath string) {
	file, err := os.Open(path.Join(config.App().ProjectRoot, envPath))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, _ := io.ReadAll(file)
	envString := string(b)
	regex := regexp.MustCompile("JWT_PRIVATE_KEY.*")
	envString = regex.ReplaceAllString(envString, fmt.Sprintf("JWT_PRIVATE_KEY=%s", privateKey))
	regex = regexp.MustCompile("JWT_PUBLIC_KEY.*")
	envString = regex.ReplaceAllString(envString, fmt.Sprintf("JWT_PUBLIC_KEY=%s", publicKey))
	file.Close()

	if err := os.WriteFile(envPath, []byte(envString), 0644); err != nil {
		panic(err)
	}
}
