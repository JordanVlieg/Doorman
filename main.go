package main

import (
	"doorman/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server.StartServer()
}
