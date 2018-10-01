package kubernetes

import (
	"math/rand"

	"github.com/golang/glog"
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
	glog.Infof("Evicting node %s", n.node.Name)
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

func EvictNodes(client kubernetes.Interface, count int) error {
	nodes, err := GetNodes(client, NodeFlagged())
	if err != nil {
		return err
	}

	if count > len(nodes) {
		count = len(nodes)
	}

	evicted := map[string]*Node{}

	for len(evicted) < count {
		eNode := nodes[rand.Intn(len(nodes))]
		if evicted[eNode.Name] != nil {
			continue
		}

		ev := &nodeEvicter{client, eNode}
		err := ev.Process()
		if err != nil {
			return err
		}
		evicted[eNode.Name] = eNode
	}

	return nil
}
