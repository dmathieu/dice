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

	podInformer     coreinformers.PodInformer
	podListerSynced cache.InformerSynced
}

func NewDeleteNodeController(kClient kube.Interface, cClient cloudprovider.CloudProvider, podInformer coreinformers.PodInformer) *DeleteNodeController {
	controller := &DeleteNodeController{
		kubeClient:  kClient,
		cloudClient: cClient,
		podInformer: podInformer,
	}

	controller.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addPod,
		UpdateFunc: controller.updatePod,
		DeleteFunc: controller.deletePod,
	})
	controller.podListerSynced = controller.podInformer.Informer().HasSynced

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

func (c *DeleteNodeController) handlePodChange(pod *corev1.Pod) {
	switch pod.Status.Phase {
	case corev1.PodSucceeded, corev1.PodFailed:
		// Continue
	default:
		return
	}

	pods, err := kubernetes.GetPods(c.kubeClient, kubernetes.PodNotTerminated(), kubernetes.PodOnNode(pod.Spec.NodeName))
	if err != nil {
		utilruntime.HandleError(err)
	}

	if len(pods) == 0 {
		node, err := kubernetes.FindNode(c.kubeClient, pod.Spec.NodeName)
		if err != nil {
			utilruntime.HandleError(err)
		}

		if node.Spec.Unschedulable && node.IsFlagged() {
			glog.Infof("Deleting node %s", node.Name)
			err := c.cloudClient.Delete(node)
			if err != nil {
				utilruntime.HandleError(err)
			}
		}
	}
}
