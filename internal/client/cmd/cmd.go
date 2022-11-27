package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/paramonies/ya-gophkeeper/internal/client"
	"github.com/paramonies/ya-gophkeeper/internal/client/config"
	"github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
	"github.com/paramonies/ya-gophkeeper/pkg/logger"
)

var (
	log     *logger.Logger
	cfg     *config.Config
	cliUser gophkeeper.UserServiceClient
	cliPass gophkeeper.PasswordServiceClient
	cliText gophkeeper.TextServiceClient
	cliBin  gophkeeper.BinaryServiceClient
	cliCard gophkeeper.CardServiceClient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   commands[RootCommand].Use,
	Short: commands[RootCommand].Long,
	Long:  commands[RootCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client ready")
	},
}

func Init(l *logger.Logger, cf *config.Config, cliSet *client.ClientSet) error {
	log = l
	cfg = cf

	cliUser = cliSet.UserClient
	cliPass = cliSet.PwdClient
	cliText = cliSet.TextClient
	cliBin = cliSet.BinClient
	cliCard = cliSet.CardClient

	err := rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
