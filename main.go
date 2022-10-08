package main

import (
	"flag"
	"log"
	"os"

	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
)

func main() {
	file := flag.String("file", "", "The path to the GoMakeFile to load")
	job := flag.String("job", "", "The name of the job to execute")

	flag.Parse()

	var fileContent []byte
	var err error

	if *file != "" {
		fileContent, err = os.ReadFile(*file)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fileContent, err = os.ReadFile("./GoMakeFile")
		if err != nil {
			log.Fatalln(err)
		}
	}

	parser, err := parser.Parse(fileContent)
	if err != nil {
		log.Fatalln(err)
	}

	rn, err := runner.NewRunner(parser)
	if err != nil {
		log.Fatalln(err)
	}

	if *job != "" {
		_, err := rn.ExecuteJobByName(*job)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if err := rn.ExecuteAllJobs(); err != nil {
		log.Fatalln(err)
	}
}
