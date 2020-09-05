package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/google/logger"
)

var log *logger.Logger

type server struct {
	git       *Git
	docker    *Docker
	mu        sync.Mutex
	baseDir   string
	stackFile string
}

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

	log = logger.Init("LoggerExample", true, false, ioutil.Discard)
	defer log.Close()

	if *gitRepo == "" {
		log.Fatalln("--git-repo should not be empty")
	}
	if *gitBranch == "" {
		log.Fatalln("--git-branch should not be empty")
	}

	git, err := NewGit(*gitRepo, *gitBranch, *repoDir, *privateKeyFile)
	if err != nil {
		log.Fatalf("NewGit: %v", err)
	}

	docker := NewDocker()

	server := &server{
		git:       git,
		docker:    docker,
		baseDir:   *repoDir,
		stackFile: *stackFile,
	}

	log.Infof("sync interval: %v", *syncInterval)
	quit := make(chan struct{})
	if *syncInterval > 0 {
		ticker := time.NewTicker(*syncInterval)
		go func() {
			for {
				select {
				case <-ticker.C:
					err := server.doSync(false)
					if err != nil {
						log.Infof("timed sync failed: %v", err)
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}
	h := server.httpHandler()
	s := &http.Server{Addr: ":" + *port, Handler: h}

	go func() {
		log.Infof("Server started at port %s", *port)
		if err := s.ListenAndServe(); err != nil {
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

func (s *server) doSync(force bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Infoln("Sync started")
	outOfSync, err := s.git.gitSync()
	if err != nil {
		return fmt.Errorf("gitSync: %w", err)
	}
	if outOfSync || force {
		cfg, err := parseConfig(s.baseDir, s.stackFile)
		if err != nil {
			return fmt.Errorf("parseConfig: %w", err)
		}
		err = s.docker.runDeploy(cfg)
		if err != nil {
			return fmt.Errorf("runDeploy: %w", err)
		}
	}
	log.Infoln("Sync completed")
	return nil
}
