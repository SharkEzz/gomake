package runner

import (
	"testing"

	"github.com/SharkEzz/gomake/pkg/parser"
)

const fakeFileContent = `version: '1'

jobs:
  create_test:
    run:
      - echo test > test.txt

  cat_test:
    deps:
      - create_test
    run:
      - cat test.txt

  test_circular_dep:
    deps:
      - test_circular_dep
    run:
      - echo test > test.txt
      - cat test.txt

  test_wrong_dep:
    deps:
      - nonexisting_dep
    run:
      - cat test.txt

  del_test:
    deps:
      - cat_test
    run:
      - rm test.txt
`

var defaultConfig = &Config{
	Dry: true,
}

func TestRunnerWithOneJob(t *testing.T) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		t.Error("Error while creating runner:", err)
	}
	count, err := rn.ExecuteJobByName("create_test")

	if count != 1 {
		t.Error("expected executed job count to be 1:", count)
	}

	if err != nil {
		t.Error("Error while executing job 'test':", err)
	}
}

func TestRunnerWithDependencies(t *testing.T) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		t.Error("Error while creating runner:", err)
	}
	count, err := rn.ExecuteJobByName("del_test")

	if count != 3 {
		t.Error("expected executed job count to be 3:", count)
	}

	if err != nil {
		t.Error("Error while executing job 'test':", err)
	}
}

func TestRunnerWithCircularDependency(t *testing.T) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		t.Error("Error while creating runner:", err)
	}
	count, err := rn.ExecuteJobByName("test_circular_dep")
	if count != 1 {
		t.Error("expected executed job count to be 1:", count)
	}

	if err != nil {
		t.Error("expected circular reference to be dropped: ", err)
	}
}

func TestRunnerWithNonExistingJob(t *testing.T) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		t.Error("Error while creating runner:", err)
	}
	count, err := rn.ExecuteJobByName("dummy")
	if count != 0 {
		t.Error("count should be 0:", count)
	}
	if err == nil || err.Error() != "job 'dummy' not found" {
		t.Error("wrong error", err)
	}
}

func TestRunnerWithNonExistingDependency(t *testing.T) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		t.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		t.Error("Error while creating runner:", err)
	}
	count, err := rn.ExecuteJobByName("test_wrong_dep")
	if count != 0 {
		t.Error("count should be 0:", count)
	}
	if err == nil || err.Error() != "dependency 'nonexisting_dep' not found for job 'test_wrong_dep'" {
		t.Error("wrong error", err)
	}
}

func BenchmarkComputeExecutionOrder(b *testing.B) {
	file, err := parser.Parse([]byte(fakeFileContent))
	if err != nil {
		b.Error("Error while parsing content:", err)
	}

	rn, err := NewRunner(file, defaultConfig)
	if err != nil {
		b.Error("Error while creating runner:", err)
	}

	job := rn.file.Jobs["del_test"]
	deps := []string{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := rn.resolveJobDependencies("del_test", &job, &deps)
		if err != nil {
			b.Error(err)
		}
	}
}
