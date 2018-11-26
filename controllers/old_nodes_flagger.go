package controllers

import (
	"time"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kube "k8s.io/client-go/kubernetes"
)

// OldNodesFlaggerController is a controller that flags old nodes to be drained
// and watches when nodes reach an uptime too big
type OldNodesFlaggerController struct {
	kubeClient kube.Interface
	interval   time.Duration
}

// NewOldNodesFlaggerController creates a new controller which flags old nodes for replacement
func NewOldNodesFlaggerController(client kube.Interface, i time.Duration) *OldNodesFlaggerController {
	return &OldNodesFlaggerController{kubeClient: client, interval: i}
}

// Run executes the actions from the controller
func (c *OldNodesFlaggerController) Run(doneCh chan struct{}, maxUptime time.Duration) {
	ticker := time.NewTicker(c.interval)

	for {
		select {
		case <-ticker.C:
			nodes, err := kubernetes.GetNodes(c.kubeClient, kubernetes.NodeNotFlagged())
			if err != nil {
				utilruntime.HandleError(err)
			}

			t := metav1.NewTime(time.Now().Add(0 - maxUptime))
			for _, n := range nodes {
				if n.ObjectMeta.CreationTimestamp.Before(&t) {
					glog.Infof("Flagging node %q", n.ObjectMeta.Name)
					err = kubernetes.FlagNode(c.kubeClient, n)
					if err != nil {
						utilruntime.HandleError(err)
					}
				}

			}
		case <-doneCh:
			return
		}
	}
}
