package main

import (
	"flag"
	"log"
	"os"

	"github.com/SharkEzz/gomake/pkg/parser"
	"github.com/SharkEzz/gomake/pkg/runner"
)

func main() {
	flag.Parse()

	file := flag.Arg(0)
	job := flag.Arg(1)

	var fileContent []byte
	var err error

	if file != "" {
		fileContent, err = os.ReadFile(file)
	} else {
		fileContent, err = os.ReadFile("./GoMakeFile")
	}
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
