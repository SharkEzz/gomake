package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/SharkEzz/gomake/pkg/parser"
)

type Runner struct {
	file *parser.GoMakefile
}

func NewRunner(gmFile *parser.GoMakefile) (*Runner, error) {
	if gmFile == nil {
		return nil, errors.New("gmFile is nil")
	}

	return &Runner{gmFile}, nil
}

func (r *Runner) ExecuteJobByName(jobName string) (int, error) {
	job, ok := r.file.Jobs[jobName]
	if !ok {
		return 0, fmt.Errorf("job '%s' not found", jobName)
	}

	deps := []string{jobName}

	err := r.computeExecutionOrder(jobName, &job, &deps)
	if err != nil {
		return 0, err
	}

	for i := len(deps) - 1; i >= 0; i-- {
		depName := deps[i]
		job := r.file.Jobs[depName]

		err := r.executeJob(depName, &job)
		if err != nil {
			return 0, err
		}
	}

	return len(deps), nil
}

func (r *Runner) executeJob(jobName string, job *parser.Job) error {
	if checkSkip(jobName, job) {
		// Skip current job
		return nil
	}

	log.Printf("executing job '%s'", jobName)

	// TODO: check env for shell to use
	cmd := exec.Command("sh", "-c", os.ExpandEnv(strings.TrimSpace(job.Run)))

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if job.Silent {
		return nil
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr != "" {
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			log.Printf("output for job '%s': %s", jobName, line)
		}
	}

	return nil
}

func checkSkip(jobName string, job *parser.Job) bool {
	if job.SkipIf != "" {
		if _, err := os.Stat(job.SkipIf); !os.IsNotExist(err) {
			// Skip the current job as the checked file / directory already exist
			if !job.Silent {
				log.Printf("skipping job '%s'\n", jobName)
			}
			return true
		}
	}

	if job.SkipIfNot != "" {
		if _, err := os.Stat(job.SkipIfNot); os.IsNotExist(err) {
			// Skip the current job as the checked file / directory doesn't exist
			if !job.Silent {
				log.Printf("skipping job '%s'\n", jobName)
			}
			return true
		}
	}

	return false
}

func (r *Runner) ExecuteAllJobs() error {
	for jobName, job := range r.file.Jobs {
		err := r.executeJob(jobName, &job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) computeExecutionOrder(startJobName string, startJob *parser.Job, deps *[]string) error {
OUTER:
	for _, depName := range startJob.Deps {
		loopDep, ok := r.file.Jobs[depName]
		if !ok {
			return fmt.Errorf("dependency '%s' not found for job '%s'", depName, startJobName)
		}

		for _, dep := range *deps {
			if depName == dep {
				log.Printf("dropped circular reference for job '%s'", dep)

				continue OUTER
			}
		}

		*deps = append(*deps, depName)

		if err := r.computeExecutionOrder(depName, &loopDep, deps); err != nil {
			return err
		}
	}

	return nil
}
