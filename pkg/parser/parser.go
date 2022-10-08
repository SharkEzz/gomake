package parser

import (
	"github.com/SharkEzz/gomake/pkg/gomakefile"
	"gopkg.in/yaml.v3"
)

func Parse(content []byte) (*gomakefile.GoMakefile, error) {
	gmFile := &gomakefile.GoMakefile{}

	err := yaml.Unmarshal(content, gmFile)
	if err != nil {
		return nil, err
	}

	return gmFile, nil
}
