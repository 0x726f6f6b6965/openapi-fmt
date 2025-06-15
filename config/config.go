package config

type Config struct {
	Remove RemoveConfig `yaml:"remove"`
}

type RemoveConfig struct {
	InputFile  string              `yaml:"input_file"`
	OutputFile string              `yaml:"output_file"`
	Exclude    map[string]struct{} `yaml:"exclude"`
}
