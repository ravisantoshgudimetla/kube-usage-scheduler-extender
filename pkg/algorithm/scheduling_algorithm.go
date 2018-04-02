package algorithm

import (
	"fmt"
	"github.com/kube-usage-scheduler-extender/pkg/metrics"
	"k8s.io/api/core/v1"
	"math"
)

// This file has the scheduling algorithm related functions.

// NodeCostInfo is nodename to cost mapping. Cost is combination of cpu utilization and cost of node in cloud.
type NodeCostInfo map[string]int64

// findOptimizedNode finds the optimized node in the cluster.
func findOptimizedNode(nodeTotalCostInfo NodeCostInfo, nodeList *v1.NodeList) []v1.Node {
	minCost := int64(math.MaxInt64)
	var nodeNeeded string
	for node, cost := range nodeTotalCostInfo {
		if cost < minCost {
			minCost = cost
			nodeNeeded = node
		}
	}
	var neededNodeList = make([]v1.Node, 0)
	if minCost == int64(math.MaxInt64) {
		return neededNodeList
	}
	fmt.Printf("The node %v is the node with minimum cpu utilization of %v", nodeNeeded, minCost)
	for _, node := range nodeList.Items {
		if node.Name == nodeNeeded {
			neededNodeList = append(neededNodeList, node)
		}
	}
	return neededNodeList
}

// FindOptimizedNodeInCluster takes nodeList, nodeCostInfo and nodeMetricsInfo and returns node with least cost.
func FindOptimizedNodeInCluster(nodeList *v1.NodeList, nodeMetricsInfo metrics.NodeMetricsInfo) []v1.Node {
	nodeTotalCost := NodeCostInfo{}
	var cloudCost, cpuUtil int64
	cloudCost = 0
	var ok bool
	// The optimization function is sum of cpuUtil and cloudCost should be minimum.
	for _, node := range nodeList.Items {
		/*cloudCost, ok = nodeCostInfo[node.Name]
		if !ok {
			continue // This node will not be taken into consideration for algorithm.
		}*/
		cpuUtil, ok = nodeMetricsInfo[node.Name]
		if !ok {
			continue // This node will not be taken into consideration for algorithm.
		}
		fmt.Printf("CpuUtil Cost is %v\n", cpuUtil)
		// This could cause a buffer overflow. Need to have a check.
		totalCost := cloudCost + cpuUtil
		nodeTotalCost[node.Name] = totalCost
	}
	return findOptimizedNode(nodeTotalCost, nodeList)
}
