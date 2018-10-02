package controllers

import (
	"testing"
	"time"

	"github.com/dmathieu/dice/cloudprovider"
	cloudtest "github.com/dmathieu/dice/cloudprovider/test"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/controller"
)

func newDeleteNodeController(kubeClient kube.Interface, cloudClient cloudprovider.CloudProvider) *DeleteNodeController {
	i := informers.NewSharedInformerFactory(kubeClient, controller.NoResyncPeriodFunc())
	controller := NewDeleteNodeController(kubeClient, cloudClient, i.Core().V1().Pods(), i.Core().V1().Nodes())
	controller.podListerSynced = alwaysReady
	controller.nodeListerSynced = alwaysReady
	return controller
}

func TestDeleteNodeController(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kubeClient, cClient)

	doneCh := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(doneCh)
	}()

	controller.Run(doneCh)
}

func TestDeleteNodeControllerAddPod(t *testing.T) {
	pod := &corev1.Pod{}

	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kubeClient, cClient)
	controller.addPod(pod)
}

func TestDeleteNodeControllerUpdatePod(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "my-node",
			Labels: map[string]string{},
		},
		Spec: corev1.NodeSpec{
			Unschedulable: true,
		},
	}
	flaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "flagged-node",
			Labels: map[string]string{"dice": "roll"},
		},
		Spec: corev1.NodeSpec{
			Unschedulable: true,
		},
	}
	schedulableNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "schedulable_node",
		},
	}
	nodeWithPods := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node_with_pods",
			Labels: map[string]string{},
		},
	}
	podOnNode := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod_on_node",
		},
		Spec: corev1.PodSpec{
			NodeName: nodeWithPods.Name,
		},
	}

	t.Run("when pod is still running", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset()
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kubeClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has other pods", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset(nodeWithPods, podOnNode)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kubeClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: nodeWithPods.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable but node is not flagged", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset(node)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kubeClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: node.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods and is schedulable", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset(schedulableNode)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kubeClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: schedulableNode.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable and is flagged", func(t *testing.T) {
		kubeClient := fake.NewSimpleClientset(flaggedNode)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kubeClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: flaggedNode.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 1, len(cClient.DeletedNodes))
		assert.Equal(t, flaggedNode.Name, cClient.DeletedNodes[0].Name)
	})
}

func TestDeleteNodeControllerAddNode(t *testing.T) {
	node := &corev1.Node{}

	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kubeClient, cClient)
	controller.addNode(node)
}

func TestDeleteNodeControllerDeleteNode(t *testing.T) {
	node := &corev1.Node{}

	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kubeClient, cClient)
	controller.deleteNode(node)
}

func TestDeleteNodeControllerUpdateNode(t *testing.T) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "my-node",
			Labels: map[string]string{},
		},
		Spec: corev1.NodeSpec{
			Unschedulable: true,
		},
	}
	flaggedNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "flagged-node",
			Labels: map[string]string{"dice": "roll"},
		},
		Spec: corev1.NodeSpec{
			Unschedulable: true,
		},
	}
	schedulableNode := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "schedulable_node",
		},
	}
	nodeWithPods := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node_with_pods",
			Labels: map[string]string{},
		},
	}

	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kubeClient, cClient)

	t.Run("when still has running pods", func(t *testing.T) {

		controller.updateNode(nodeWithPods, nodeWithPods)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable but is not flagged", func(t *testing.T) {
		controller.updateNode(node, node)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods and is schedulable", func(t *testing.T) {
		controller.updateNode(schedulableNode, schedulableNode)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node was already unschedulable before", func(t *testing.T) {
		controller.updateNode(schedulableNode, schedulableNode)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable and is flagged", func(t *testing.T) {
		controller.updateNode(schedulableNode, flaggedNode)
		assert.Equal(t, 1, len(cClient.DeletedNodes))
		assert.Equal(t, flaggedNode.Name, cClient.DeletedNodes[0].Name)
	})
}
