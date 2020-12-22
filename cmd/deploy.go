package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blesswinsamuel/swarmops/internal/config"
	"github.com/blesswinsamuel/swarmops/internal/docker"
	"github.com/spf13/cobra"
)

type DeployCmdFlags struct {
	stackFile string
	baseDir   string
}

var deployCmdFlags DeployCmdFlags

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DeployExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVar(&deployCmdFlags.stackFile, "stack-file", "stack.yaml", "Stack file")
	deployCmd.PersistentFlags().StringVar(&deployCmdFlags.baseDir, "base-dir", "./", "Base directory")
}

func DeployExecute() error {
	absPath, err := filepath.Abs(deployCmdFlags.baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	deployCmdFlags.baseDir = absPath

	c := deployCmdFlags
	cfg, err := config.ParseConfig(c.baseDir, c.stackFile)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	d := docker.NewDockerStackCmd()
	err = d.Deploy(cfg)
	if err != nil {
		return fmt.Errorf("failed to run deploy: %w", err)
	}
	return nil
}
