package controllers

import (
	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	kube "k8s.io/client-go/kubernetes"
)

// AllNodesFlaggerController is a controller that flags all nodes to be drained.
type AllNodesFlaggerController struct {
	kubeClient kube.Interface
}

// NewAllNodesFlaggerController creates a new controller which flags all nodes for replacement
func NewAllNodesFlaggerController(client kube.Interface) *AllNodesFlaggerController {
	return &AllNodesFlaggerController{kubeClient: client}
}

// Run executes the actions from the controller
func (c *AllNodesFlaggerController) Run() error {
	nodes, err := kubernetes.GetNodes(c.kubeClient)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		if n.IsFlagged() {
			glog.Infof("Found flagged nodes. Continuing with them.")
			return nil
		}
	}

	for _, n := range nodes {
		err = kubernetes.FlagNode(c.kubeClient, n)
		if err != nil {
			return err
		}
	}

	return nil
}
