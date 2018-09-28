package controllers

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller"
)

func newEvictNodeController(kClient kube.Interface) *EvictNodeController {
	i := informers.NewSharedInformerFactory(kClient, controller.NoResyncPeriodFunc())
	controller := NewEvictNodeController(kClient, i.Core().V1().Nodes())
	controller.nodeListerSynced = alwaysReady
	return controller
}

func TestEvictNodeController(t *testing.T) {
	kClient := fake.NewSimpleClientset()
	controller := newEvictNodeController(kClient)

	doneCh := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(doneCh)
	}()
	controller.Run(doneCh)
}

func TestEvictNodeControllerDeleteNode(t *testing.T) {
	node := &corev1.Node{}

	kClient := fake.NewSimpleClientset()
	controller := newEvictNodeController(kClient)
	controller.deleteNode(node)
}
