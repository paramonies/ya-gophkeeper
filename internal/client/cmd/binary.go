package cmd

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/paramonies/ya-gophkeeper/internal/client/storage"
	"github.com/paramonies/ya-gophkeeper/internal/model"
	pb "github.com/paramonies/ya-gophkeeper/pkg/gen/api/gophkeeper/v1"

	"github.com/spf13/cobra"
)

var (
	createBinary pb.CreateBinaryRequest
	getBinary    pb.GetBinaryRequest
	deleteBinary pb.DeleteBinaryRequest
)

func init() {
	rootCmd.AddCommand(createBinaryCmd)
	createBinaryCmd.Flags().StringVarP(&createBinary.Title, "title", "t", "", "Title to save.")
	createBinaryCmd.Flags().StringVarP(&createBinary.Data, "data", "d", "",
		"Binary data to save.")
	createBinaryCmd.Flags().StringVarP(&createBinary.Meta, "meta", "m", "",
		"Meta info for the saved binary. Optional.")
	createBinaryCmd.MarkFlagRequired("title")
	createBinaryCmd.MarkFlagRequired("data")

	rootCmd.AddCommand(getBinaryCmd)
	getBinaryCmd.Flags().StringVarP(&getBinary.Title, "title", "t", "",
		"Title for binary to search for.")
	getBinaryCmd.MarkFlagRequired("title")

	rootCmd.AddCommand(deletebinaryCmd)
	deletebinaryCmd.Flags().StringVarP(&deleteBinary.Title, "title", "t", "",
		"Title for binary to delete.")
	deletebinaryCmd.MarkFlagRequired("title")
}

var createBinaryCmd = &cobra.Command{
	Use:   commands[CreateBinaryCommand].Use,
	Short: commands[CreateBinaryCommand].Short,
	Long:  commands[CreateBinaryCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		binary, ok := store.Binary[createBinary.GetTitle()]
		if ok {
			createBinary.Version = binary.Version + 1
		} else {
			createBinary.Version = 1
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+(*jwt))

		res, err := cliBin.CreateBinary(newCtx, &createBinary)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to created binary: %w", err), "code", st.Code(), "message", st.Message())
			return
		}

		store.Binary[createBinary.GetTitle()] = &model.Binary{
			Title:   createBinary.GetTitle(),
			Data:    createBinary.GetData(),
			Meta:    createBinary.GetMeta(),
			Version: createBinary.GetVersion(),
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("create binary for title %s", createBinary.GetTitle()), "status", res.GetStatus())
		return
	},
}

var getBinaryCmd = &cobra.Command{
	Use:   commands[GetBinaryCommand].Use,
	Short: commands[GetBinaryCommand].Short,
	Long:  commands[GetBinaryCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		binary, ok := store.Binary[getBinary.Title]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Local version for binary data: title: %s, data: %s, meta: %s, version: %d. "+
				"Make sure you have the latest version by synchronizing local storage",
				binary.Title, binary.Data, binary.Meta, binary.Version)
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliBin.GetBinary(newCtx, &pb.GetBinaryRequest{
			Title: getBinary.Title,
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to get binary: %w", err), "code", st.Code())
			return
		}

		store.Binary[getBinary.Title] = &model.Binary{
			Title:   res.GetTitle(),
			Data:    res.GetData(),
			Meta:    res.GetMeta(),
			Version: res.GetVersion(),
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		msg := fmt.Sprintf("Server version for binary: title: %s, data: %s, meta: %s, version: %d. "+
			"Make sure you have the latest version by synchronizing local storage",
			res.GetTitle(), res.GetData(), res.GetMeta(), res.GetVersion())
		log.Info(msg)
		return
	},
}

var deletebinaryCmd = &cobra.Command{
	Use:   commands[DeleteBinaryCommand].Use,
	Short: commands[DeleteBinaryCommand].Short,
	Long:  commands[DeleteBinaryCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		u, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		_, ok := store.Binary[deleteBinary.GetTitle()]
		// local version doesn't exist: nothing to delete
		if !ok {
			msg := fmt.Sprintf("Nothing found for title: %s. "+
				"Make sure you have the latest version by synchronizing your local storage.",
				deleteBinary.GetTitle())
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliBin.DeleteBinary(newCtx, &pb.DeleteBinaryRequest{
			Title: deleteBinary.Title,
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to delete binary: %w", err), "code", st.Code())
			return
		}

		delete(storage.Objects[u.Username].Binary, deleteBinary.GetTitle())

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("delete binary for title %s!", deleteBinary.GetTitle()), "status", res.GetStatus())
		return
	},
}
