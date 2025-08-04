package flame
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	clusterLastActivity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kubernetes_cluster_last_activity_timestamp",
			Help: "Timestamp of the last activity in the Kubernetes cluster (based on audit log modification time)",
		},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(clusterLastActivity)
}

func main() {
	auditLogPath := os.Getenv("AUDIT_LOG_PATH")
	pushgatewayURL := os.Getenv("PUSHGATEWAY_URL")
	jobName := os.Getenv("JOB_NAME")

	log.Printf("Starting Kubernetes Cluster Activity Monitor")
	log.Printf("Audit log path: %s", auditLogPath)
	log.Printf("Pushgateway URL: %s", pushgatewayURL)
	log.Printf("Job name: %s", jobName)

	// Check cluster last activity
	if err := checkClusterActivity(auditLogPath); err != nil {
		log.Printf("Error checking cluster activity: %v", err)
		// clusterActivityCheck.WithLabelValues("error").Inc()
	} 
	// else {
	// 	clusterActivityCheck.WithLabelValues("success").Inc()
	// }

	// Push metrics to Pushgateway
	if err := pushMetricsToPushgateway(pushgatewayURL); err != nil {
		log.Fatalf("Error pushing metrics to Pushgateway: %v", err)
	}
	log.Println("Cluster activity check completed successfully")
}

// checkClusterActivity 
func checkClusterActivity(auditLogPath string) error {
	// Get file info
	fileInfo, err := os.Stat(auditLogPath)
	if err != nil {
		return fmt.Errorf("failed to access audit log file: %w", err)
	}

	// Get the last modification time
	lastModTime := fileInfo.ModTime()

	// Set the metric to the timestamp
	clusterLastActivity.Set(float64(lastModTime.Unix()))

	log.Printf("Cluster last activity detected at: %v (Unix timestamp: %d)",
		lastModTime.Format(time.RFC3339), lastModTime.Unix())

	return nil
}

// pushMetricsToPushgateway
func pushMetricsToPushgateway(pushgatewayURL) error {
	pusher := push.New(pushgatewayURL).
		Collector(clusterLastActivity).

	if err := pusher.Push(); err != nil {
		return fmt.Errorf("failed to push metrics to pushgateway: %w", err)
	}

	log.Printf("Successfully pushed metrics to pushgateway at %s", pushgatewayURL)
	return nil
}
