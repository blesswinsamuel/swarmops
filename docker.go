package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func execDockerCommand(args ...string) error {
	log.Printf("Running docker %v", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	cmdReader := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(cmdReader)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	return cmd.Wait()
}

func runDeploy(cfg *StackConfig) error {
	args := []string{"stack", "deploy"}
	for _, f := range cfg.ComposeFiles {
		args = append(args, "--compose-file", f)
	}
	if cfg.Prune {
		args = append(args, "--prune")
	}
	if cfg.WithRegistryAuth {
		args = append(args, "--with-registry-auth")
	}
	args = append(args, "--resolve-image", cfg.ResolveImage)
	args = append(args, cfg.StackName)
	return execDockerCommand(args...)
}

func runRemove(cfg *StackConfig) error {
	args := []string{"stack", "remove"}
	args = append(args, cfg.StackName)
	return execDockerCommand(args...)
}
