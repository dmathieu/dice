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

func GetNodes(client kubernetes.Interface, opts ...func(*metav1.ListOptions)) (*corev1.NodeList, error) {
	options := &metav1.ListOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return client.CoreV1().Nodes().List(*options)
}

func FlagNode(client kubernetes.Interface, node *corev1.Node) error {
	if node.ObjectMeta.Labels == nil {
		node.ObjectMeta.Labels = map[string]string{}
	}

	node.ObjectMeta.Labels[flagName] = flagValue
	_, err := client.CoreV1().Nodes().Update(node)
	return err
}
