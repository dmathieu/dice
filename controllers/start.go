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
	nodes, err := kubernetes.GetNodes(c.kubeClient) //, kubernetes.NodeFlagged())
	if err != nil {
		return err
	}

	if concurrency > len(nodes) {
		concurrency = len(nodes)
	}

	evicted := map[string]*kubernetes.Node{}

	for len(evicted) < concurrency {
		eNode := nodes[rand.Intn(len(nodes))]
		if evicted[eNode.Name] != nil {
			continue
		}

		err = kubernetes.EvictNode(c.kubeClient, eNode)
		if err != nil {
			return err
		}
		evicted[eNode.Name] = eNode
	}

	return nil
}
