
/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	resourceclient "k8s.io/metrics/pkg/client/clientset_generated/clientset/typed/metrics/v1beta1"
	//"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/labels"
	"time"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	config1, _ := resourceclient.NewForConfig(config)
	metricsClient := NewRESTMetricsClient(config1)
	metricsClient.GetResourceMetric()
}


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
	//GetResourceMetric(resource v1.ResourceName, namespace string, selector labels.Selector) (NodeMetricsInfo, time.Time, error)
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
		fmt.Print("Cost of system")
		return nil, time.Time{}, fmt.Errorf("no metrics returned from heapster")
	}

	//res := make(NodeMetricsInfo, len(metrics.Items))

	for _, m := range metrics.Items {
		fmt.Printf("Node Usage %v", m.Name, m.Usage.Cpu())
	}

	timestamp := metrics.Items[0].Timestamp.Time

	return nil, timestamp, nil
}


func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
