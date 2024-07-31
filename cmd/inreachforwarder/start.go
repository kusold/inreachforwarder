package main

import (
	"fmt"

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
	startCmd.Flags().String("flagname", "default", "Description of the flag")
	viper.BindPFlag("flagname", startCmd.Flags().Lookup("flagname"))
}

func start(flagValue string) {
	fmt.Printf("Start command executed with flag: %s\n", flagValue)
}
