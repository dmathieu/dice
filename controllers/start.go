package controllers

import (
	"errors"
	"math/rand"
	"time"

	"github.com/dmathieu/dice/kubernetes"
	kube "k8s.io/client-go/kubernetes"
)

type StartController struct {
	kubeClient kube.Interface
}

func NewStartController(kClient kube.Interface) *StartController {
	return &StartController{kubeClient: kClient}
}

func (c *StartController) Run(concurrency int) error {
	rand.Seed(time.Now().Unix())
	err := c.flagNodes()
	if err != nil {
		return err
	}

	return c.evictNodes(concurrency)
}

func (c *StartController) flagNodes() error {
	nodes, err := kubernetes.GetNodes(c.kubeClient, kubernetes.NodeFlagged())
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return errors.New("found already flagged nodes. Looks like a roll process is already running")
	}

	nodes, err = kubernetes.GetNodes(c.kubeClient)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		err = kubernetes.FlagNode(c.kubeClient, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *StartController) evictNodes(concurrency int) error {
	return kubernetes.EvictNodes(c.kubeClient, concurrency)
}
