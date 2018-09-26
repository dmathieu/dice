package builder

import (
	"fmt"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/cloudprovider/test"
)

// NewCloudProvider creates a new cloud provider by name.
// We need this in it's own package to avoid circular imports
func NewCloudProvider(name string) (cloudprovider.CloudProvider, error) {
	switch name {
	case test.ProviderName:
		return test.NewTestCloudProvider(), nil
	default:
		return nil, fmt.Errorf("Unknown cloud provider %q", name)
	}
}
