package docker

import (
	"bufio"
	"docker_swarm_gitops/internal/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func execDockerCommand(env map[string]string, args ...string) error {
	log.Infof("Running docker %v", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
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

func execDockerCommandCaptureOutput(env map[string]string, args ...string) ([]string, error) {
	log.Infof("Running docker %v", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	cmdReader := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(cmdReader)
	var lines []string
	for scanner.Scan() {
		m := scanner.Text()
		lines = append(lines, m)
	}
	if err := cmd.Wait(); err != nil {
		return nil, errors.New("cmd.Wait failed")
	}
	return lines, nil
}

type DockerStackCmd struct {
}

func NewDockerStackCmd() *DockerStackCmd {
	return &DockerStackCmd{}
}

func (d *DockerStackCmd) Deploy(cfg *config.StackConfig) error {
	for _, stack := range cfg.Stacks {
		args := []string{"stack", "deploy"}
		for _, f := range stack.ComposeFiles {
			args = append(args, "--compose-file", f)
		}
		if cfg.Prune {
			args = append(args, "--prune")
		}
		if cfg.WithRegistryAuth {
			args = append(args, "--with-registry-auth")
		}
		args = append(args, "--resolve-image", cfg.ResolveImage)
		args = append(args, stack.StackName)

		if err := execDockerCommand(stack.Environment, args...); err != nil {
			return fmt.Errorf("execDockerCommand: %w", err)
		}
	}
	return nil
}

func (d *DockerStackCmd) Remove(cfg *config.StackConfig) error {
	for _, stack := range cfg.Stacks {
		args := []string{"stack", "remove"}
		args = append(args, stack.StackName)
		if err := execDockerCommand(nil, args...); err != nil {
			return fmt.Errorf("execDockerCommand: %w", err)
		}
	}
	return nil
}

type DockerStack struct {
	Name         string
	Namespace    string
	Orchestrator string
	Services     string
}

type DockerService struct {
	ID       string
	Image    string
	Mode     string
	Name     string
	Ports    string
	Replicas string
}

type DockerStackPs struct {
	CurrentState string
	DesiredState string
	Error        string
	ID           string
	Image        string
	Name         string
	Node         string
	Ports        string
}

func (d *DockerStackCmd) Services(stackName string) ([]*DockerService, error) {
	args := []string{"stack", "services"}
	args = append(args, "--format", "{{ json . }}")
	args = append(args, stackName)
	if lines, err := execDockerCommandCaptureOutput(nil, args...); err != nil {
		return nil, fmt.Errorf("execDockerCommand: %w", err)
	} else {
		var objs []*DockerService
		for _, line := range lines {
			obj := DockerService{}
			if err := json.Unmarshal([]byte(line), &obj); err != nil {
				return nil, fmt.Errorf("failed to parse json: %w", err)
			}
			objs = append(objs, &obj)
		}
		return objs, nil
	}
}

func (d *DockerStackCmd) Ls() ([]*DockerStack, error) {
	args := []string{"stack", "ls"}
	args = append(args, "--format", "{{ json . }}")
	if lines, err := execDockerCommandCaptureOutput(nil, args...); err != nil {
		return nil, fmt.Errorf("execDockerCommand: %w", err)
	} else {
		var objs []*DockerStack
		for _, line := range lines {
			obj := DockerStack{}
			if err := json.Unmarshal([]byte(line), &obj); err != nil {
				return nil, fmt.Errorf("failed to parse json: %w", err)
			}
			objs = append(objs, &obj)
		}
		return objs, nil
	}
}
