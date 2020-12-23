package cmd

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type SyncCommandFlags struct {
	force bool
	port  string
}

var syncCmdFlags SyncCommandFlags

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Reconcile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := SyncExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&syncCmdFlags.force, "force", false, "Force sync")
	syncCmd.Flags().StringVar(&syncCmdFlags.port, "port", "8080", "Server port")
}

func SyncExecute() error {
	c := syncCmdFlags
	url := fmt.Sprintf("http://localhost:%s/api/sync", c.port)
	if c.force {
		url += "?force=true"
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Got status %s", resp.Status)
	}
	log.Info("Sync successful")
	return nil
}
