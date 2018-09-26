package test

import "github.com/dmathieu/dice/cloudprovider"

const ProviderName = "test"

// TestCloudProvider is a dummy cloud provider to be used in tests
type TestCloudProvider struct {
}

func (t *TestCloudProvider) Name() string {
	return ProviderName
}

func (t *TestCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	group := make([]cloudprovider.NodeGroup, 0)

	return group
}

func (t *TestCloudProvider) Refresh() error {
	return nil
}

type TestNodeGroup struct {
}
