package controllers

import (
	"math/rand"
	"time"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	kube "k8s.io/client-go/kubernetes"
)

// StartController is a controller that does all the bootstrap actions.
// Those include flagging all nodes as needing to be restarted, and triggering the first eviction.
type StartController struct {
	kubeClient kube.Interface
}

// NewStartController creates a new Start Controller
func NewStartController(client kube.Interface) *StartController {
	return &StartController{kubeClient: client}
}

// Run executes the actions from StartController
func (c *StartController) Run(concurrency int) error {
	rand.Seed(time.Now().Unix())
	return c.flagNodes()
}

func (c *StartController) flagNodes() error {
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
