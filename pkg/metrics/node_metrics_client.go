package metrics

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	resourceclient "k8s.io/metrics/pkg/client/clientset_generated/clientset/typed/metrics/v1beta1"
	"time"

	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// NodeMetricsInfo contains pod metric values as a map from pod names to
// metric values (the metric values are expected to be the metric as a milli-value)
type NodeMetricsInfo map[string]int64

func NewRESTMetricsClient(resourceClient resourceclient.NodeMetricsesGetter) MetricsClient {
	return &restMetricsClient{
		&resourceMetricsClient{resourceClient},
	}
}

// restMetricsClient is a client which supports fetching
// metrics from both the resource metrics API and the
// custom metrics API.
type restMetricsClient struct {
	*resourceMetricsClient
}

// MetricsClient knows how to query a remote interface to retrieve container-level
// resource metrics as well as pod-level arbitrary metrics
type MetricsClient interface {
	// GetResourceMetric gets the given resource metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetResourceMetric() (NodeMetricsInfo, time.Time, error)
}

type resourceMetricsClient struct {
	client resourceclient.NodeMetricsesGetter
}

// GetResourceMetric gets the given resource metric (and an associated oldest timestamp)
// for all pods matching the specified selector in the given namespace
func (c *resourceMetricsClient) GetResourceMetric() (NodeMetricsInfo, time.Time, error) {
	metrics, err := c.client.NodeMetricses().List(metav1.ListOptions{})
	if err != nil {
		fmt.Printf("unable to fetch metrics from API: %v", err)
		return nil, time.Time{}, fmt.Errorf("unable to fetch metrics from API: %v", err)
	}

	if len(metrics.Items) == 0 {
		return nil, time.Time{}, fmt.Errorf("no metrics returned from metric-server")
	}
	// In the list of nodes, find the one which has least CPU
	getNodeWithLeastCPU(metrics)
	timestamp := metrics.Items[0].Timestamp.Time
	return nil, timestamp, nil
}

func getNodeWithLeastCPU(metrics *v1beta1.NodeMetricsList) {
	for _, m := range metrics.Items {
		// Implement logic to find node which has least CPU utilization.
		fmt.Printf("Node %v, CPU Usage %v, Memory Usage %v", m.Name, m.Usage.Cpu(), m.Usage.Memory())
	}
}