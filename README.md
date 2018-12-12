# Dice

Dice will roll all instances within a kubernetes cluster, using a zero-downtime strategy.

Whenever running the process, it will:

* Flag all running instances as `dice=roll`. All those instances will be rolled.
* Evict all pods from the number of parallel instances required.
* Listen for all pods stopping on a node.
  * When a node has no pods anymore, it will delete it.
* Listen for all new nodes arriving on the cluster.
  * When a new node comes up, it will move on to evicting the pods on another one.

Dice assumes the cluster has an auto-scaler running, so when a node is deleted, another one can be booted.

## Supported Providers

Only AWS is currently supported

## Usage

### In Cluster

In order to run dice within the cluster, kubernetes needs to be able to delete AWS instances. That can be done with the following IAM policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:TerminateInstances"
            ],
            "Resource": "*"
        }
    ]
}
```

#### As a one-off

Running dice as a one-off is a good use case when you have changed the boot
configuration for your instances for examples, and you need new ones with the
appropriate config.

```
kubeclt apply -f examples/dice-aws.yaml
```

#### Regularly roll instances

You may want to regularly roll instances if they are too old. This allows rolling out the fleet on a regular cadence.

```
kubeclt apply -f examples/loop-aws.yaml
```

### Out of Cluster

You can run dice from your own machines (good for testing, but it shouldn't be used on production workloads).

```
go get -u github.com/dmathieu/dice
dice run -c aws
```

Note: you can run the permanent loop out of the cluster with the `dice loop -c aws` command.
This is not recommended however, as you would then need to have the process running permanently.
