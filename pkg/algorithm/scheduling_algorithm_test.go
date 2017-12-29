package algorithm

import (
	"testing"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kube-cab/pkg/metrics"
	"reflect"

)


func makeNodeList(nodeNames []string) []v1.Node {
	result := make([]v1.Node, 0, len(nodeNames))
	for _, nodeName := range nodeNames {
		result = append(result, v1.Node{ObjectMeta: metav1.ObjectMeta{Name: nodeName}})
	}
	return result
}

func makeNodeCost(nodeNames []string, costs []int64) NodeCostInfo {
	nodeCostInfo := NodeCostInfo{}
	for index, nodeName := range nodeNames {
		nodeCostInfo[nodeName] = costs[index]
	}
	return nodeCostInfo
}

func makeNodeMetricsInfo(nodeNames []string, metricsInfo []int64) metrics.NodeMetricsInfo {
	nodeMetricsInfo := metrics.NodeMetricsInfo{}
	for index, nodeName := range nodeNames {
		nodeMetricsInfo[nodeName] = metricsInfo[index]
	}
	return nodeMetricsInfo
}

func TestFindOptimizedNodeInCluster(t *testing.T) {
	tests := []struct{
		description string
		nodes *v1.NodeList
		nodeCostInfo NodeCostInfo
		nodeMetricsInfo metrics.NodeMetricsInfo
		expectedNodeList []v1.Node
	}{
		{
			description: "Test which returns node with least score",
			nodes: &v1.NodeList{Items: makeNodeList([]string{"node1", "node2"})},
			nodeCostInfo: makeNodeCost([]string{"node1", "node2"}, []int64{int64(35), int64(30)}),
			nodeMetricsInfo: makeNodeMetricsInfo([]string{"node1", "node2"}, []int64{int64(35), int64(30)}),
			expectedNodeList: makeNodeList([]string{"node2"}),
		},
		{
			description: "Test which returns node with least score but if node cloud cost is not " +
				"available it will not be taken into consideration, so node2 won't be taken into consideration",
			nodes: &v1.NodeList{Items: makeNodeList([]string{"node1", "node2"})},
			nodeCostInfo: makeNodeCost([]string{"node1"}, []int64{int64(35)}),
			nodeMetricsInfo: makeNodeMetricsInfo([]string{"node1", "node2"}, []int64{int64(35), int64(30)}),
			expectedNodeList: makeNodeList([]string{"node1"}),
		},
		{
			description: "Test which returns node with least score but if metrics are not " +
				"available that node will not be taken into consideration, so node2 won't be taken into consideration",
			nodes: &v1.NodeList{Items: makeNodeList([]string{"node1", "node2"})},
			nodeCostInfo: makeNodeCost([]string{"node1", "node2"}, []int64{int64(35), int64(30)}),
			nodeMetricsInfo: makeNodeMetricsInfo([]string{"node1"}, []int64{int64(35)}),
			expectedNodeList: makeNodeList([]string{"node1"}),
		},


	}
	for _, test:= range tests{
		if !reflect.DeepEqual(test.expectedNodeList, FindOptimizedNodeInCluster(test.nodes, test.nodeCostInfo, test.nodeMetricsInfo)) {
			t.Errorf("Test %v failed, Expected the node list to be same", test.description)
		}
	}
}
