package main

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type StackConfig struct {
	StackName        string   `yaml:"stack_name"`
	ComposeFiles     []string `yaml:"compose_files"`
	ResolveImage     string   `yaml:"resolve_image"`
	WithRegistryAuth bool     `yaml:"with_registry_auth"`
	Prune            bool     `yaml:"prune"`
}

func parseConfig() (*StackConfig, error) {
	var stackConfig StackConfig
	data, err := ioutil.ReadFile(path.Join(*repoDir, "stack.yml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &stackConfig)
	if err != nil {
		return nil, err
	}
	for i, f := range stackConfig.ComposeFiles {
		stackConfig.ComposeFiles[i] = path.Join(*repoDir, f)
	}
	return &stackConfig, nil
}
