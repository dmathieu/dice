package cmd

import (
	"log"
	"time"

	"github.com/dmathieu/dice/controllers"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var watchFrequency string

// loopCmd represents the loop command
var loopCmd = &cobra.Command{
	Use:   "loop",
	Short: "Continuously roll the instances",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		k8Client, cloudClient, err := buildClients(cloud)
		if err != nil {
			log.Fatal(err)
		}

		wf, err := parseWatchFrequency(watchFrequency)

		glog.Infof("Starting controllers")
		doneCh := make(chan struct{})
		flagger := controllers.NewOldNodesFlaggerController(k8Client, *wf)
		go flagger.Run(doneCh, 24*time.Hour)

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

	rootCmd.PersistentFlags().StringVarP(&watchFrequency, "watch-frequency", "w", "10d", "How frequently the watcher will look for nodes to destroy")
}
