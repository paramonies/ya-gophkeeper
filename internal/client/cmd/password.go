package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"

	"github.com/spf13/cobra"
)

// savePairCmd represents the savePair command
var savePairCmd = &cobra.Command{
	Use:   "createPassword",
	Short: "Create a new Password of login&password",
	Long: `
This command allows to the authenticated user to save new password data.
Usage: gophkeeperclient createPassword --login=<login_to_save> --password=<password_to_save> --meta=<meta_info_for_saved_login&password>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to get current linux user", err)
		}

		jwt, ok := storage.Users[u.Username]
		if !ok {
			log.Fatal("user not authenticated", err)
		}

		store, ok := storage.Objects[u.Username]
		if !ok {
			log.Fatal("user not found. Please register", nil)
		}
		password, ok := store.Password[createPassword.Login]
		if ok {
			createPassword.Version = password.Version + 1
		} else {
			createPassword.Version = 1
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

		res, err := cliPass.CreatePassword(newCtx, &createPassword)
		if err != nil {
			log.Error("failed to created password", err)
			return
		}

		store.Password[createPassword.Login] = &model.Password{
			Login:    createPassword.Login,
			Password: createPassword.Password,
			Meta:     createPassword.Meta,
			Version:  createPassword.Version,
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("create password for %s login!", createPassword.Login), "status", res.GetStatus())

	},
}

var (
	createPassword pb.CreatePasswordRequest
)

func init() {
	rootCmd.AddCommand(savePairCmd)
	savePairCmd.Flags().StringVarP(&createPassword.Login, "login", "l", "", "Login to save.")
	savePairCmd.Flags().StringVarP(&createPassword.Password, "password", "p", "", "Password to save.")
	savePairCmd.Flags().StringVarP(&createPassword.Meta, "meta", "m", "", "Meta info for the saved password. Optional.")
	savePairCmd.MarkFlagRequired("login")
	savePairCmd.MarkFlagRequired("password")
}
