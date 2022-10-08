package main

import (
	"flag"
	"log"

	"github.com/SharkEzz/gomake/pkg/gomakefile/read"
	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
)

func main() {
	dry := flag.Bool("dry", false, "Set to true to disable the execution of the commands")

	flag.Parse()

	job := flag.Arg(0)

	fileContent, err := read.FindAndReadGoMakefile()

	if err != nil {
		log.Fatalln(err)
	}

	parser, err := parser.Parse(fileContent)
	if err != nil {
		log.Fatalln(err)
	}

	rn, err := runner.NewRunner(parser, *dry)
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
