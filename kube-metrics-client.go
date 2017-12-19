
package main

import (
	"fmt"
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	resourceclient "k8s.io/metrics/pkg/client/clientset_generated/clientset/typed/metrics/v1beta1"
	"github.com/kube-metrics-test/pkg/metrics"
)

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "/var/run/kubernetes/admin.kubeconfig", "absolute path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// metricsConfig is used for resourceclient
	metricsConfig, err := resourceclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	metricsClient := metrics.NewRESTMetricsClient(metricsConfig)
	leastUtilizedNode, timeStamp, err := metricsClient.GetResourceMetric()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("At %v time, node %v is the least utilized one\n", timeStamp, leastUtilizedNode)
}
