package server

import (
	"swarmops/internal/config"
	"swarmops/internal/docker"
	"swarmops/internal/git"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	git       *git.Git
	docker    *docker.DockerStackCmd
	mu        sync.Mutex
	baseDir   string
	stackFile string
}

func NewServer(git *git.Git, repoDir string, stackFile string) *Server {
	docker := docker.NewDockerStackCmd()

	return &Server{
		git:       git,
		docker:    docker,
		baseDir:   repoDir,
		stackFile: stackFile,
	}
}

func (s *Server) HttpHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/sync", s.syncHandler)
	r.HandleFunc("/api/docker/stacks", s.dockerStackListHandler)
	r.HandleFunc("/api/docker/services", s.dockerServiceListHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
	return r
}

func (s *Server) syncHandler(w http.ResponseWriter, r *http.Request) {
	forceStr := r.URL.Query().Get("force")
	force := false
	if forceStr != "" {
		var err error
		force, err = strconv.ParseBool(forceStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid param force: %v", err)
			return
		}
	}
	fmt.Println(force)

	err := s.doSync(force)
	if err != nil {
		log.Infof("Sync failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Sync failed: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}

func (s *Server) dockerStackListHandler(w http.ResponseWriter, r *http.Request) {
	stacks, err := docker.NewDockerStackCmd().Ls()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	err = json.NewEncoder(w).Encode(stacks)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
}

func (s *Server) dockerServiceListHandler(w http.ResponseWriter, r *http.Request) {
	stackName := r.URL.Query().Get("stack")
	stacks, err := docker.NewDockerStackCmd().Services(stackName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	err = json.NewEncoder(w).Encode(stacks)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")
}

func (s *Server) BackgroundSync(syncInterval time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(syncInterval)
	for {
		select {
		case <-ticker.C:
			err := s.doSync(false)
			if err != nil {
				log.Infof("timed sync failed: %v", err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (s *Server) doSync(force bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Infoln("Sync started")
	outOfSync, err := s.git.Sync()
	if err != nil {
		return fmt.Errorf("gitSync: %w", err)
	}
	if outOfSync || force {
		cfg, err := config.ParseConfig(s.baseDir, s.stackFile)
		if err != nil {
			return fmt.Errorf("parseConfig: %w", err)
		}
		err = s.docker.Deploy(cfg)
		if err != nil {
			return fmt.Errorf("runDeploy: %w", err)
		}
	}
	log.Infoln("Sync completed")
	return nil
}
