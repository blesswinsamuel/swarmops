package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
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

type composeFile struct {
	Configs map[string]struct {
		File string `yaml:"file"`
	} `yaml:"configs"`
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
	h := sha256.New()
	for _, stack := range stackConfig.Stacks {
		for i, f := range stack.ComposeFiles {
			data, err := ioutil.ReadFile(path.Join(repoDir, f))
			if err != nil {
				return nil, err
			}
			var composeFile composeFile
			err = yaml.Unmarshal([]byte(data), &composeFile)
			if err != nil {
				log.Warn("Ignoring error: %v", err)
			}
			for key, c := range composeFile.Configs {
				configFile := strings.ReplaceAll(c.File, "${PWD}", repoDir)
				f, err := os.Open(configFile)
				if err != nil {
					return nil, fmt.Errorf("failed to open config file '%s': %v", configFile, err)
				}
				if _, err := io.Copy(h, f); err != nil {
					return nil, fmt.Errorf("failed to hash config file '%s': %v", configFile, err)
				}
				stack.Environment[key+"_hash"] = hex.EncodeToString(h.Sum(nil))[:7]
			}
			stack.ComposeFiles[i] = path.Join(repoDir, f)
		}
		stack.Environment["PWD"] = repoDir
		fmt.Println(stack.Environment)
	}
	return &stackConfig, nil
}
