package inreachforwarder

import (
	"fmt"

	"github.com/kusold/inreachforwarder/internal/inreachparser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Long:  `Start the application with the specified flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		flagValue := viper.GetString("flagname")
		start(flagValue)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("url", "", "The Garmin InReach URL that was emailed to you")
	viper.BindPFlag("url", startCmd.Flags().Lookup("url"))
}

func start(url string) {
	fmt.Printf("Start command executed with flag: %s\n", url)

	inreachparser.SendMessageToInReach(url, "Hello, World!")
}
