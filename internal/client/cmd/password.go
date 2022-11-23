package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/metadata"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"

	"github.com/spf13/cobra"
)

var (
	createPassword pb.CreatePasswordRequest
	getPassword    pb.GetPasswordRequest
	deletePassword pb.DeletePasswordRequest
)

func init() {
	rootCmd.AddCommand(savePairCmd)
	savePairCmd.Flags().StringVarP(&createPassword.Login, "login", "l", "", "Login to save.")
	savePairCmd.Flags().StringVarP(&createPassword.Password, "password", "p", "", "Password to save.")
	savePairCmd.Flags().StringVarP(&createPassword.Meta, "meta", "m", "", "Meta info for the saved password. Optional.")
	savePairCmd.MarkFlagRequired("login")
	savePairCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(getPasswordCmd)
	getPasswordCmd.Flags().StringVarP(&getPassword.Login, "login", "l", "", "Login for password to search for.")
	getPasswordCmd.MarkFlagRequired("login")

	rootCmd.AddCommand(deletePasswordCmd)
	deletePasswordCmd.Flags().StringVarP(&deletePassword.Login, "login", "l", "", "Login for password to delete.")
	deletePasswordCmd.MarkFlagRequired("login")
}

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
		return
	},
}

// getPasswordCmd represents the getPassword command
var getPasswordCmd = &cobra.Command{
	Use:   "getPassword",
	Short: "Get a password data by login",
	Long: `
This command returns to the authenticated user the password data requested by login.
Usage: gophkeeperclient getPassword --login=<login>.`,
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

		pwd, ok := store.Password[getPassword.Login]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Local version for password data: login: %s, password: %s, meta: %s, version: %d. Make sure you have the latest version by synchronizing local storage",
				pwd.Login, pwd.Password, pwd.Meta, pwd.Version)
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliPass.GetPassword(newCtx, &pb.GetPasswordRequest{
			Login: getPassword.Login,
		})
		if err != nil {
			st, _ := status.FromError(err)
			msg := fmt.Sprintf("statusCode: %v, message: %s", st.Code(), st.Message())
			log.Error(msg, err)
			return
		}

		store.Password[getPassword.Login] = &model.Password{
			Login:    res.Login,
			Password: res.Password,
			Meta:     res.Meta,
			Version:  res.Version,
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		msg := fmt.Sprintf("Server version for password data: login: %s, password: %s, meta: %s, version: %d. Make sure you have the latest version by synchronizing local storage",
			res.Login, res.Password, res.Meta, res.Version)
		log.Info(msg)
		return
	},
}

// deletePasswordCmd represents the deletePassword command
var deletePasswordCmd = &cobra.Command{
	Use:   "deletePassword",
	Short: "Delete the password data for login",
	Long: `
This command allows to the authenticated user to delete the password data.
Usage: gophkeeperclient deletePassword --login=<login>.`,
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

		_, ok = store.Password[deletePassword.Login]
		// local version doesn't exist: nothing to delete
		if !ok {
			msg := fmt.Sprintf("Nothing found for login: %s.Make sure you have the latest version by synchronizing your local storage.",
				deletePassword.GetLogin())
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		_, err = cliPass.DeletePassword(newCtx, &pb.DeletePasswordRequest{
			Login: deletePassword.Login,
		})
		if err != nil {
			st, _ := status.FromError(err)
			msg := fmt.Sprintf("request failed. statusCode: %v, message: %s", st.Code(), st.Message())
			log.Error(msg, err)
			return
		}

		// successful response
		// delete local version
		delete(storage.Objects[u.Username].Password, deletePassword.GetLogin())

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		msg := fmt.Sprintf("Password data for login %s deleted", deletePassword.GetLogin())
		log.Info(msg)
		return
	},
}
