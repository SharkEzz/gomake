package parser

type Job struct {
	Run    string   `yaml:"run"`
	Deps   []string `yaml:"deps"`
	Silent bool     `yaml:"silent"`
	Check  string   `yaml:"check"`
}

type GoMakefile struct {
	Version string         `yaml:"version"`
	Jobs    map[string]Job `yaml:"jobs"`
}
