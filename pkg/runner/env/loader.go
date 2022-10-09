package env

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariablesFromMap(variables map[string]string) error {
	for key, value := range variables {
		return os.Setenv(key, value)
	}

	return nil
}

func LoadEnvVariablesFromFiles(files ...string) error {
	if len(files) > 0 {
		return godotenv.Load(files...)
	}

	return nil
}
