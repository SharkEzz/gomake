package read

import (
	"fmt"
	"os"
	"strings"
)

func FindAndReadGoMakefile() ([]byte, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryName := strings.ToLower(entry.Name())

		if !entry.IsDir() && (entryName == "gomakefile.yml" || entryName == "gomakefile.yaml") {
			return os.ReadFile(entry.Name())
		}
	}

	return nil, fmt.Errorf("gomakefile not found")
}
