package algorithm

import (
	"fmt"
	"github.com/kube-metrics-test/pkg/metrics"
	"k8s.io/api/core/v1"
	"math"
	"math/rand"
	"time"
)

// This file has the scheduling algorithm related functions.

// NodeCost is nodename to cost mapping. cost is combination of
type NodeCostInfo map[string]int64

// random - just some random function which returns values between ranges.
func random(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// PopulateCostForEachNode - As of now, generates random value for each node in cloud. This f(n) could be replaced with
// dynamic cost function.
func PopulateCostForEachNode(nodeList *v1.NodeList) NodeCostInfo {
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

// FindNodeWithLeastCostInCluster takes nodeList, nodeCostInfo and nodeMetricsInfo and returns node with least cost.
func FindOptimizedNodeInCluster(nodeList *v1.NodeList, nodeCostInfo NodeCostInfo, nodeMetricsInfo metrics.NodeMetricsInfo) []v1.Node {
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
