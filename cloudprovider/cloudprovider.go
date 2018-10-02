package cloudprovider

import "github.com/dmathieu/dice/kubernetes"

// CloudProvider is a generic interface used by all providers
type CloudProvider interface {
	// Name returns name of the cloud provider.
	Name() string

	// Delete finds and deletes the specified node
	Delete(*kubernetes.Node) error
}
