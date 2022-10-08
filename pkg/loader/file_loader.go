package loader

import (
	"fmt"
	"os"
	"strings"
)

func LoadGoMakefileContent() ([]byte, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.ToLower(entry.Name()) == "gomakefile" {
			return os.ReadFile(entry.Name())
		}
	}

	return nil, fmt.Errorf("gomakefile not found")
}
