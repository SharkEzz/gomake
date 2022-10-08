package parser

type Job struct {
	Run       []string `yaml:"run"`
	Deps      []string `yaml:"deps"`
	Silent    bool     `yaml:"silent"`
	SkipIf    string   `yaml:"skipIf"`
	SkipIfNot string   `yaml:"skipIfNot"`
}

type GoMakefile struct {
	Version string            `yaml:"version"`
	Jobs    map[string]Job    `yaml:"jobs"`
	Env     map[string]string `yaml:"env"`
	Dotenv  []string          `yaml:"dotenv"`
}
