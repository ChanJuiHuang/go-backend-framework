package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	filename := ""
	flag.StringVar(&filename, "env", ".env", "dot env file path")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	envString := string(b)
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	regex := regexp.MustCompile("JWT_PRIVATE_KEY.*")
	envString = regex.ReplaceAllString(envString, fmt.Sprintf("JWT_PRIVATE_KEY=%s", base64.RawURLEncoding.EncodeToString(privateKey)))
	regex = regexp.MustCompile("JWT_PUBLIC_KEY.*")
	envString = regex.ReplaceAllString(envString, fmt.Sprintf("JWT_PUBLIC_KEY=%s", base64.RawURLEncoding.EncodeToString(publicKey)))

	if err := os.WriteFile(filename, []byte(envString), 0644); err != nil {
		panic(err)
	}
}
