package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/blesswinsamuel/swarmops/internal/config"
	"github.com/blesswinsamuel/swarmops/internal/docker"
	"github.com/blesswinsamuel/swarmops/internal/git"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	reconcileDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "swarmops_reconcile_duration_seconds", Help: "Time taken to receoncile"},
		[]string{"force"},
	)
)

func init() {
	prometheus.MustRegister(reconcileDuration)
}

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
	r.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
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

func (s *Server) BackgroundSync(syncInterval time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(syncInterval)
	for {
		select {
		case <-ticker.C:
			err := s.doSync(false)
			if err != nil {
				log.Errorf("timed sync failed: %v", err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (s *Server) doSync(force bool) error {
	startTime := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Debug("Sync started")
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
	log.Debug("Sync completed")
	reconcileDuration.WithLabelValues(strconv.FormatBool(force)).Observe(time.Since(startTime).Seconds())
	return nil
}
