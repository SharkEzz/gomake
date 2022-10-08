package parser_test

import (
	"strings"
	"testing"

	"github.com/SharkEzz/gomake/pkg/parser"
)

const fakeFileContent = `version: '1'

jobs:
  gomake:
    run:
      - go build .
    check: /bin/sh
  
  test_dep:
    deps:
      - gomake
    silent: true
`

func TestParser(t *testing.T) {
	gmfile, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content", err)
	}

	if gmfile.Version != "1" {
		t.Error("Version is not 1")
	}

	if len(gmfile.Jobs) != 2 {
		t.Error("Jobs length is not 2")
	}

	job1, ok := gmfile.Jobs["gomake"]
	if !ok {
		t.Error("Job gomake not found")
	}

	if strings.TrimSpace(job1.Run[0]) != "go build ." {
		t.Error("Job gomake run is not 'go build .'")
	}

	if len(job1.Deps) != 0 {
		t.Error("Job gomake deps length is not 0")
	}

	if job1.Silent {
		t.Error("expected 'test_dep' job to not be silent")
	}

	job2, ok := gmfile.Jobs["test_dep"]
	if !ok {
		t.Error("Job test_dep not found")
	}

	if len(job2.Deps) != 1 {
		t.Error("Job test_dep deps length is not 1")
	}

	if !job2.Silent {
		t.Error("expected 'test_dep' job to be silent")
	}
}
