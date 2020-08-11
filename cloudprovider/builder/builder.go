package builder

import (
	"fmt"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/cloudprovider/aws"
	"github.com/dmathieu/dice/cloudprovider/test"
)

// NewCloudProvider creates a new cloud provider by name.
// We need this in it's own package to avoid circular imports
func NewCloudProvider(name string) (cloudprovider.CloudProvider, error) {
	switch name {
	case aws.ProviderName:
		return aws.NewAWSCloudProvider()
	case test.ProviderName:
		return test.NewTestCloudProvider(), nil
	default:
		return nil, fmt.Errorf("Unknown cloud provider %q", name)
	}
}
