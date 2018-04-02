package algorithm

import (
	"github.com/kube-usage-scheduler-extender/pkg/metrics"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func makeNodeList(nodeNames []string) []v1.Node {
	result := make([]v1.Node, 0, len(nodeNames))
	for _, nodeName := range nodeNames {
		result = append(result, v1.Node{ObjectMeta: metav1.ObjectMeta{Name: nodeName}})
	}
	return result
}

func makeNodeMetricsInfo(nodeNames []string, metricsInfo []int64) metrics.NodeMetricsInfo {
	nodeMetricsInfo := metrics.NodeMetricsInfo{}
	for index, nodeName := range nodeNames {
		nodeMetricsInfo[nodeName] = metricsInfo[index]
	}
	return nodeMetricsInfo
}

func TestFindOptimizedNodeInCluster(t *testing.T) {
	tests := []struct {
		description      string
		nodes            *v1.NodeList
		nodeMetricsInfo  metrics.NodeMetricsInfo
		expectedNodeList []v1.Node
	}{
		{
			description:      "Test which returns node with least score",
			nodes:            &v1.NodeList{Items: makeNodeList([]string{"node1", "node2"})},
			nodeMetricsInfo:  makeNodeMetricsInfo([]string{"node1", "node2"}, []int64{int64(35), int64(30)}),
			expectedNodeList: makeNodeList([]string{"node2"}),
		},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.expectedNodeList, FindOptimizedNodeInCluster(test.nodes, test.nodeMetricsInfo)) {
			t.Errorf("Test %v failed, Expected the node list to be same", test.description)
		}
	}
}
