package controllers

import (
	"fmt"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
)

type DeleteNodeController struct {
	kubeClient  kube.Interface
	cloudClient cloudprovider.CloudProvider

	podInformer      coreinformers.PodInformer
	podListerSynced  cache.InformerSynced
	nodeInformer     coreinformers.NodeInformer
	nodeListerSynced cache.InformerSynced
}

func NewDeleteNodeController(kClient kube.Interface, cClient cloudprovider.CloudProvider, podInformer coreinformers.PodInformer, nodeInformer coreinformers.NodeInformer) *DeleteNodeController {
	controller := &DeleteNodeController{
		kubeClient:   kClient,
		cloudClient:  cClient,
		podInformer:  podInformer,
		nodeInformer: nodeInformer,
	}

	controller.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addPod,
		UpdateFunc: controller.updatePod,
		DeleteFunc: controller.deletePod,
	})
	controller.podListerSynced = controller.podInformer.Informer().HasSynced

	controller.nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addNode,
		UpdateFunc: controller.updateNode,
		DeleteFunc: controller.deleteNode,
	})
	controller.nodeListerSynced = controller.nodeInformer.Informer().HasSynced

	return controller
}

func (c *DeleteNodeController) Run(doneCh chan struct{}) {
	defer utilruntime.HandleCrash()
	if !controller.WaitForCacheSync("delete node", doneCh, c.podListerSynced) {
		return
	}
	<-doneCh
}

func (c *DeleteNodeController) addPod(obj interface{}) {
	// We have nothing to handle on add
}

func (c *DeleteNodeController) updatePod(old, cur interface{}) {
	pod, ok := cur.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get pod %#v", cur))
		return
	}
	c.handlePodChange(pod)
}

func (c *DeleteNodeController) deletePod(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		utilruntime.HandleError(fmt.Errorf("Couldn't get pod %#v", obj))
		return
	}
	c.handlePodChange(pod)
}

func (c *DeleteNodeController) addNode(obj interface{}) {
	// We have nothing to handle on add
}

func (c *DeleteNodeController) updateNode(old, cur interface{}) {
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

	if oldNode.Spec.Unschedulable {
		// Node was already non-ready before. We don't need to delete again.
		return
	}

	err := c.handleNodeDeletion(&kubernetes.Node{Node: node})
	if err != nil {
		utilruntime.HandleError(err)
	}
}

func (c *DeleteNodeController) deleteNode(obj interface{}) {
	// We have nothing to handle on delete
}

func (c *DeleteNodeController) handlePodChange(pod *corev1.Pod) {
	switch pod.Status.Phase {
	case corev1.PodSucceeded, corev1.PodFailed:
		// Continue
	default:
		return
	}

	node, err := kubernetes.FindNode(c.kubeClient, pod.Spec.NodeName)
	if err != nil {
		utilruntime.HandleError(err)
	}
	err = c.handleNodeDeletion(node)
	if err != nil {
		utilruntime.HandleError(err)
	}
}

func (c *DeleteNodeController) handleNodeDeletion(node *kubernetes.Node) error {
	if !node.Spec.Unschedulable || !node.IsFlagged() {
		return nil
	}

	pods, err := kubernetes.GetPods(c.kubeClient, kubernetes.PodNotTerminated(), kubernetes.PodOnNode(node.Name))
	if err != nil {
		return err
	}

	if len(pods) == 0 {
		glog.Infof("Deleting node %s", node.Name)
		err := c.cloudClient.Delete(node)
		if err != nil {
			return err
		}
	}
	return nil
}
