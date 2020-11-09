package config

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type StackConfig struct {
	ResolveImage     string  `yaml:"resolve_image"`
	WithRegistryAuth bool    `yaml:"with_registry_auth"`
	Prune            bool    `yaml:"prune"`
	Stacks           []Stack `yaml:"stacks"`
}
type Stack struct {
	StackName    string            `yaml:"stack_name"`
	ComposeFiles []string          `yaml:"compose_files"`
	Environment  map[string]string `yaml:"environment"`
}

func ParseConfig(repoDir, stackFile string) (*StackConfig, error) {
	var stackConfig StackConfig
	data, err := ioutil.ReadFile(path.Join(repoDir, stackFile))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &stackConfig)
	if err != nil {
		return nil, err
	}
	for _, stack := range stackConfig.Stacks {
		for i, f := range stack.ComposeFiles {
			stack.ComposeFiles[i] = path.Join(repoDir, f)
		}
		stack.Environment["PWD"] = repoDir
	}
	return &stackConfig, nil
}
