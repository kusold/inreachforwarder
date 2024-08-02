package inreachforwarder

import (
	"github.com/kusold/inreachforwarder/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a server that monitors for pagerduty alerts",
	Long:  `Start a server that sends pagerduty alerts to your Garmin InReach device`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("api-token", "", "Pager Duty API User Token. Go to Profile -> My Profile -> User Settings -> API Access -> Create New API User Token")
	viper.BindPFlag("pagerduty.user-api-token", serverCmd.Flags().Lookup("api-token"))
}
