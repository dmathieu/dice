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

func newDeleteNodeController(kClient kube.Interface, cClient cloudprovider.CloudProvider) *DeleteNodeController {
	i := informers.NewSharedInformerFactory(kClient, controller.NoResyncPeriodFunc())
	controller := NewDeleteNodeController(kClient, cClient, i.Core().V1().Pods())
	controller.podListerSynced = alwaysReady
	return controller
}

func TestDeleteNodeController(t *testing.T) {
	kClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kClient, cClient)

	doneCh := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(doneCh)
	}()

	controller.Run(doneCh)
}

func TestDeleteNodeControllerAddPod(t *testing.T) {
	pod := &corev1.Pod{}

	kClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()
	controller := newDeleteNodeController(kClient, cClient)
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
	flagged_node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "flagged-node",
			Labels: map[string]string{"dice": "roll"},
		},
		Spec: corev1.NodeSpec{
			Unschedulable: true,
		},
	}
	schedulable_node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "schedulable_node",
		},
	}
	node_with_pods := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node_with_pods",
			Labels: map[string]string{},
		},
	}
	pod_on_node := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod_on_node",
		},
		Spec: corev1.PodSpec{
			NodeName: node_with_pods.Name,
		},
	}

	t.Run("when pod is still running", func(t *testing.T) {
		kClient := fake.NewSimpleClientset()
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kClient, cClient)

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
		kClient := fake.NewSimpleClientset(node_with_pods, pod_on_node)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: node_with_pods.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable but node is not flagged", func(t *testing.T) {
		kClient := fake.NewSimpleClientset(node)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kClient, cClient)

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
		kClient := fake.NewSimpleClientset(schedulable_node)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: schedulable_node.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 0, len(cClient.DeletedNodes))
	})

	t.Run("when node has no other pods, is unschedulable and is flagged", func(t *testing.T) {
		kClient := fake.NewSimpleClientset(flagged_node)
		cClient := cloudtest.NewTestCloudProvider()
		controller := newDeleteNodeController(kClient, cClient)

		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-pod",
			},
			Spec: corev1.PodSpec{
				NodeName: flagged_node.Name,
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodSucceeded,
			},
		}
		controller.updatePod(pod, pod)
		assert.Equal(t, 1, len(cClient.DeletedNodes))
		assert.Equal(t, flagged_node.Name, cClient.DeletedNodes[0].Name)
	})
}
