# Dice

**WIP**

Dice will roll all instances within a kubernetes cluster, using a zero-downtime strategy.

Whenever running the process, it will:

* Flag all running instances as `dice=roll`. All those instances will be rolled.
* Evict all pods from the number of parallel instances required.
* Listen for all pods stopping on a node.
  * When a node has no pods anymore, it will delete it.
* Listen for all new nodes arriving on the cluster.
  * When a new node comes up, it will move on to evicting the pods on another one.

The only supported cloud provider at the moment is AWS. If you need support for another provider, help is more than welcome.  
Dice assumes the cluster has an auto-scaler running, so when a node is deleted, another one can be booted.
