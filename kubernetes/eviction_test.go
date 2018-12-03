package kubernetes

import (
	"fmt"
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
	t.Run("evicts the set number of nodes", func(t *testing.T) {
		client := fake.NewSimpleClientset(GenerateTestNodes(10)...)
		count, err := EvictNodes(client, 3)
		assert.Nil(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("does not fail if the concurrency is higher than the number of nodes", func(t *testing.T) {
		client := fake.NewSimpleClientset(GenerateTestNodes(10)...)
		count, err := EvictNodes(client, 100)
		assert.Nil(t, err)
		assert.Equal(t, 10, count)
	})

	t.Run("ignores nodes that were already evicted", func(t *testing.T) {
		nodes := GenerateTestNodes(10)
		client := fake.NewSimpleClientset(nodes...)
		ev := &nodeEvicter{client, &Node{nodes[0].(*corev1.Node)}}
		err := ev.Process()
		assert.Nil(t, err)

		count, err := EvictNodes(client, 3)
		assert.Nil(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("reduces count if there are non-ready nodes", func(t *testing.T) {
		nodes := GenerateTestNodes(10)
		nodes = append(nodes, &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "Non-ready node",
				Labels: map[string]string{flagName: flagValue},
			},
		})

		client := fake.NewSimpleClientset(nodes...)
		count, err := EvictNodes(client, 3)
		assert.Nil(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("does nothing if there are too many non-ready nodes", func(t *testing.T) {
		nodes := GenerateTestNodes(10)
		for i := 1; i <= 10; i++ {
			nodes = append(nodes, &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name:   fmt.Sprintf("Non-ready node %d", i),
					Labels: map[string]string{flagName: flagValue},
				},
			})
		}

		client := fake.NewSimpleClientset(nodes...)
		count, err := EvictNodes(client, 3)
		assert.Nil(t, err)
		assert.Equal(t, 0, count)
	})
}
