package processors

import (
	"context"
	"errors"
	"testing"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestFlagNodes(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)
	processor := &FlagNodesProcessor{kubeClient: client}

	err := processor.Process(context.Background())
	assert.Nil(t, err)

	nodes, err := kubernetes.GetNodes(client, kubernetes.NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
}

func TestFlagNodesAlreadyFlagged(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)
	processor := &FlagNodesProcessor{kubeClient: client}

	err := kubernetes.FlagNode(client, &kubernetes.Node{Node: node})
	assert.Nil(t, err)

	err = processor.Process(context.Background())
	assert.Equal(t, errors.New("found already flagged nodes. Looks like a roll process is already running"), err)
}
