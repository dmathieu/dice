package controllers

import (
	"testing"
	"time"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestEvictNodeControllerNewNoFlagged(t *testing.T) {
	nonFlaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "non-flagged-node",
		},
	}
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-node",
		},
		Spec: corev1.NodeSpec{
			Unschedulable: false,
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				corev1.NodeCondition{Status: corev1.ConditionTrue},
			},
		},
	}

	kClient := fake.NewSimpleClientset(nonFlaggedNode)
	controller := newEvictNodeController(kClient)
	controller.addNode(node)

	nodes, err := kubernetes.GetNodes(kClient)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, false, nodes[0].Spec.Unschedulable)
}

func TestEvictNodeControllerNew(t *testing.T) {
	flaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "flagged-node",
			Labels: map[string]string{"dice": "roll"},
		},
	}
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-node",
		},
		Spec: corev1.NodeSpec{
			Unschedulable: false,
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				corev1.NodeCondition{Status: corev1.ConditionTrue},
			},
		},
	}

	kClient := fake.NewSimpleClientset(flaggedNode)
	controller := newEvictNodeController(kClient)
	controller.addNode(node)

	nodes, err := kubernetes.GetNodes(kClient)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, true, nodes[0].Spec.Unschedulable)
}
