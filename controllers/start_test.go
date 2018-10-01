package controllers

import (
	"errors"
	"testing"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestStartControllerFlagNodes(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)
	controller := &StartController{kubeClient: client}

	err := controller.Run()
	assert.Nil(t, err)

	nodes, err := kubernetes.GetNodes(client, kubernetes.NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
}

func TestStartControllerFlagNodesAlreadyFlagged(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)
	controller := &StartController{kubeClient: client}

	err := kubernetes.FlagNode(client, &kubernetes.Node{Node: node})
	assert.Nil(t, err)

	err = controller.Run()
	assert.Equal(t, errors.New("found already flagged nodes. Looks like a roll process is already running"), err)
}
