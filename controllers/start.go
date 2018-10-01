package controllers

import (
	"errors"

	"github.com/dmathieu/dice/kubernetes"
	kube "k8s.io/client-go/kubernetes"
)

type StartController struct {
	kubeClient kube.Interface
}

func (p *StartController) Run() error {
	err := p.flagNodes()
	if err != nil {
		return err
	}

	return nil
}

func (p *StartController) flagNodes() error {
	nodes, err := kubernetes.GetNodes(p.kubeClient, kubernetes.NodeFlagged())
	if err != nil {
		return err
	}
	if len(nodes) > 0 {
		return errors.New("found already flagged nodes. Looks like a roll process is already running")
	}

	nodes, err = kubernetes.GetNodes(p.kubeClient)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		err = kubernetes.FlagNode(p.kubeClient, n)
		if err != nil {
			return err
		}
	}

	return nil
}
