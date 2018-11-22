package cmd

import (
	"log"

	"github.com/dmathieu/dice/controllers"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the instances rolling",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		k8Client, cloudClient, err := buildClients(cloud)
		if err != nil {
			log.Fatal(err)
		}

		glog.Infof("Starting controllers")

		flagger := controllers.NewAllNodesFlaggerController(k8Client)
		err = flagger.Run()
		if err != nil {
			log.Fatal(err)
		}

		c, err := runWatchControllers(k8Client, cloudClient, concurrency, false)
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		glog.Infof("Started all controllers")

		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
