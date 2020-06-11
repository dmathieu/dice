package cmd

import (
	"log"

	"github.com/dmathieu/dice/controllers"
	"github.com/dmathieu/dice/kubernetes"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	watchFrequency string
	maxUptime      string
	dieAfter       string
)

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

		wf, err := parseStringDuration(watchFrequency)
		if err != nil {
			log.Fatalf("watch frequency: %q", err)
		}
		uptime, err := parseStringDuration(maxUptime)
		if err != nil {
			log.Fatalf("uptime: %q", err)
		}
		da, err := parseStringDuration(dieAfter)
		if err != nil {
			log.Fatalf("uptime: %q", err)
		}

		kubernetes.Setup(kubernetes.FlagValue("roll-loop"))
		glog.Infof("Starting controllers")

		doneCh := make(chan struct{})
		flagger := controllers.NewOldNodesFlaggerController(k8Client, *wf)
		go flagger.Run(doneCh, *uptime)
		defer close(doneCh)

		c, err := runWatchControllers(k8Client, cloudClient, concurrency, true)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		glog.Infof("Started all controllers")

		c.Run(da)
	},
}

func init() {
	rootCmd.AddCommand(loopCmd)

	rootCmd.PersistentFlags().StringVarP(&watchFrequency, "watch-frequency", "w", "5m", "How frequently the watcher will look for nodes to destroy")
	rootCmd.PersistentFlags().StringVarP(&maxUptime, "max-uptime", "u", "10d", "The uptime after which dice will kill an instance")
	rootCmd.PersistentFlags().StringVarP(&dieAfter, "die-after", "d", "11d", "The uptime after which dice will consider the watchers have crashed and will kill itself")
}
