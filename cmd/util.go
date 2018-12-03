package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/dmathieu/dice/cloudprovider"
	"github.com/dmathieu/dice/cloudprovider/builder"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var wfRegex = regexp.MustCompile("^([0-9]+)(s|m|h|d)$")

func parseStringDuration(wf string) (*time.Duration, error) {
	s := wfRegex.FindStringSubmatch(wf)
	if len(s) != 3 {
		return nil, fmt.Errorf("invalid duration %q", wf)
	}
	v, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, err
	}

	var d time.Duration
	switch s[2] {
	case "s":
		d = time.Duration(v) * time.Second
	case "m":
		d = time.Duration(v) * time.Minute
	case "h":
		d = time.Duration(v) * time.Hour
		return &d, nil
	case "d":
		d = time.Duration(v) * 24 * time.Hour
	default:
		return nil, fmt.Errorf("couldn't find any matching frequencies for %q", wf)
	}

	return &d, nil
}

func k8Config() (*restclient.Config, error) {
	if len(os.Getenv("KUBERNETES_SERVICE_PORT")) > 0 && len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0 {
		return restclient.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfig)
}

func getK8Client() (*kubernetes.Clientset, error) {
	k8Config, err := k8Config()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(k8Config)
}

func buildClients(cloud string) (*kubernetes.Clientset, cloudprovider.CloudProvider, error) {
	k8Client, err := getK8Client()
	if err != nil {
		return nil, nil, err
	}

	cloudClient, err := builder.NewCloudProvider(cloud)
	if err != nil {
		return nil, nil, err
	}

	return k8Client, cloudClient, nil
}
