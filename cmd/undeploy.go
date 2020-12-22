package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blesswinsamuel/swarmops/internal/config"
	"github.com/blesswinsamuel/swarmops/internal/docker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type UndeployCmdFlags struct {
	stackFile string
	baseDir   string
}

var undeployCmdFlags UndeployCmdFlags

var undeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := UndeployExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(undeployCmd)

	undeployCmd.PersistentFlags().StringVar(&undeployCmdFlags.stackFile, "stack-file", "stack.yaml", "Stack file")
	undeployCmd.PersistentFlags().StringVar(&undeployCmdFlags.baseDir, "base-dir", "./", "Base directory")
}

func UndeployExecute() error {
	absPath, err := filepath.Abs(undeployCmdFlags.baseDir)
	if err != nil {
		log.Errorf("failed to resolve path: %s", err)
	}
	undeployCmdFlags.baseDir = absPath

	c := undeployCmdFlags
	cfg, err := config.ParseConfig(c.baseDir, c.stackFile)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %s", err)
	}
	d := docker.NewDockerStackCmd()
	err = d.Remove(cfg)
	if err != nil {
		return fmt.Errorf("Failed to run undeploy: %s", err)
	}
	return nil
}
