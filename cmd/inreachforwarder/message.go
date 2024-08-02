package inreachforwarder

import (
	"log"

	"github.com/kusold/inreachforwarder/internal/inreachparser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Messaging CLI",
	Long:  `Send or Receive one-off messages`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This is a strange viper/cobra interaction.
		// https://github.com/spf13/cobra/issues/875#issuecomment-536647568
		viper.BindPFlag("url", cmd.Flags().Lookup("url"))
		return nil
	},
}

var url string

func init() {
	rootCmd.AddCommand(messageCmd)
	messageCmd.AddCommand(messageSendCmd)
	messageCmd.AddCommand(messageReadCmd)
	messageCmd.PersistentFlags().String("url", "", "The Garmin InReach URL that was emailed to you")
	// messageCmd.PersistentFlags().StringVar(&url, "url", "", "The Garmin InReach URL that was emailed to you")
	// messageCmd.MarkPersistentFlagRequired("url")

	//		inreachparser.SendMessageToInReach(url, "Hello, World!")
	//	}
}

var messageSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message",
	Long:  `Send a message to your InReach device`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		// log.Printf("Sending message to %s\n", url)

		inreachparser.SendMessageToInReach(url, "Hello, World!")
	},
}

var messageReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a message",
	Long:  `Read a message sent by your InReach device`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("url")
		// log.Printf("Sending message to %s\n", url)

		msg, err := inreachparser.ReadMessageFromInReach(url)
		if err != nil {
			panic(err)
		}
		log.Print(msg)
	},
}
