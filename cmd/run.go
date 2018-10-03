package cmd

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dmathieu/dice/cloudprovider/builder"
	"github.com/dmathieu/dice/controllers"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig  string
	cloud       string
	concurrency int
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the instances rolling",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		k8Client, err := getK8Client()
		if err != nil {
			log.Fatal(err)
		}

		cloudClient, err := builder.NewCloudProvider(cloud)
		if err != nil {
			log.Fatal(err)
		}

		glog.Infof("Starting controllers")

		doneCh := make(chan struct{})
		i := informers.NewSharedInformerFactory(k8Client, time.Second*30)
		evict := controllers.NewEvictNodeController(k8Client, i.Core().V1().Nodes(), concurrency)
		go evict.Run(doneCh)

		delete := controllers.NewDeleteNodeController(k8Client, cloudClient, i.Core().V1().Pods(), i.Core().V1().Nodes())
		go delete.Run(doneCh)

		i.Start(doneCh)

		start := controllers.NewStartController(k8Client)
		err = start.Run(concurrency)
		if err != nil {
			log.Fatal(err)
		}
		glog.Infof("Started all controllers")

		for {
			select {
			case <-doneCh:
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	defaultKubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if os.Getenv("KUBECONFIG") != "" {
		defaultKubeConfig = os.Getenv("KUBECONFIG")
	}

	runCmd.Flags().StringVarP(&kubeConfig, "kube-config", "k", defaultKubeConfig, "Path to the kubernetes config, when running out of the cluster")
	runCmd.Flags().StringVarP(&cloud, "cloud", "c", "", "Cloud Provider used by the cluster")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "p", 1, "Number of instances to roll concurrently")
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
