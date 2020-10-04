package main

import (
	"fmt"
)

// func execDockerCommand(env map[string]string, args ...string) error {
// 	log.Infof("Running docker %v", strings.Join(args, " "))
// 	cmd := exec.Command("docker", args...)
// 	for k, v := range env {
// 		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
// 	}
// 	stdout, _ := cmd.StdoutPipe()
// 	stderr, _ := cmd.StderrPipe()
// 	cmd.Start()

// 	cmdReader := io.MultiReader(stdout, stderr)
// 	scanner := bufio.NewScanner(cmdReader)
// 	for scanner.Scan() {
// 		m := scanner.Text()
// 		fmt.Println(m)
// 	}
// 	return cmd.Wait()
// }

type DockerServiceCmd struct {
}

func NewDockerServiceCmd() *DockerServiceCmd {
	return &DockerServiceCmd{}
}

func (d *DockerServiceCmd) ls(cfg *StackConfig) error {
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