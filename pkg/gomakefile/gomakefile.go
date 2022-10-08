package gomakefile

type GoMakefile struct {
	Version string            `yaml:"version"`
	Jobs    map[string]Job    `yaml:"jobs"`
	Env     map[string]string `yaml:"env"`
	Dotenv  []string          `yaml:"dotenv"`
}
