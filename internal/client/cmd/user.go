package cmd

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"
)

var (
	registerUser pb.RegisterUserRequest
	loginUser    pb.LoginUserRequest
)

func init() {
	rootCmd.AddCommand(registerUserCmd)
	registerUserCmd.Flags().StringVarP(&registerUser.ServiceLogin, "login", "l", "", "New user login value.")
	registerUserCmd.Flags().StringVarP(&registerUser.ServicePass, "password", "p", "", "New user password value.")
	registerUserCmd.MarkFlagRequired("login")
	registerUserCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringVarP(&loginUser.Login, "login", "l", "", "New user login value.")
	loginUserCmd.Flags().StringVarP(&loginUser.Password, "password", "p", "", "New user password value.")
	loginUserCmd.MarkFlagRequired("login")
	loginUserCmd.MarkFlagRequired("password")
}

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
		log.Info(fmt.Sprintf("user %s registered!", registerUser.GetServiceLogin()), "userID", response.GetUserID())
	},
}

// loginUserCmd represents the loginUser command
var loginUserCmd = &cobra.Command{
	Use:   "loginUser",
	Short: "Login user to the service",
	Long: `
This command login user.
Usage: gophkeeperclient loginUser --login=<login> --password=<password>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatal("failed to get current linux user", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		response, err := cli.LoginUser(ctx, &loginUser)
		if err != nil {
			log.Error("failed to login user", err)
			return
		}

		storage.Users[u.Username] = response.GetJwt()
		localObjectsd, ok := storage.Objects[u.Username]
		if !ok {
			localObjectsd = storage.CreateStorage()
			storage.Objects[u.Username] = localObjectsd
		}

		// after successful login - get JWT and send to server to synchronize data.
		_ = metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+response.GetJwt())

		//syncResp, err := c.SyncVault(ctxWTKN, &syncData)
		//if err != nil {
		//	log.Println(`[ERROR]:`, err)
		//	fmt.Println("request failed. please try again.")
		//	return
		//}
		//
		//fmt.Println("Get latest data from server: ", syncResp.GetStatus())
		//
		//fmt.Print("Synchronizing: ")
		//updVault := clserv.CombineVault(clstor.Local[u.Username], clserv.VaultSyncConvert(syncResp))
		////save actual data
		//clstor.Local[u.Username] = updVault

		log.Info(fmt.Sprintf("user %s loged in!", loginUser.GetLogin()), "userID", response.GetUserID())
	},
}
