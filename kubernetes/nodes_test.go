package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNodeIsFlagged(t *testing.T) {
	node := &Node{&corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{},
		},
	}}
	assert.False(t, node.IsFlagged())
	node.Labels[flagName] = flagValue
	assert.True(t, node.IsFlagged())
}

func TestNodeIsReady(t *testing.T) {
	node := &Node{&corev1.Node{}}
	assert.False(t, node.IsReady())

	node.Node.Spec.Unschedulable = true
	assert.False(t, node.IsReady())

	node.Status.Conditions = []corev1.NodeCondition{
		corev1.NodeCondition{Status: corev1.ConditionTrue},
	}
	assert.False(t, node.IsReady())

	node.Node.Spec.Unschedulable = false
	assert.True(t, node.IsReady())

	node.Status.Conditions = []corev1.NodeCondition{
		corev1.NodeCondition{Status: corev1.ConditionTrue},
		corev1.NodeCondition{Status: corev1.ConditionFalse},
	}
	assert.False(t, node.IsReady())
}

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
		assert.Equal(t, 2, len(nodes))
	})

	t.Run("get all flagged nodes", func(t *testing.T) {
		nodes, err := GetNodes(client, NodeFlagged())
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes))
		assert.Equal(t, "flagged-node", nodes[0].Name)
	})

	t.Run("get all non-flagged nodes", func(t *testing.T) {
		nodes, err := GetNodes(client, NodeNotFlagged())
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes))
		assert.Equal(t, "not-flagged-node", nodes[0].Name)
	})
}

func TestFindNode(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-node",
		},
	}
	client := fake.NewSimpleClientset(node)

	t.Run("when the node is found", func(t *testing.T) {
		n, err := FindNode(client, "my-node")
		assert.Nil(t, err)
		assert.Equal(t, &Node{node}, n)
	})

	t.Run("when the node is not found", func(t *testing.T) {
		n, err := FindNode(client, "unknown-node")
		assert.IsType(t, &errors.StatusError{}, err)
		assert.Nil(t, n)
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
	assert.Equal(t, 0, len(nodes))

	err = FlagNode(client, &Node{node})
	assert.Nil(t, err)

	nodes, err = GetNodes(client, NodeFlagged())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, "node", nodes[0].Name)
}
