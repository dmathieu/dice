package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	flagName  = "dice"
	flagValue = "roll"
)

// Node represents a kubernetes node
type Node struct {
	*corev1.Node
}

// IsFlagged checks whether the node has the dice label
func (n *Node) IsFlagged() bool {
	return n.Labels[flagName] == flagValue
}

// IsReady checks whether the node is ready to accept pods
func (n *Node) IsReady() bool {
	if n.Spec.Unschedulable || len(n.Status.Conditions) == 0 {
		return false
	}

	conditionMap := make(map[corev1.NodeConditionType]*corev1.NodeCondition)
	conditions := []corev1.NodeConditionType{corev1.NodeReady}

	for i := range n.Status.Conditions {
		cond := n.Status.Conditions[i]
		conditionMap[cond.Type] = &cond
	}

	for _, validCondition := range conditions {
		if condition, ok := conditionMap[validCondition]; ok {
			if condition.Status == corev1.ConditionFalse {
				return false
			}
		}
	}

	return true
}

// NodeFlagged allows filtering to find only the flagged nodes in GetNodes
func NodeFlagged() func(*metav1.ListOptions) {
	return func(o *metav1.ListOptions) {
		o.LabelSelector = fmt.Sprintf("%s=%s", flagName, flagValue)
	}
}

// GetNodes lists nodes with optional filters
func GetNodes(client kubernetes.Interface, opts ...func(*metav1.ListOptions)) ([]*Node, error) {
	options := &metav1.ListOptions{}
	for _, opt := range opts {
		opt(options)
	}

	kn, err := client.CoreV1().Nodes().List(*options)
	if err != nil {
		return nil, err
	}

	var nodes []*Node
	for _, n := range kn.Items {
		nodes = append(nodes, &Node{n.DeepCopy()})
	}

	return nodes, nil
}

// FindNode finds a specific node by name
func FindNode(client kubernetes.Interface, name string) (*Node, error) {
	node, err := client.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &Node{node}, err
}

// FlagNode marks a specific node as needing to be rolled
func FlagNode(client kubernetes.Interface, node *Node) error {
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}

	node.Labels[flagName] = flagValue
	_, err := client.CoreV1().Nodes().Update(node.Node)
	return err
}
