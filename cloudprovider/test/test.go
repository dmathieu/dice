package test

import (
	"sync"

	"github.com/dmathieu/dice/cloudprovider"
)

const ProviderName = "test"

func NewTestCloudProvider() *TestCloudProvider {
	return &TestCloudProvider{
		groups: make(map[string]cloudprovider.NodeGroup),
	}
}

// TestCloudProvider is a dummy cloud provider to be used in tests
type TestCloudProvider struct {
	sync.Mutex
	groups map[string]cloudprovider.NodeGroup
}

func (t *TestCloudProvider) Name() string {
	return ProviderName
}

func (t *TestCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	t.Lock()
	defer t.Unlock()

	group := make([]cloudprovider.NodeGroup, 0)

	for _, g := range t.groups {
		group = append(group, g)
	}

	return group
}

func (t *TestCloudProvider) InsertNodeGroup(g *TestNodeGroup) {
	t.Lock()
	defer t.Unlock()
	t.groups[g.Id()] = g
}

func (t *TestCloudProvider) Refresh() error {
	return nil
}

type TestNodeGroup struct {
	sync.Mutex
	id string
}

func (g *TestNodeGroup) Id() string {
	g.Lock()
	defer g.Unlock()
	return g.id
}
