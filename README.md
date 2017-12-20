# Kube-CaB
Kubernetes's stock scheduler won't take into account the current utilization of nodes in cluster. We intend to solve this problem with Kube-CaB. Kube-CaB is designed to be an add-on service on top of kubernetes to enhance the  scheduler with resource management capabilities. This has two componets:
- Kube metrics client to get node level metrics in a kubernetes cluster.
- A sample kubernetes scheduler extender which returns the node with least cost.

## Build and Run

 Do a git clone of this repo and then run:

```
$ make
```
and then run kube-cab:

```
$ _output/bin/kube-cab --kubeconfig <path to kubeconfig file>
```

Sample Output:

At 2017-12-19 19:39:00 +0530 IST time, node 172.17.0.1 is the least utilized one


