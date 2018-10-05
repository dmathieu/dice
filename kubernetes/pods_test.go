package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func newTrue() *bool {
	b := true
	return &b
}

func TestGetPods(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod",
		},
	}
	daemonSetPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "daemonset-pod",
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					Controller: newTrue(),
					Kind:       "DaemonSet",
				},
			},
		},
	}
	client := fake.NewSimpleClientset(pod, daemonSetPod)

	t.Run("get all pods", func(t *testing.T) {
		pods, err := GetPods(client)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(pods))
	})
}
