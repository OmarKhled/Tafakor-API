package config

import (
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func LoadEnv() {
	// Loading enviroment from .env (if any)
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Env file couldn't be found")
	}
}
