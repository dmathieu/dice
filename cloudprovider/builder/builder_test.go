package builder

import (
	"errors"
	"testing"

	"github.com/dmathieu/dice/cloudprovider/aws"
	"github.com/dmathieu/dice/cloudprovider/test"
	"github.com/stretchr/testify/assert"
)

func TestNewCloudProvider(t *testing.T) {
	t.Run("with an unknown cloud provider", func(t *testing.T) {
		c, err := NewCloudProvider("foobar")
		assert.Nil(t, c)
		assert.Equal(t, errors.New("Unknown cloud provider \"foobar\""), err)
	})

	t.Run("with the aws cloud provider", func(t *testing.T) {
		c, err := NewCloudProvider("aws")
		assert.Nil(t, err)
		assert.IsType(t, &aws.AWSCloudProvider{}, c)
	})

	t.Run("with the test cloud provider", func(t *testing.T) {
		c, err := NewCloudProvider("test")
		assert.Nil(t, err)
		assert.IsType(t, &test.TestCloudProvider{}, c)
	})
}
