package main

import (
	"context"
	"docker_swarm_gitops/internal/git"
	"docker_swarm_gitops/internal/server"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		repoDir        = flag.String("repo-dir", "", "Repo clone directory")
		privateKeyFile = flag.String("private-key-file", "", "Key persistense directory")
		gitRepo        = flag.String("git-repo", "", "SSH git clone repo URL")
		gitBranch      = flag.String("git-branch", "master", "SSH git branch")
		stackFile      = flag.String("stack-file", "stack.yaml", "Stack file")
		syncInterval   = flag.Duration("sync-interval", 5*time.Minute, "Sync Interval")
		port           = flag.String("port", "8080", "Server port")
	)
	flag.Parse()

	if *gitRepo == "" {
		log.Fatalln("--git-repo should not be empty")
	}
	if *gitBranch == "" {
		log.Fatalln("--git-branch should not be empty")
	}

	git, err := git.NewGit(*gitRepo, *gitBranch, *repoDir, *privateKeyFile)
	if err != nil {
		log.Fatalf("NewGit: %v", err)
	}

	server := server.NewServer(git, *repoDir, *stackFile)

	log.Infof("sync interval: %v", *syncInterval)
	quit := make(chan struct{})
	if *syncInterval > 0 {
		go server.BackgroundSync(*syncInterval, quit)
	}
	h := server.HttpHandler()
	s := &http.Server{Addr: ":" + *port, Handler: h}

	go func() {
		log.Infof("Server started at port %s", *port)
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
}
