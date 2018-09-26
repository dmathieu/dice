package test

import (
	"testing"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestMatchInterfaces(t *testing.T) {
	assert.Implements(t, (*cloudprovider.CloudProvider)(nil), &TestCloudProvider{})
}

func TestName(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Equal(t, "test", p.Name())
}

func TestDelete(t *testing.T) {
	node := &kubernetes.Node{}
	p := NewTestCloudProvider()
	assert.Nil(t, p.Delete(node))
}

func TestRefresh(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Nil(t, p.Refresh())
}
