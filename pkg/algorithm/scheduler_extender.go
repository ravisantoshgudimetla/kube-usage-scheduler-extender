package algorithm

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kube-metrics-test/pkg/metrics"
	"k8s.io/api/core/v1"
	restclient "k8s.io/client-go/rest"
	resourceclient "k8s.io/metrics/pkg/client/clientset_generated/clientset/typed/metrics/v1beta1"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// NodeCost is nodename to cost mapping. cost is combination of
type NodeCostInfo map[string]int64

// All these types could be exported from scheduler once we import k8s.io/kubernetes as scheduler is not separated as
// repo.
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

// random - just some random function which returns values between ranges.
func random(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// populateCostForEachNode - As of now, generates random value for each node in cloud. This f(n) could be replaced with
// dynamic cost function.
func populateCostForEachNode(nodeList *v1.NodeList) NodeCostInfo {
	nodeCloudCostInfo := NodeCostInfo{}
	rand.Seed(time.Now().Unix())
	// Replace this to get cost for each node in cloud.
	for _, node := range nodeList.Items {
		nodeCloudCostInfo[node.Name] = random(int64(1), int64(100))
	}
	return nodeCloudCostInfo
}

func findOptimizedNode(nodeTotalCostInfo NodeCostInfo, nodeList *v1.NodeList) []v1.Node {
	minCost := int64(math.MaxInt64)
	var nodeNeeded string
	for node, cost := range nodeTotalCostInfo {
		if cost < minCost {
			minCost = cost
			nodeNeeded = node
		}
	}
	fmt.Printf("The node %v is the node with minimum cost %v", nodeNeeded, minCost)
	var neededNodeList []v1.Node
	for _, node := range nodeList.Items {
		if node.Name == nodeNeeded {
			neededNodeList = append(neededNodeList, node)
		}
	}
	return neededNodeList
}

// findNodeWithLeastCostInCluster takes nodeList, nodeCostInfo and nodeMetricsInfo and returns node with least cost.
func findOptimizedNodeInCluster(nodeList *v1.NodeList, nodeCostInfo NodeCostInfo, nodeMetricsInfo metrics.NodeMetricsInfo) []v1.Node {
	nodeTotalCost := NodeCostInfo{}
	var cloudCost, cpuUtil int64
	var ok bool
	// The optimization function is sum of cpuUtil and cloudCost should be minimum.
	for _, node := range nodeList.Items {
		cloudCost, ok = nodeCostInfo[node.Name]
		if !ok {
			continue
		}
		cpuUtil, ok = nodeMetricsInfo[node.Name]
		if !ok {
			continue
		}
		fmt.Printf("Cloud Cost is %v, cpuUtil Cost is %v\n", cloudCost, cpuUtil)
		totalCost := cloudCost + cpuUtil
		nodeTotalCost[node.Name] = totalCost
	}
	return findOptimizedNode(nodeTotalCost, nodeList)
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
	// Populate cost for each node from cloud.
	nodeCostInfo := populateCostForEachNode(args.Nodes)
	// Find the totalCost of each node.
	nodesWithLeastCost := findOptimizedNodeInCluster(args.Nodes, nodeCostInfo, nodeUtilInfo)
	return &ExtenderFilterResult{Nodes: &v1.NodeList{Items: nodesWithLeastCost}, NodeNames: nil, FailedNodes: nil}
}

func StartHttpServer(config *restclient.Config) {
	router := mux.NewRouter()
	router.HandleFunc("/scheduler/filter", func(w http.ResponseWriter, r *http.Request) { schedule(w, r, config) }).Methods("POST")
	http.ListenAndServe(":9000", router)
}
