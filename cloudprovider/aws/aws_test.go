package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMatchInterfaces(t *testing.T) {
	assert.Implements(t, (*cloudprovider.CloudProvider)(nil), &AWSCloudProvider{})
}

func TestName(t *testing.T) {
	p := &AWSCloudProvider{&mockEC2Client{}}
	assert.Equal(t, "aws", p.Name())
}

func TestDelete(t *testing.T) {
	node := &kubernetes.Node{
		Node: &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-node",
			},
		},
	}
	client := &mockEC2Client{
		reservations: []*ec2.Reservation{
			&ec2.Reservation{
				Instances: []*ec2.Instance{&ec2.Instance{}},
			},
		},
	}
	p := &AWSCloudProvider{client}
	assert.Nil(t, p.Delete(node))
}

func TestInstanceFromNode(t *testing.T) {
	node := &kubernetes.Node{
		Node: &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-node",
			},
		},
	}

	t.Run("when there is no node", func(t *testing.T) {
		p := &AWSCloudProvider{&mockEC2Client{}}
		i, err := p.findInstanceFromNode(node)
		assert.Nil(t, i)
		assert.Equal(t, errors.New("no instances found matching node my-node"), err)
	})

	t.Run("when there are several reservations", func(t *testing.T) {
		client := &mockEC2Client{
			reservations: []*ec2.Reservation{
				&ec2.Reservation{},
				&ec2.Reservation{},
			},
		}
		p := &AWSCloudProvider{client}
		i, err := p.findInstanceFromNode(node)
		assert.Nil(t, i)
		assert.Equal(t, errors.New("found 2 reservations matching node my-node"), err)
	})

	t.Run("when there are several instances", func(t *testing.T) {
		client := &mockEC2Client{
			reservations: []*ec2.Reservation{
				&ec2.Reservation{
					Instances: []*ec2.Instance{
						&ec2.Instance{},
						&ec2.Instance{},
					},
				},
			},
		}
		p := &AWSCloudProvider{client}
		i, err := p.findInstanceFromNode(node)
		assert.Nil(t, i)
		assert.Equal(t, errors.New("found 2 instances matching node my-node"), err)
	})

	t.Run("fetch the instance", func(t *testing.T) {
		instance := &ec2.Instance{}
		client := &mockEC2Client{
			reservations: []*ec2.Reservation{
				&ec2.Reservation{
					Instances: []*ec2.Instance{instance},
				},
			},
		}
		p := &AWSCloudProvider{client}
		i, err := p.findInstanceFromNode(node)
		assert.Nil(t, err)
		assert.Equal(t, instance, i)
	})
}

func TestRefresh(t *testing.T) {
	p := &AWSCloudProvider{&mockEC2Client{}}
	assert.Nil(t, p.Refresh())
}
