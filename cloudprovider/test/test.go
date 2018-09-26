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
}

func (t *TestCloudProvider) Name() string {
	return ProviderName
}

func (t *TestCloudProvider) Delete(*kubernetes.Node) error {
	return nil
}

func (t *TestCloudProvider) Refresh() error {
	return nil
}
