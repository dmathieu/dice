package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseStringDuration(t *testing.T) {
	t.Run("with an invalid value", func(t *testing.T) {
		f, err := parseStringDuration("foobar")
		assert.Nil(t, f)
		assert.Equal(t, errors.New("invalid duration \"foobar\""), err)
	})

	t.Run("with a value of seconds", func(t *testing.T) {
		f, err := parseStringDuration("27s")
		assert.Nil(t, err)
		assert.Equal(t, 27*time.Second, *f)
	})

	t.Run("with a value of minutes", func(t *testing.T) {
		f, err := parseStringDuration("32m")
		assert.Nil(t, err)
		assert.Equal(t, 32*time.Minute, *f)
	})

	t.Run("with a value of hours", func(t *testing.T) {
		f, err := parseStringDuration("164h")
		assert.Nil(t, err)
		assert.Equal(t, 164*time.Hour, *f)
	})

	t.Run("with a value of days", func(t *testing.T) {
		f, err := parseStringDuration("65d")
		assert.Nil(t, err)
		assert.Equal(t, 65*24*time.Hour, *f)
	})
}

func TestBuildClients(t *testing.T) {
	t.Run("with a kube config not found", func(t *testing.T) {
		tmpfile, err := generateKubeConfig()
		assert.Nil(t, err)
		os.Remove(tmpfile.Name())
		kubeConfig = tmpfile.Name()

		k8Client, cloudClient, err := buildClients("test")
		assert.NotNil(t, err)
		assert.Regexp(t, "no such file or directory", err.Error())
		assert.Nil(t, k8Client)
		assert.Nil(t, cloudClient)
	})

	t.Run("with an invalid cloud client", func(t *testing.T) {
		tmpfile, err := generateKubeConfig()
		assert.Nil(t, err)
		defer os.Remove(tmpfile.Name())
		kubeConfig = tmpfile.Name()

		k8Client, cloudClient, err := buildClients("foobar")
		assert.Equal(t, errors.New("Unknown cloud provider \"foobar\""), err)
		assert.Nil(t, k8Client)
		assert.Nil(t, cloudClient)
	})

	t.Run("with a valid client", func(t *testing.T) {
		tmpfile, err := generateKubeConfig()
		assert.Nil(t, err)
		defer os.Remove(tmpfile.Name())
		kubeConfig = tmpfile.Name()

		k8Client, cloudClient, err := buildClients("test")
		assert.Nil(t, err)

		if err != nil {
			fmt.Fprintf(os.Stdout, "%q\n", err.Error())
		}

		assert.NotNil(t, k8Client)
		assert.NotNil(t, cloudClient)
	})
}

func generateKubeConfig() (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		return nil, err
	}

	_, err = tmpfile.Write([]byte(`{
"apiVersion": "v1",
"kind": "Config",
"current-context": "default",
"clusters": [
	{
		"cluster": {"server": "foobar"},
		"name": "kubernetes"
	}
],
"contexts": [
	{
		"context": {"cluster": "kubernetes"},
		"name": "default"
	}
]
}`))
	if err != nil {
		return nil, err
	}

	return tmpfile, tmpfile.Close()
}
