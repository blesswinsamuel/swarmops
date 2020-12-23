package docker

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	dockerStacks = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "swarmops_docker_stack", Help: "Docker stacks"},
		[]string{"name", "namespace", "orchestrator"},
	)
	dockerServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "swarmops_docker_service", Help: "Docker services"},
		[]string{"id", "image", "mode", "name", "ports", "replicas", "stack"},
	)
	dockerMetricsDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "swarmops_docker_metrics_duration_seconds", Help: "Time taken to get docker metrics"},
	)
)

func init() {
	prometheus.MustRegister(dockerStacks)
	prometheus.MustRegister(dockerServices)
	prometheus.MustRegister(dockerMetricsDuration)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func UpdateDockerMetricsLoop(interval time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := updateDockerMetrics()
			if err != nil {
				log.Errorf("failed to update docker metrics: %v", err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func updateDockerMetrics() error {
	startTime := time.Now()
	sc := NewDockerStackCmd()
	stacks, err := sc.Ls()
	if err != nil {
		return err
	}
	dockerStacks.Reset()
	dockerServices.Reset()
	for _, stack := range stacks {
		serviceCount, err := strconv.Atoi(stack.Services)
		if err != nil {
			return err
		}
		dockerStacks.With(prometheus.Labels{
			"name":         stack.Name,
			"namespace":    stack.Namespace,
			"orchestrator": stack.Orchestrator,
		}).Set(float64(serviceCount))

		services, err := sc.Services(stack.Name)
		if err != nil {
			return err
		}
		for _, service := range services {
			dockerServices.With(prometheus.Labels{
				"stack":    stack.Name,
				"id":       service.ID,
				"image":    service.Image,
				"mode":     service.Mode,
				"name":     service.Name,
				"ports":    service.Ports,
				"replicas": service.Replicas,
			}).Set(1)
		}
	}
	dockerMetricsDuration.Observe(time.Since(startTime).Seconds())
	return nil
}
