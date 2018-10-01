package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEvictNode(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)

	err := EvictNode(client, node)
	assert.Nil(t, err)
	assert.Equal(t, true, node.Spec.Unschedulable)
}

func TestEvictNodeWithPods(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
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
	client := fake.NewSimpleClientset(node, podOnNode, podOnOtherNode)

	err := EvictNode(client, node)
	assert.Nil(t, err)

	/*pods, err := client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pods.Items))*/
}
