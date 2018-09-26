package test

import (
	"testing"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMatchInterfaces(t *testing.T) {
	assert.Implements(t, (*cloudprovider.CloudProvider)(nil), &TestCloudProvider{})
}

func TestName(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Equal(t, "test", p.Name())
}

func TestDelete(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
		},
	}

	p := NewTestCloudProvider()
	assert.Nil(t, p.Delete(node))
}

func TestRefresh(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Nil(t, p.Refresh())
}
