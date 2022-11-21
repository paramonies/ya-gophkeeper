package cmd

import (
	"context"
	"os/user"
	"time"

	"github.com/spf13/cobra"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var (
	registerUser pb.RegisterUserRequest
)

// registerUserCmd represents the registerUser command
var registerUserCmd = &cobra.Command{
	Use:   "registerUser",
	Short: "Register new user in the service.",
	Long: `
This command register a new user.
Usage: gophkeeperclient registerUser --login=<login> --password=<password>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to get current linux user", err)
			return
		}

		log.Debug("register user", "user", u.Username, u.Name) //user vitretyakov Владимир Третьяков

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		response, err := cli.RegisterUser(ctx, &registerUser)
		if err != nil {
			log.Error("failed to register user", err)
			return
		}

		storage.Users[u.Username] = response.GetJwt()
		// init for the new user local storage
		storage.Objects[u.Username] = storage.CreateStorage()
		log.Debug("userID", response.GetUserID())
	},
}

func init() {
	rootCmd.AddCommand(registerUserCmd)
	registerUserCmd.Flags().StringVarP(&registerUser.ServiceLogin, "login", "l", "", "New user login value.")
	registerUserCmd.Flags().StringVarP(&registerUser.ServicePass, "password", "p", "", "New user password value.")
	registerUserCmd.MarkFlagRequired("login")
	registerUserCmd.MarkFlagRequired("password")
}
