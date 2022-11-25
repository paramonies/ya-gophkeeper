package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var (
	registerUser pb.RegisterUserRequest
	loginUser    pb.LoginUserRequest
	syncUserData pb.GetAllUserDataFromDBResponse
)

func init() {
	rootCmd.AddCommand(registerUserCmd)
	registerUserCmd.Flags().StringVarP(&registerUser.Login, "login", "l", "", "New user login value.")
	registerUserCmd.Flags().StringVarP(&registerUser.Password, "password", "p", "", "New user password value.")
	registerUserCmd.MarkFlagRequired("login")
	registerUserCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringVarP(&loginUser.Login, "login", "l", "", "New user login value.")
	loginUserCmd.Flags().StringVarP(&loginUser.Password, "password", "p", "", "New user password value.")
	loginUserCmd.MarkFlagRequired("login")
	loginUserCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(syncUserDataCmd)
}

// registerUserCmd represents the registerUser command
var registerUserCmd = &cobra.Command{
	Use:   commands[RegisterUserCommand].Use,
	Short: commands[RegisterUserCommand].Short,
	Long:  commands[RegisterUserCommand].Long,
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

		response, err := cliUser.RegisterUser(ctx, &registerUser)
		if err != nil {
			log.Error("failed to register user", err)
			return
		}

		storage.Users[u.Username] = response.GetJwt()
		// init for the new user local storage
		storage.Objects[u.Username] = storage.CreateStorage()
		log.Info(fmt.Sprintf("user %s registered!", registerUser.GetLogin()), "userID", response.GetUserID())
	},
}

// loginUserCmd represents the loginUser command
var loginUserCmd = &cobra.Command{
	Use:   commands[LoginUserCommand].Use,
	Short: commands[LoginUserCommand].Short,
	Long:  commands[LoginUserCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to get current linux user", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		response, err := cliUser.LoginUser(ctx, &loginUser)
		if err != nil {
			log.Error("failed to login user", err)
			return
		}

		storage.Users[u.Username] = response.GetJwt()
		localObjects, ok := storage.Objects[u.Username]
		if !ok {
			localObjects = storage.CreateStorage()
			storage.Objects[u.Username] = localObjects
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		// after successful login - get JWT and send to server to synchronize data.
		newCtx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+response.GetJwt())

		// send data to server and receive JWT in case of success. then save it in Users
		dataDB, err := cliUser.GetAllUserDataFromDB(newCtx, &pb.GetAllUserDataFromDBRequest{})
		if err != nil {
			st, _ := status.FromError(err)
			msg := fmt.Sprintf("request failed. statusCode: %v, message: %s", st.Code(), st.Message())
			log.Error(msg, err)
			return
		}

		fmt.Println("!!!")
		for _, p := range dataDB.Passwords {
			fmt.Printf("%s %s %s %d \n", p.GetLogin(), p.GetPassword(), p.GetMeta(), p.GetVersion())
		}

		for _, p := range storage.Objects[u.Username].Password {
			info := fmt.Sprintf("!!!login: %s, password: %s, meta: %s, version: %d", p.Login, p.Password, p.Meta, p.Version)
			fmt.Println(info)
		}

		// check for latest version data
		lastVerData := storage.SyncData(storage.Objects[u.Username], model.ProtoToLocalStorage(dataDB))
		//
		// save actual data
		storage.Objects[u.Username] = lastVerData

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info("user data synchronized!")
		for _, p := range storage.Objects[u.Username].Password {
			info := fmt.Sprintf("login: %s, password: %s, meta: %s, version: %d", p.Login, p.Password, p.Meta, p.Version)
			fmt.Println(info)
		}

		log.Info(fmt.Sprintf("user %s loged in!", loginUser.GetLogin()), "userID", response.GetUserID())
		return

	},
}

// syncUserDataCmd represents the syncUserData command
var syncUserDataCmd = &cobra.Command{
	Use:   commands[SyncDataCommand].Use,
	Short: commands[SyncDataCommand].Short,
	Long:  commands[SyncDataCommand].Long,
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

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		dataDB, err := cliUser.GetAllUserDataFromDB(newCtx, &pb.GetAllUserDataFromDBRequest{})
		if err != nil {
			st, _ := status.FromError(err)
			msg := fmt.Sprintf("request failed. statusCode: %v, message: %s", st.Code(), st.Message())
			log.Error(msg, err)
			return
		}

		fmt.Println("!!!")
		for _, p := range dataDB.Passwords {
			fmt.Printf("%s %s %s %d \n", p.GetLogin(), p.GetPassword(), p.GetMeta(), p.GetVersion())
		}

		for _, p := range storage.Objects[u.Username].Password {
			info := fmt.Sprintf("!!!login: %s, password: %s, meta: %s, version: %d", p.Login, p.Password, p.Meta, p.Version)
			fmt.Println(info)
		}

		// check for latest version data
		lastVerData := storage.SyncData(storage.Objects[u.Username], model.ProtoToLocalStorage(dataDB))
		//
		// save actual data
		storage.Objects[u.Username] = lastVerData

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info("user data synchronized!")
		for _, p := range storage.Objects[u.Username].Password {
			info := fmt.Sprintf("login: %s, password: %s, meta: %s, version: %d", p.Login, p.Password, p.Meta, p.Version)
			fmt.Println(info)
		}

		return
	},
}
