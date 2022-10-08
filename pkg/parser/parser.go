package parser

import (
	"gopkg.in/yaml.v3"
)

func Parse(content []byte) (*GoMakefile, error) {
	gmFile := &GoMakefile{}

	err := yaml.Unmarshal(content, gmFile)
	if err != nil {
		return nil, err
	}

	return gmFile, nil
}
