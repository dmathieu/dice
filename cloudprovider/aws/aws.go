package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/kubernetes"
)

// ProviderName is the name of this provider
const ProviderName = "aws"

// CloudProvider is a cloud provider allowing manipulating an AWS cluster
type CloudProvider struct {
	svc ec2iface.EC2API
}

// NewAWSCloudProvider instantiates a new CloudProvider
func NewAWSCloudProvider() cloudprovider.CloudProvider {
	svc := ec2.New(session.New())
	return &CloudProvider{svc}
}

// Name is the name of that cloud provider
func (t *CloudProvider) Name() string {
	return ProviderName
}

// Delete sends a delete request to the specified instance
func (t *CloudProvider) Delete(node *kubernetes.Node) error {
	instance, err := t.findInstanceFromNode(node)
	if err != nil {
		return err
	}

	_, err = t.svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance.InstanceId},
	})

	return err
}

func (t *CloudProvider) findInstanceFromNode(node *kubernetes.Node) (*ec2.Instance, error) {
	var addresses []*string
	for _, a := range node.Status.Addresses {
		addresses = append(addresses, &a.Address)
	}

	result, err := t.svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("private-dns-name"),
				Values: addresses,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	err = t.checkReservations(node, result.Reservations)
	if err != nil {
		return nil, err
	}
	return result.Reservations[0].Instances[0], nil
}

func (t *CloudProvider) checkReservations(node *kubernetes.Node, r []*ec2.Reservation) error {
	if len(r) == 0 {
		return fmt.Errorf("no reservations found matching node %s", node.Name)
	}
	if len(r) > 1 {
		return fmt.Errorf("found %d reservations matching node %s", len(r), node.Name)
	}
	if len(r[0].Instances) > 1 {
		return fmt.Errorf("found %d instances matching node %s", len(r[0].Instances), node.Name)
	}
	return nil
}
