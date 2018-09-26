package test

import (
	"testing"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/stretchr/testify/assert"
)

func TestMatchInterfaces(t *testing.T) {
	assert.Implements(t, (*cloudprovider.CloudProvider)(nil), &TestCloudProvider{})
	assert.Implements(t, (*cloudprovider.NodeGroup)(nil), &TestNodeGroup{})
}

func TestNodeGroups(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Equal(t, 0, len(p.NodeGroups()))

	p.InsertNodeGroup(&TestNodeGroup{id: "my-node"})
	assert.Equal(t, 1, len(p.NodeGroups()))
	assert.Equal(t, "my-node", p.NodeGroups()[0].Id())
}

func TestName(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Equal(t, "test", p.Name())
}

func TestRefresh(t *testing.T) {
	p := NewTestCloudProvider()
	assert.Nil(t, p.Refresh())
}
