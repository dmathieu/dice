package controllers

import (
	"testing"
	"time"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestOldNodesFlaggerControllerFlagNodes(t *testing.T) {
	now := time.Now()
	firstNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "first node",
			CreationTimestamp: metav1.NewTime(now.Add(0 - time.Hour)),
		},
	}
	secondNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "second node",
			CreationTimestamp: metav1.NewTime(now.Add(0 - time.Second)),
		},
	}
	client := fake.NewSimpleClientset(firstNode, secondNode)
	controller := &OldNodesFlaggerController{
		kubeClient: client,
		interval:   time.Millisecond,
	}

	doneCh := make(chan struct{})
	go controller.Run(doneCh, time.Minute)
	time.Sleep(3 * time.Millisecond)
	close(doneCh)

	nodes, err := kubernetes.GetNodes(client, kubernetes.NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, "first node", nodes[0].ObjectMeta.Name)
}
