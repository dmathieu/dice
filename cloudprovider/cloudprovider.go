package cloudprovider

type CloudProvider interface {
	// Name returns name of the cloud provider.
	Name() string

	// NodeGroups returns all node groups configured for this cloud provider.
	NodeGroups() []NodeGroup

	// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
	// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
	Refresh() error
}

// NodeGroup contains configuration info and functions to control a set
// of nodes that have the same capacity and set of labels.
type NodeGroup interface {
	Id() string
}
