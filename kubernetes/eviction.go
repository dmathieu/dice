package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

type nodeEvicter struct {
	client kubernetes.Interface
	node   *Node
}

func (n *nodeEvicter) Process() error {
	err := n.markNodeUnschedulable()
	if err != nil {
		return err
	}

	pods, err := n.nodePods()
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		err = n.evictPod(&p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *nodeEvicter) markNodeUnschedulable() error {
	node := n.node
	node.Spec.Unschedulable = true
	_, err := n.client.CoreV1().Nodes().Update(node.Node)
	return err
}

func (n *nodeEvicter) evictPod(pod *corev1.Pod) error {
	return n.client.PolicyV1beta1().Evictions(pod.Namespace).Evict(&policy.Eviction{
		ObjectMeta: pod.ObjectMeta,
	})
}

func (n *nodeEvicter) nodePods() (*corev1.PodList, error) {
	nodeName := n.node.ObjectMeta.Name
	return n.client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodeName}).String(),
	})
}

func EvictNode(client kubernetes.Interface, node *Node) error {
	ev := &nodeEvicter{client, node}
	return ev.Process()
}
