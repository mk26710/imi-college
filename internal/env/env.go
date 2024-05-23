package env

import (
	"os"
)

func IsProduction() bool {
	value := os.Getenv("PROD")
	return value == "true"
}

func Addr() string {
	value := os.Getenv("ADDR")
	if len(value) <= 0 {
		panic("ADDR environment variable is unset!")
	}
	return value
}

func DSN() string {
	value := os.Getenv("DB_DSN")
	if len(value) <= 0 {
		panic("DB_DSN environment variable is unset!")
	}
	return value
}
