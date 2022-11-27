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
	createText pb.CreateTextRequest
	getText    pb.GetTextRequest
	deleteText pb.DeleteTextRequest
)

func init() {
	rootCmd.AddCommand(createTextCmd)
	createTextCmd.Flags().StringVarP(&createText.Title, "title", "t", "", "Title to save.")
	createTextCmd.Flags().StringVarP(&createText.Data, "data", "d", "",
		"Text data to save.")
	createTextCmd.Flags().StringVarP(&createText.Meta, "meta", "m", "",
		"Meta info for the saved text. Optional.")
	createTextCmd.MarkFlagRequired("title")
	createTextCmd.MarkFlagRequired("data")

	rootCmd.AddCommand(getTextCmd)
	getTextCmd.Flags().StringVarP(&getText.Title, "title", "t", "",
		"Title for text to search for.")
	getTextCmd.MarkFlagRequired("title")

	rootCmd.AddCommand(deleteTextCmd)
	deleteTextCmd.Flags().StringVarP(&deleteText.Title, "title", "t", "",
		"Title for text to delete.")
	deleteTextCmd.MarkFlagRequired("title")
}

var createTextCmd = &cobra.Command{
	Use:   commands[CreateTextCommand].Use,
	Short: commands[CreateTextCommand].Short,
	Long:  commands[CreateTextCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		text, ok := store.Text[createText.GetTitle()]
		if ok {
			createText.Version = text.Version + 1
		} else {
			createText.Version = 1
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+(*jwt))

		res, err := cliText.CreateText(newCtx, &createText)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to created text: %w", err), "code", st.Code(), "message", st.Message())
			return
		}

		store.Text[createText.GetTitle()] = &model.Text{
			Title:   createText.GetTitle(),
			Data:    createText.GetData(),
			Meta:    createText.GetMeta(),
			Version: createText.GetVersion(),
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("create text for title %s", createText.GetTitle()), "status", res.GetStatus())
		return
	},
}

var getTextCmd = &cobra.Command{
	Use:   commands[GetTextCommand].Use,
	Short: commands[GetTextCommand].Short,
	Long:  commands[GetTextCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		text, ok := store.Text[getText.Title]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Local version for text data: title: %s, data: %s, meta: %s, version: %d. "+
				"Make sure you have the latest version by synchronizing local storage",
				text.Title, text.Data, text.Meta, text.Version)
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliText.GetText(newCtx, &pb.GetTextRequest{
			Title: getText.Title,
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to get text: %w", err), "code", st.Code())
			return
		}

		store.Text[getText.Title] = &model.Text{
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

		msg := fmt.Sprintf("Server version for text: title: %s, data: %s, meta: %s, version: %d. "+
			"Make sure you have the latest version by synchronizing local storage",
			res.GetTitle(), res.GetData(), res.GetMeta(), res.GetVersion())
		log.Info(msg)
		return
	},
}

var deleteTextCmd = &cobra.Command{
	Use:   commands[DeleteTextCommand].Use,
	Short: commands[DeleteTextCommand].Short,
	Long:  commands[DeleteTextCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		u, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		_, ok := store.Text[deleteText.GetTitle()]
		// local version doesn't exist: nothing to delete
		if !ok {
			msg := fmt.Sprintf("Nothing found for title: %s. "+
				"Make sure you have the latest version by synchronizing your local storage.",
				deleteText.GetTitle())
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliText.DeleteText(newCtx, &pb.DeleteTextRequest{
			Title: deleteText.Title,
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to delete text: %w", err), "code", st.Code())
			return
		}

		delete(storage.Objects[u.Username].Text, deleteText.GetTitle())

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("delete text for title %s!", deleteText.GetTitle()), "status", res.GetStatus())
		return
	},
}
