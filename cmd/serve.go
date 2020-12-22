package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/blesswinsamuel/swarmops/internal/git"
	"github.com/blesswinsamuel/swarmops/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ServeCommandFlags struct {
	repoDir        string
	privateKeyFile string
	gitRepo        string
	gitBranch      string
	stackFile      string
	syncInterval   time.Duration
	port           string
}

var serveCmdFlags ServeCommandFlags

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ServeExecute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&serveCmdFlags.repoDir, "repo-dir", "", "Repo clone directory")
	serveCmd.Flags().StringVar(&serveCmdFlags.privateKeyFile, "private-key-file", "", "Key persistense directory")
	serveCmd.Flags().StringVar(&serveCmdFlags.gitRepo, "git-repo", "", "SSH git clone repo URL")
	serveCmd.MarkFlagRequired("git-repo")
	serveCmd.Flags().StringVar(&serveCmdFlags.gitBranch, "git-branch", "master", "SSH git branch")
	serveCmd.MarkFlagRequired("git-branch")
	serveCmd.Flags().StringVar(&serveCmdFlags.stackFile, "stack-file", "stack.yaml", "Stack file")
	serveCmd.Flags().DurationVar(&serveCmdFlags.syncInterval, "sync-interval", 5*time.Minute, "Sync Interval")
	serveCmd.Flags().StringVar(&serveCmdFlags.port, "port", "8080", "Server port")
}

func ServeExecute() error {
	c := serveCmdFlags

	git, err := git.NewGit(c.gitRepo, c.gitBranch, c.repoDir, c.privateKeyFile)
	if err != nil {
		log.Fatalf("NewGit: %v", err)
	}

	server := server.NewServer(git, c.repoDir, c.stackFile)

	log.Infof("sync interval: %v", c.syncInterval)
	quit := make(chan struct{})
	if c.syncInterval > 0 {
		go server.BackgroundSync(c.syncInterval, quit)
	}
	h := server.HttpHandler()
	s := &http.Server{Addr: ":" + c.port, Handler: h}

	go func() {
		log.Infof("Server started at port %s", c.port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Infof("Server ListenAndServe failed: %v", err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	close(quit)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Infof("Server Shutdown failed: %v", err)
	}
	return nil
}
