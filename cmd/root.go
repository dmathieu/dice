package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	kubeConfig  string
	cloud       string
	concurrency int

	rootCmd = &cobra.Command{
		Use:   "dice",
		Short: "Zero-downtime rolling of a kubernetes' cluster's instances",
		Long:  `Replace all instances within a cluster with zero downtime.`,
	}
)

// Execute runs the content of the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultKubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if os.Getenv("KUBECONFIG") != "" {
		defaultKubeConfig = os.Getenv("KUBECONFIG")
	}

	rand.Seed(time.Now().Unix())

	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kube-config", "k", defaultKubeConfig, "Path to the kubernetes config, when running out of the cluster")
	rootCmd.PersistentFlags().StringVarP(&cloud, "cloud", "c", "", "Cloud Provider used by the cluster")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "p", 1, "Number of instances to roll concurrently")
}
