package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodeEvicter(t *testing.T) {
	node := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}}
	client := fake.NewSimpleClientset(node.Node)

	ev := &nodeEvicter{client, node}
	err := ev.Process()
	assert.Nil(t, err)
	assert.Equal(t, true, node.Spec.Unschedulable)
}

func TestNodeEvicterWithPods(t *testing.T) {
	node := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}}
	podOnNode := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-on-node",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName: "node",
		},
	}
	podOnOtherNode := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-on-other-node",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName: "other-node",
		},
	}
	client := fake.NewSimpleClientset(node.Node, podOnNode, podOnOtherNode)

	ev := &nodeEvicter{client, node}
	err := ev.Process()
	assert.Nil(t, err)

	pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pods.Items))
}

func TestEvictNodes(t *testing.T) {
	firstNode := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "first-node",
			Labels: map[string]string{flagName: flagValue},
		},
	}}
	secondNode := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "second-node",
			Labels: map[string]string{flagName: flagValue},
		},
	}}
	thirdNode := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "third-node",
			Labels: map[string]string{flagName: flagValue},
		},
	}}
	client := fake.NewSimpleClientset(firstNode.Node, secondNode.Node, thirdNode.Node)

	t.Run("evicts the set number of nodes", func(t *testing.T) {
		err := EvictNodes(client, 2)
		assert.Nil(t, err)

		nodes, err := GetNodes(client)
		assert.Nil(t, err)
		evictedCount := 0

		for _, n := range nodes {
			if n.Spec.Unschedulable {
				evictedCount = evictedCount + 1
			}
		}
		assert.Equal(t, 2, evictedCount)
	})

	t.Run("ignores nodes that were already evicted", func(t *testing.T) {
		err := EvictNodes(client, 2)
		assert.Nil(t, err)

		ev := &nodeEvicter{client, thirdNode}
		err = ev.Process()
		assert.Nil(t, err)

		nodes, err := GetNodes(client)
		assert.Nil(t, err)
		evictedCount := 0

		for _, n := range nodes {
			if n.Spec.Unschedulable {
				evictedCount = evictedCount + 1
			}
		}
		assert.Equal(t, 2, evictedCount)
	})
}
