package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/stack/loader"
	"github.com/docker/cli/cli/command/stack/options"
	"github.com/docker/cli/cli/command/stack/swarm"
	"github.com/docker/cli/cli/flags"
	"gopkg.in/yaml.v2"
)

var baseDir = "./sample-stack"

type StackConfig struct {
	StackName        string   `yaml:"stack_name"`
	ComposeFiles     []string `yaml:"compose_files"`
	ResolveImage     string   `yaml:"resolve_image"`
	WithRegistryAuth bool     `yaml:"with_registry_auth"`
	Prune            bool     `yaml:"prune"`
}

func parseConfig() (*StackConfig, error) {
	var stackConfig StackConfig
	data, err := ioutil.ReadFile(path.Join(baseDir, "stack.yml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &stackConfig)
	if err != nil {
		return nil, err
	}
	return &stackConfig, nil
}

func main() {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dockerCli.Initialize(flags.NewClientOptions())

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = runDeploy(dockerCli, cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = runRemove(dockerCli, cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func runDeploy(dockerCli *command.DockerCli, cfg *StackConfig) error {
	composeFiles := []string{}
	for _, f := range cfg.ComposeFiles {
		composeFiles = append(composeFiles, path.Join(baseDir, f))
	}
	opts := options.Deploy{
		Composefiles:     composeFiles,
		Namespace:        cfg.StackName,
		ResolveImage:     cfg.ResolveImage,
		SendRegistryAuth: cfg.WithRegistryAuth,
		Prune:            cfg.Prune,
	}
	config, err := loader.LoadComposefile(dockerCli, opts)
	if err != nil {
		return err
	}
	err = swarm.RunDeploy(dockerCli, opts, config)
	if err != nil {
		return err
	}
	return nil
}

func runRemove(dockerCli *command.DockerCli, cfg *StackConfig) error {
	opts := options.Remove{
		Namespaces: []string{cfg.StackName},
	}
	err := swarm.RunRemove(dockerCli, opts)
	if err != nil {
		return err
	}
	return nil
}
