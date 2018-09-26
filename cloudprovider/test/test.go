package test

import (
	"sync"

	"github.com/dmathieu/dice/kubernetes"
)

const ProviderName = "test"

func NewTestCloudProvider() *TestCloudProvider {
	return &TestCloudProvider{}
}

// TestCloudProvider is a dummy cloud provider to be used in tests
type TestCloudProvider struct {
	sync.Mutex

	DeletedNodes []*kubernetes.Node
}

func (t *TestCloudProvider) Name() string {
	return ProviderName
}

func (t *TestCloudProvider) Delete(node *kubernetes.Node) error {
	t.DeletedNodes = append(t.DeletedNodes, node)
	return nil
}

func (t *TestCloudProvider) Refresh() error {
	return nil
}
