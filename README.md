[![Build Status](https://travis-ci.org/ravisantoshgudimetla/kube-cab.svg?branch=master)](https://travis-ci.org/ravisantoshgudimetla/kube-cab)
[![Go Report Card](https://goreportcard.com/badge/github.com/ravisantoshgudimetla/kube-usage-scheduler-extender)](https://goreportcard.com/badge/github.com/ravisantoshgudimetla/kube-usage-scheduler-extender)
# kube-usage-scheduler-extender
Kubernetes's stock scheduler won't take into account the current utilization of nodes in cluster while making scheduling decisions. We intend to solve this problem with kube-usage-scheduler-extender. It is designed to be an add-on service on top of kubernetes to enhance scheduler with resource management capabilities. It has two components:
- Kube metrics client to get node level metrics in a kubernetes cluster. As of now, the code returns the node with least CPU utilization but this could be extended to any resource(like memory, GPU etc). We are relying on metrics-server to get information related to node usage(https://github.com/kubernetes-incubator/metrics-server).
- A sample kubernetes scheduler extender which returns the node with least cost in the cloud. As of now, the algorithm is very simple with static hardcoding of node costs in the cloud.

## Architecture

![](https://github.com/ravisantoshgudimetla/kube-CaB/blob/master/Kube-CaB%20Arch.png)

### Flow
- Kubernetes scheduler has the concept of scheduler extender where it sends the request to HTTP server before binding a pod to node if extender is enabled. The request sent to HTTP server includes pod and nodelist that are filtered. 
- The extender as of now has a filter function which further filters the nodes from nodelist supplied. The filtering is based on computation algorithm. The computation algorithm as of now is a very simple algorithm which talks to metrics server and get the node information for CPU, memory usage.
- Once the filtering happens, the HTTP server responds back with node which has the least CPU utilization.

## Build and Run

 - Make sure that metrics-server is running as deployment and getting node level metrics. You can test this using:
 kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq. This should return information related to all the node along with current usage on the nodes.
 - Do a git clone of this repo and then run:

```
$ make
```
and then run kube-usage-scheduler-extender:

```
$ _output/bin/kube-ext --kubeconfig <path to kubeconfig file>
```

## Note:
This is not ready for production use yet. But give it a spin and provide some feedback.

