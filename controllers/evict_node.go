package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
)

type EvictNodeController struct {
	kubeClient  kube.Interface
	cloudClient cloudprovider.CloudProvider

	nodeInformer     coreinformers.NodeInformer
	nodeListerSynced cache.InformerSynced
}

func NewEvictNodeController(kClient kube.Interface, nodeInformer coreinformers.NodeInformer) *EvictNodeController {
	rand.Seed(time.Now().Unix())
	controller := &EvictNodeController{
		kubeClient:   kClient,
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

func (c *EvictNodeController) Run(doneCh chan struct{}) {
	defer utilruntime.HandleCrash()
	if !controller.WaitForCacheSync("dice", doneCh, c.nodeListerSynced) {
		return
	}
	<-doneCh
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
	eNode := nodes[rand.Intn(len(nodes))]
	err = kubernetes.EvictNode(c.kubeClient, eNode)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
}
