package cmd

import (
	"log"
	"time"

	"github.com/dmathieu/dice/cloudprovider/builder"
	"github.com/dmathieu/dice/controllers"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// loopCmd represents the loop command
var loopCmd = &cobra.Command{
	Use:   "loop",
	Short: "Continuously roll the instances",
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
		flagger := controllers.NewOldNodesFlaggerController(k8Client, 5*time.Minute)
		flagger.Run(doneCh, 24*time.Hour)

		c, err := runWatchControllers(k8Client, cloudClient, concurrency, true)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		glog.Infof("Started all controllers")

		c.Run()
		close(doneCh)
	},
}

func init() {
	rootCmd.AddCommand(loopCmd)
}
