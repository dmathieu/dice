package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateTestNodes generated a set number of valid and flagged nodes for tests
func GenerateTestNodes(count int) []runtime.Object {
	var nodes []runtime.Object
	for i := 1; i <= count; i++ {
		nodes = append(nodes, &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("Node %d", i),
				Labels: map[string]string{flagName: flagValue},
			},
			Spec: corev1.NodeSpec{
				Unschedulable: false,
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					corev1.NodeCondition{Type: corev1.NodeReady, Status: corev1.ConditionTrue},
				},
			},
		})
	}

	return nodes
}
