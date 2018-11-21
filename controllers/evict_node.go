package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
)

// EvictNodeController is a controller which performs node eviction.
// It listsns on node events.
//
// When a new node comes online and ready to accept pods, it triggers an
// eviction for another node found randomly.
type EvictNodeController struct {
	kubeClient  kube.Interface
	concurrency int
	infinite    bool
	doneCh      chan struct{}

	nodeInformer     coreinformers.NodeInformer
	nodeListerSynced cache.InformerSynced
}

// NewEvictNodeController instantiates a new eviction controller
func NewEvictNodeController(client kube.Interface, nodeInformer coreinformers.NodeInformer, c int, i bool) *EvictNodeController {
	rand.Seed(time.Now().Unix())
	controller := &EvictNodeController{
		kubeClient:   client,
		concurrency:  c,
		infinite:     i,
		nodeInformer: nodeInformer,
	}

	controller.nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addNode,
		UpdateFunc: controller.updateNode,
		DeleteFunc: controller.deleteNode,
	})
	controller.nodeListerSynced = controller.nodeInformer.Informer().HasSynced

	return controller
}

// Run starts the controller
func (c *EvictNodeController) Run(doneCh chan struct{}) {
	defer utilruntime.HandleCrash()
	if !controller.WaitForCacheSync("evict node", doneCh, c.nodeListerSynced) {
		return
	}

	err := kubernetes.EvictNodes(c.kubeClient, c.concurrency)
	if err != nil {
		utilruntime.HandleError(err)
	}

	c.doneCh = make(chan struct{})

	for {
		select {
		case <-c.doneCh:
			close(doneCh)
			return
		case <-doneCh:
			close(c.doneCh)
			return
		}
	}
}

func (c *EvictNodeController) addNode(obj interface{}) {
	node, ok := obj.(*corev1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get node %#v", obj))
		return
	}
	c.handleNodeChange(node)
}

func (c *EvictNodeController) updateNode(old, cur interface{}) {
	oldNode, ok := old.(*corev1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get old node %#v", cur))
		return
	}
	node, ok := cur.(*corev1.Node)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get node %#v", cur))
		return
	}

	if (&kubernetes.Node{Node: oldNode}).IsReady() {
		// Node was already ready before. We don't need to evict another one.
		return
	}

	c.handleNodeChange(node)
}

func (c *EvictNodeController) deleteNode(obj interface{}) {
	// We have nothing to handle on delete
}

func (c *EvictNodeController) handleNodeChange(n *corev1.Node) {
	node := &kubernetes.Node{Node: n}

	if !node.IsReady() || node.IsFlagged() {
		return
	}

	nodes, err := kubernetes.GetNodes(c.kubeClient, kubernetes.NodeFlagged())
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	if len(nodes) == 0 && c.doneCh != nil && !c.infinite {
		glog.Infof("My job here is done!")
		close(c.doneCh)
		return
	}

	err = kubernetes.EvictNodes(c.kubeClient, c.concurrency)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
}
