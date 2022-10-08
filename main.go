package main

import (
	"flag"
	"log"

	"github.com/SharkEzz/gomake/pkg/loader"
	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
)

func main() {
	flag.Parse()

	job := flag.Arg(0)

	fileContent, err := loader.LoadGoMakefileContent()
	if err != nil {
		log.Fatalln(err)
	}

	parser, err := parser.Parse(fileContent)
	if err != nil {
		log.Fatalln(err)
	}

	rn, err := runner.NewRunner(parser)
	if err != nil {
		log.Fatalln(err)
	}

	if job != "" {
		_, err := rn.ExecuteJobByName(job)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if err := rn.ExecuteAllJobs(); err != nil {
		log.Fatalln(err)
	}
}
