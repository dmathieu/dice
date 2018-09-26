package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
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

func (c *DeleteNodeController) Run(ctx context.Context) error {
	defer utilruntime.HandleCrash()

	doneCh := make(chan struct{})
	if !controller.WaitForCacheSync("dice", doneCh, c.podListerSynced) {
		return errors.New("couldn't wait for cache sync")
	}

	ctx, cancel := context.WithCancel(ctx)
	for {
		select {
		case <-doneCh:
			cancel()
			return nil
		case <-ctx.Done():
			cancel()
			return nil
		}
	}
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
			err := c.cloudClient.Delete(node)
			if err != nil {
				utilruntime.HandleError(err)
			}
		}
	}
}
