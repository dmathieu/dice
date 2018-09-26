package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/dmathieu/dice/cloudprovider"
	corev1 "k8s.io/api/core/v1"
)

const ProviderName = "aws"

type AWSCloudProvider struct {
	svc ec2iface.EC2API
}

func NewAWSCloudProvider() cloudprovider.CloudProvider {
	svc := ec2.New(session.New())
	return &AWSCloudProvider{svc}
}

func (t *AWSCloudProvider) Name() string {
	return ProviderName
}

func (t *AWSCloudProvider) Delete(node *corev1.Node) error {
	instance, err := t.findInstanceFromNode(node)
	if err != nil {
		return err
	}

	_, err = t.svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance.InstanceId},
	})

	return err
}

func (t *AWSCloudProvider) findInstanceFromNode(node *corev1.Node) (*ec2.Instance, error) {
	var addresses []*string
	for _, a := range node.Status.Addresses {
		addresses = append(addresses, &a.Address)
	}

	result, err := t.svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("network-interface.private-dns-name"),
				Values: addresses,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Reservations) == 0 {
		return nil, fmt.Errorf("no instances found matching node %s", node.Name)
	}
	if len(result.Reservations) > 1 {
		return nil, fmt.Errorf("found %d reservations matching node %s", len(result.Reservations), node.Name)
	}
	if len(result.Reservations[0].Instances) > 1 {
		return nil, fmt.Errorf("found %d instances matching node %s", len(result.Reservations[0].Instances), node.Name)
	}

	return result.Reservations[0].Instances[0], nil
}

func (t *AWSCloudProvider) Refresh() error {
	return nil
}
