package kubernetes

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Pod represents a kubernetes pod
type Pod struct {
	*corev1.Pod
}

// PodOnNode filters the pods to find by a specific node name
func PodOnNode(name string) func(*metav1.ListOptions) {
	return func(o *metav1.ListOptions) {
		o.FieldSelector = mergeSelectors(o.FieldSelector, fmt.Sprintf("spec.nodeName=%q", name))
	}
}

// PodNotTerminated filters the specific pods to find by not being terminated
func PodNotTerminated() func(*metav1.ListOptions) {
	return func(o *metav1.ListOptions) {
		o.FieldSelector = mergeSelectors(o.FieldSelector, fmt.Sprintf("status.phase!=%s,status.phase!=%s", string(corev1.PodSucceeded), string(corev1.PodFailed)))
	}
}

// GetPods lists all pods with optional filters
func GetPods(client kubernetes.Interface, opts ...func(*metav1.ListOptions)) ([]*Pod, error) {
	options := &metav1.ListOptions{}
	for _, opt := range opts {
		opt(options)
	}

	kp, err := client.CoreV1().Pods(metav1.NamespaceAll).List(*options)
	if err != nil {
		return nil, err
	}

	var pods []*Pod
	for _, p := range kp.Items {
		pods = append(pods, &Pod{&p})
	}

	return pods, nil
}

func mergeSelectors(s ...string) string {
	result := []string{}

	for _, v := range s {
		if v == "" {
			continue
		}
		result = append(result, v)
	}

	return strings.Join(result, ",")
}
