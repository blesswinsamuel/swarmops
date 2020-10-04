package docker

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

func (d *DockerServiceCmd) ls() error {
	args := []string{"service", "ls"}
	args = append(args, "--format", "{{ json . }}")

	if err := execDockerCommand(nil, args...); err != nil {
		return fmt.Errorf("execDockerCommand: %w", err)
	}
	return nil
}
