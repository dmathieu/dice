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

func newEvictNodeController(client kube.Interface) *EvictNodeController {
	i := informers.NewSharedInformerFactory(client, controller.NoResyncPeriodFunc())
	controller := NewEvictNodeController(client, i.Core().V1().Nodes(), 1, false)
	controller.nodeListerSynced = alwaysReady
	controller.doneCh = make(chan struct{})
	return controller
}

func TestEvictNodeController(t *testing.T) {
	client := fake.NewSimpleClientset()
	controller := newEvictNodeController(client)

	finishedCh := make(chan struct{})
	doneCh := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(finishedCh)
	}()
	controller.Run(finishedCh, doneCh)
}

func TestEvictNodeControllerDeleteNode(t *testing.T) {
	node := &corev1.Node{}

	client := fake.NewSimpleClientset()
	controller := newEvictNodeController(client)
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
	}

	client := fake.NewSimpleClientset(nonFlaggedNode)
	controller := newEvictNodeController(client)
	controller.addNode(node)

	nodes, err := kubernetes.GetNodes(client)
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
		Spec: corev1.NodeSpec{
			Unschedulable: false,
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				corev1.NodeCondition{Status: corev1.ConditionTrue},
			},
		},
	}
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-node",
		},
	}

	client := fake.NewSimpleClientset(flaggedNode)
	controller := newEvictNodeController(client)
	controller.addNode(node)

	nodes, err := kubernetes.GetNodes(client)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, true, nodes[0].Spec.Unschedulable)
}
