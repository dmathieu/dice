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

// EvictNodes finds a set number of random nodes to evict, and drains them of all their pods
func EvictNodes(client kubernetes.Interface, count int) (int, error) {
	nodes, err := GetNodes(client)
	if err != nil {
		return 0, err
	}

	flagged := []*Node{}
	notReady := map[string]*Node{}
	for _, n := range nodes {
		if !n.IsReady() {
			notReady[n.Name] = n
		} else if n.IsFlagged() {
			flagged = append(flagged, n)
		}
	}

	var evictedCount int
	for len(flagged) > 0 && len(notReady) < count {
		i := rand.Intn(len(flagged))
		eNode := flagged[i]
		ev := &nodeEvicter{client, eNode}
		err := ev.Process()
		if err != nil {
			return evictedCount, err
		}
		notReady[eNode.Name] = eNode
		flagged = append(flagged[:i], flagged[i+1:]...)
		evictedCount = evictedCount + 1
	}

	return evictedCount, nil
}
