package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kube-cab/pkg/algorithm"
	"github.com/kube-cab/pkg/metrics"
	"k8s.io/api/core/v1"
	restclient "k8s.io/client-go/rest"
	resourceclient "k8s.io/metrics/pkg/client/clientset_generated/clientset/typed/metrics/v1beta1"
	"net/http"
	"strings"
)

// TODO: All these types could be exported from scheduler once we import k8s.io/kubernetes as scheduler is not separated as
// repo.

// FailedNodesMap is needed by HTTP server response.
type FailedNodesMap map[string]string

// ExtenderArgs represents the arguments needed by the extender to filter/prioritize
// nodes for a pod.
type ExtenderArgs struct {
	// Pod being scheduled
	Pod v1.Pod
	// List of candidate nodes where the pod can be scheduled; to be populated
	// only if ExtenderConfig.NodeCacheCapable == false
	Nodes *v1.NodeList
	// List of candidate node names where the pod can be scheduled; to be
	// populated only if ExtenderConfig.NodeCacheCapable == true
	NodeNames *[]string
}

// ExtenderFilterResult stores the result from extender to be sent as response.
type ExtenderFilterResult struct {
	// Filtered set of nodes where the pod can be scheduled; to be populated
	// only if ExtenderConfig.NodeCacheCapable == false
	Nodes *v1.NodeList
	// Filtered set of nodes where the pod can be scheduled; to be populated
	// only if ExtenderConfig.NodeCacheCapable == true
	NodeNames *[]string
	// Filtered out nodes where the pod can't be scheduled and the failure messages
	FailedNodes FailedNodesMap
	// Error message indicating failure
	Error string
}

// schedule does the actual scheduling of pods on the given node.
func schedule(w http.ResponseWriter, r *http.Request, config *restclient.Config) {
	// Get the list of nodes from scheduler in request and sort them based on the cost.
	// Iterate over the list of nodes which has least CPU.
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	encoder := json.NewEncoder(w)
	var args ExtenderArgs
	if strings.Contains(r.URL.Path, "filter") {
		if err := decoder.Decode(&args); err != nil {
			http.Error(w, "Decode error", http.StatusBadRequest)
		}
		resp := filter(&args, config)
		if err := encoder.Encode(resp); err != nil {
			http.Error(w, "Encode error", http.StatusBadRequest)
		}
	}
}

// filter takes the metrics config input and returns
func filter(args *ExtenderArgs, config *restclient.Config) *ExtenderFilterResult {
	// Get the CPU utilization for each node. It returns nodename and CPU value.
	metricsConfig, err := resourceclient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	metricsClient := metrics.NewRESTMetricsClient(metricsConfig)

	// Populate the nodeUtilInfo which has node name to CPU utilization metrics.
	// TODO: Need to use timestamp for caching later.
	nodeUtilInfo, timeStamp, err := metricsClient.GetResourceMetric()
	if err != nil {
		// return all the nodes here as metrics are not yet available.
		fmt.Println("Returning nodes here")
		return &ExtenderFilterResult{Nodes: args.Nodes, NodeNames: nil, FailedNodes: nil}
	}
	fmt.Printf("At %v time, %v is the node utilization map\n", timeStamp, nodeUtilInfo)
	// Populate cost for each node from cloud. This step will be replaced later.
	//nodeCostInfo := algorithm.PopulateCostForEachNode(args.Nodes)
	// Find the totalCost of each node.
	//nodesWithLeastCost := algorithm.FindOptimizedNodeInCluster(args.Nodes, nodeCostInfo, nodeUtilInfo)
	nodesWithLeastCost := algorithm.FindOptimizedNodeInCluster(args.Nodes, nodeUtilInfo)
	if len(nodesWithLeastCost) > 0 {
		return &ExtenderFilterResult{Nodes: &v1.NodeList{Items: nodesWithLeastCost}, NodeNames: nil, FailedNodes: nil}
	}
	return &ExtenderFilterResult{Nodes: args.Nodes, NodeNames: nil, FailedNodes: nil}
}

// startHttpServer starts the HTTP server needed for scheduler.
func startHttpServer(config *restclient.Config) {
	router := mux.NewRouter()
	router.HandleFunc("/scheduler/filter", func(w http.ResponseWriter, r *http.Request) { schedule(w, r, config) }).Methods("POST")
	http.ListenAndServe(":9000", router)
}
