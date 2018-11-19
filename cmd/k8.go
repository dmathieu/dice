package cmd

import (
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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
