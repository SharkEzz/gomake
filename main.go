package main

import (
	"flag"
	"log"

	filereader "github.com/SharkEzz/gomake/pkg/file_reader"
	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
)

func main() {
	flag.Parse()

	job := flag.Arg(0)

	fileContent, err := filereader.FindAndReadGoMakefile()
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
