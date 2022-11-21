package cmd

import (
	"fmt"

	"github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/spf13/cobra"

	"github.com/paramonies/ya-gophkeeper/internal/client/config"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

var (
	log *logger.Logger
	cli gophkeeper.GophkeeperServiceClient
	cfg *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "GophKeeper is a service to store and protect your important data",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

GophKeeper is a service, that gives you the possibilities to save you data and retrieve it from different devices. 
Service is synchronized between all you devices, where you are authenticated.
This application is a CLI tool to interact with the service.
Type -help to see more.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client ready")
	},
}

func Init(l *logger.Logger, c gophkeeper.GophkeeperServiceClient, cf *config.Config) error {
	log = l
	cli = c
	cfg = cf
	err := rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
