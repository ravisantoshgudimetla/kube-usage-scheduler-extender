# Kube-‘CaB’  — Add-on services in kubernetes to make it ‘Cheap’ & ‘Balanced’
Kubernetes's stock scheduler won't take into account the current utilization of nodes in cluster while making scheduling decisions. We intend to solve this problem with Kube-CaB. Kube-CaB is designed to be an add-on service on top of kubernetes to enhance scheduler with resource management capabilities and cost related optimizations in cloud. It has two components:
- Kube metrics client to get node level metrics in a kubernetes cluster. As of now, the code returns the node with least CPU utilization but this could be extended to any resource(like memory, GPU etc).
- A sample kubernetes scheduler extender which returns the node with least cost in the cloud. As of now, the algorithm is very simple with static hardcoding of node costs in the cloud.

## Architecture

![](https://github.com/ravisantoshgudimetla/kube-CaB/blob/master/Kube-CaB%20Arch.png)

## Build and Run

 Do a git clone of this repo and then run:

```
$ make
```
and then run kube-cab:

```
$ _output/bin/kube-cab --kubeconfig <path to kubeconfig file>
```



