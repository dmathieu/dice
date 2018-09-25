package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetNodes(t *testing.T) {
	flaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "flagged-node",
			Labels: map[string]string{flagName: flagValue},
		},
	}
	notFlaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "not-flagged-node",
		},
	}
	client := fake.NewSimpleClientset(flaggedNode, notFlaggedNode)

	t.Run("get all nodes", func(t *testing.T) {
		nodes, err := GetNodes(client)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(nodes.Items))
	})

	t.Run("get all flagged nodes", func(t *testing.T) {
		nodes, err := GetNodes(client, NodeFlagged())
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))
		assert.Equal(t, "flagged-node", nodes.Items[0].ObjectMeta.Name)
	})

	t.Run("get all non-flagged nodes", func(t *testing.T) {
		nodes, err := GetNodes(client, NodeNotFlagged())
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))
		assert.Equal(t, "not-flagged-node", nodes.Items[0].ObjectMeta.Name)
	})
}

func TestFlagNode(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}
	client := fake.NewSimpleClientset(node)

	nodes, err := GetNodes(client, NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(nodes.Items))

	err = FlagNode(client, node)
	assert.Nil(t, err)

	nodes, err = GetNodes(client, NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))
	assert.Equal(t, "node", nodes.Items[0].ObjectMeta.Name)
}
