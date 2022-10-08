package gomakefile

type Job struct {
	Run       []string `yaml:"run"`
	Deps      []string `yaml:"deps"`
	Silent    bool     `yaml:"silent"`
	SkipIf    string   `yaml:"skipIf"`
	SkipIfNot string   `yaml:"skipIfNot"`
}
