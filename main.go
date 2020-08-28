package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	gitRepo      = flag.String("git-repo", "", "SSH git clone repo URL")
	gitBranch    = flag.String("git-branch", "master", "SSH git branch")
	keysDir      = flag.String("keys-dir", "", "Key persistense directory")
	repoDir      = flag.String("repo-dir", "", "Repo clone directory")
	stackFile    = flag.String("stack-file", "stack.yaml", "Stack file")
	syncInterval = flag.Duration("sync-interval", 5*time.Minute, "Sync Interval")
	port         = flag.String("port", "8080", "Server port")
)
var mx sync.Mutex

func main() {
	flag.Parse()

	if *gitRepo == "" {
		log.Fatalln("--git-repo should not be empty")
	}
	if *gitBranch == "" {
		log.Fatalln("--git-branch should not be empty")
	}

	ticker := time.NewTicker(*syncInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := doSync()
				if err != nil {
					log.Printf("Initial sync failed: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	h := httpHandler()
	server := &http.Server{Addr: ":" + *port, Handler: h}

	go func() {
		log.Printf("Server started at port %s", *port)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server ListenAndServe failed: %v", err)
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
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown failed: %v", err)
	}
}

func doSync() error {
	mx.Lock()
	defer mx.Unlock()
	log.Println("Sync started")
	_, err := gitSync()
	if err != nil {
		return fmt.Errorf("gitSync: %w", err)
	}
	// if changed {
	cfg, err := parseConfig()
	if err != nil {
		return fmt.Errorf("parseConfig: %w", err)
	}
	err = runDeploy(cfg)
	if err != nil {
		return fmt.Errorf("runDeploy: %w", err)
	}
	// }
	log.Println("Sync completed")
	return nil
}
