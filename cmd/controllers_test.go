package cmd

import (
	"testing"
	"time"

	cloudtest "github.com/dmathieu/dice/cloudprovider/test"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestRunWatchControllers(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	cClient := cloudtest.NewTestCloudProvider()

	c, err := runWatchControllers(kubeClient, cClient)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	go func() {
		time.Sleep(time.Second)
		close(c.evictDoneCh)
	}()

	c.Run()
	c.Close()
}
