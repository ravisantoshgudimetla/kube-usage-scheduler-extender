package main

import (
	"flag"
	"github.com/kube-metrics-test/pkg/algorithm"
	"k8s.io/client-go/tools/clientcmd"
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
	// Start the extender HTTP service. Need to make sure to avoid race condition.
	algorithm.StartHttpServer(config)
}
