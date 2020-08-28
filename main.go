package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	gitRepo   = flag.String("git-repo", "", "SSH git clone repo URL")
	gitBranch = flag.String("git-branch", "master", "SSH git branch")
	keysDir   = flag.String("keys-dir", "", "Key persistense directory")
	repoDir   = flag.String("repo-dir", "", "Repo clone directory")
	stackFile = flag.String("stack-file", "stack.yaml", "Stack file")
	port      = flag.String("port", "8080", "Server port")
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
	go func() {
		err := doSync()
		if err != nil {
			log.Printf("Initial sync failed: %v", err)
		}
	}()
	h := httpHandler()
	log.Printf("Server started at port %s", *port)
	http.ListenAndServe(":"+*port, h)
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
