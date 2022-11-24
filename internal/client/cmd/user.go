package cmd

import (
	"context"
	"fmt"
	"google.golang.org/grpc/status"
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
	syncUserData pb.SyncUserDataResponse
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

// syncUserDataCmd represents the syncUserData command
var syncUserDataCmd = &cobra.Command{
	Use:   "syncUserData",
	Short: "Synchronize local user data with the server database",
	Long: `
This command provides latest data from the server database.
Then the database data, with version higher, that the local version, is saved to local storage.
During the saving of local data to the database, in case of version conflict(database version is higher/newer), you will be alerted by a warning.
Usage: gophkeeperclient syncUserData`,
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
		data, err := cliUser.SyncUserData(newCtx, &pb.SyncUserDataRequest{})
		if err != nil {
			st, _ := status.FromError(err)
			msg := fmt.Sprintf("request failed. statusCode: %v, message: %s", st.Code(), st.Message())
			log.Error(msg, err)
			return
		}

		fmt.Println("!!!")
		for _, p := range data.Passwords {
			fmt.Printf("%s %s %s %d \n", p.GetLogin(), p.GetPassword(), p.GetMeta(), p.GetVersion())
		}

		//// check server response
		//if response.GetStatus() != "success" {
		//	log.Println(response.GetStatus())
		//	fmt.Println("request failed. please try again.")
		//	return
		//}
		//
		////check for latest version data
		//syncVault := clserv.CombineVault(clstor.Local[u.Username], clserv.VaultSyncConvert(response))
		//
		//// save actual data
		//clstor.Local[u.Username] = syncVault
		//fmt.Println(response.GetStatus())
	},
}
