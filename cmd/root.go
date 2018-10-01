package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dice",
	Short: "Zero-downtime rolling of a kubernetes' cluster's instances",
	Long:  `Replace all instances within a cluster with zero downtime.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
