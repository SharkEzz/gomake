package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SharkEzz/gomake/pkg/gomakefile"
	"github.com/SharkEzz/gomake/pkg/runner/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Dry bool
}

type Runner struct {
	file   *gomakefile.GoMakefile
	config *Config
}

func NewRunner(file *gomakefile.GoMakefile, config *Config) (*Runner, error) {
	if file == nil {
		return nil, errors.New("file is nil")
	}

	if err := env.LoadEnvVariablesFromFiles(file.Dotenv...); err != nil {
		return nil, err

	}
	if err := env.LoadEnvVariablesFromMap(file.Env); err != nil {
		return nil, err
	}

	return &Runner{file, config}, nil
}

func (r *Runner) ExecuteJobByName(jobName string) (int, error) {
	job, ok := r.file.Jobs[jobName]
	if !ok {
		return 0, fmt.Errorf("job '%s' not found", jobName)
	}

	deps := []string{jobName}

	err := r.resolveJobDependencies(jobName, &job, &deps)
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

func (r *Runner) executeJob(jobName string, job *gomakefile.Job) error {
	if checkSkip(jobName, job) {
		// Skip current job
		return nil
	}

	if r.config.Dry {
		log.Infof("executing job '%s' in dry mode", jobName)
	} else {
		log.Infof("executing job '%s'", jobName)
	}

	for _, run := range job.Run {
		// TODO: check env for shell to use
		cmd := exec.Command("sh", "-c", os.ExpandEnv(run))

		if r.config.Dry {
			log.WithField("cmd", cmd.String()).Info("dry run")
			continue
		}

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
				log.WithField("output", line).Infof("output for job '%s'", jobName)
			}
		}
	}

	return nil
}

func checkSkip(jobName string, job *gomakefile.Job) bool {
	if job.SkipIf != "" {
		if _, err := os.Stat(job.SkipIf); !os.IsNotExist(err) {
			// Skip the current job as the checked file / directory already exist
			if !job.Silent {
				log.Infof("skipping job '%s'", jobName)
			}
			return true
		}
	}

	if job.SkipIfNot != "" {
		if _, err := os.Stat(job.SkipIfNot); os.IsNotExist(err) {
			// Skip the current job as the checked file / directory doesn't exist
			if !job.Silent {
				log.Infof("skipping job '%s'", jobName)
			}
			return true
		}
	}

	return false
}

func (r *Runner) resolveJobDependencies(startJobName string, startJob *gomakefile.Job, deps *[]string) error {
OUTER:
	for _, depName := range startJob.Deps {
		loopDep, ok := r.file.Jobs[depName]
		if !ok {
			return fmt.Errorf("dependency '%s' not found for job '%s'", depName, startJobName)
		}

		for _, dep := range *deps {
			if depName == dep {
				log.Warnf("dropped circular reference for job '%s'", dep)

				continue OUTER
			}
		}

		*deps = append(*deps, depName)

		if err := r.resolveJobDependencies(depName, &loopDep, deps); err != nil {
			return err
		}
	}

	return nil
}
