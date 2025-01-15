package models

type APIKey struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	In    string `yaml:"in"`
}

type APIProvider struct {
	APIKey APIKey `yaml:"apiKey"`
}

type APIConfig struct {
	APIProvider APIProvider `yaml:"apiProvider"`
	API         string      `yaml:"api"`
}

type Config struct {
	APIs              APIConfig `yaml:"apis"`
	Instruction       string    `yaml:"instruction"`
	MaxIterationSteps int       `yaml:"max_iteration_steps"`
}
