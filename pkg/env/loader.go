package env

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariablesFromMap(variables map[string]string) {
	for key, value := range variables {
		os.Setenv(key, value)
	}
}

func LoadEnvVariablesFromFiles(files ...string) error {
	if len(files) > 0 {
		return godotenv.Load(files...)
	}

	return nil
}
