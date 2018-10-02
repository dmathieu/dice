package test

import (
	"sync"

	"github.com/dmathieu/dice/kubernetes"
)

// ProviderName is the name of this provider
const ProviderName = "test"

// NewTestCloudProvider instantiates a new CloudProvider
func NewTestCloudProvider() *CloudProvider {
	return &CloudProvider{}
}

// CloudProvider is a dummy cloud provider to be used in tests
type CloudProvider struct {
	sync.Mutex

	DeletedNodes []*kubernetes.Node
}

// Name is the name of that cloud provider
func (t *CloudProvider) Name() string {
	return ProviderName
}

// Delete deletes the specified node
func (t *CloudProvider) Delete(node *kubernetes.Node) error {
	t.DeletedNodes = append(t.DeletedNodes, node)
	return nil
}
