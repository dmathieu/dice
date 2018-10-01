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

type Node struct {
	*corev1.Node
}

func (n *Node) IsFlagged() bool {
	return n.Labels[flagName] == flagValue
}

func (n *Node) IsReady() bool {
	if n.Spec.Unschedulable || len(n.Status.Conditions) == 0 {
		return false
	}

	for _, c := range n.Status.Conditions {
		if c.Status != corev1.ConditionTrue {
			return false
		}
	}

	return true
}

func NodeFlagged() func(*metav1.ListOptions) {
	return func(o *metav1.ListOptions) {
		o.LabelSelector = fmt.Sprintf("%s=%s", flagName, flagValue)
	}
}

func NodeNotFlagged() func(*metav1.ListOptions) {
	return func(o *metav1.ListOptions) {
		o.LabelSelector = fmt.Sprintf("%s!=%s", flagName, flagValue)
	}
}

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

func FindNode(client kubernetes.Interface, name string) (*Node, error) {
	node, err := client.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &Node{node}, err
}

func FlagNode(client kubernetes.Interface, node *Node) error {
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}

	node.Labels[flagName] = flagValue
	_, err := client.CoreV1().Nodes().Update(node.Node)
	return err
}
