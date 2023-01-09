package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/ChanJuiHuang/go-backend-framework/app/http"
)

// @title Example API
// @version 1.0
// @schemes http https
// @host localhost:8080
func main() {
	http.RunServer()
}
