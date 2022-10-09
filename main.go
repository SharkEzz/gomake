package main

import (
	"flag"

	"github.com/SharkEzz/gomake/pkg/gomakefile/read"
	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
	log "github.com/sirupsen/logrus"
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

	config := &runner.Config{
		Dry: *dry,
	}

	rn, err := runner.NewRunner(parser, config)
	if err != nil {
		log.Fatalln(err)
	}

	if job == "" {
		log.Fatalln("No job specified")
	}

	_, err = rn.ExecuteJobByName(job)
	if err != nil {
		log.Fatalln(err)
	}
}
