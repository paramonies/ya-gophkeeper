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
	createCard pb.CreateCardRequest
	getCard    pb.GetCardRequest
	deleteCard pb.DeleteCardRequest
)

func init() {
	rootCmd.AddCommand(createCardCmd)
	createCardCmd.Flags().StringVarP(&createCard.Number, "number", "n", "", "Number to save.")
	createCardCmd.Flags().StringVarP(&createCard.Owner, "owner", "o", "",
		"Card owner to save.")
	createCardCmd.Flags().StringVarP(&createCard.ExpDate, "exp_date", "e", "",
		"Card expiration date to save.")
	createCardCmd.Flags().StringVarP(&createCard.Cvv, "cvv", "c", "", "Card cvv to save")
	createCardCmd.Flags().StringVarP(&createCard.Meta, "meta", "m", "",
		"Meta info for the saved card. Optional.")

	createCardCmd.MarkFlagRequired("number")
	createCardCmd.MarkFlagRequired("owner")
	createCardCmd.MarkFlagRequired("exp_date")
	createCardCmd.MarkFlagRequired("cvv")

	rootCmd.AddCommand(getCardCmd)
	getCardCmd.Flags().StringVarP(&getCard.Number, "number", "n", "",
		"Card number to search for.")
	getCardCmd.MarkFlagRequired("number")

	rootCmd.AddCommand(deleteCardCmd)
	deleteCardCmd.Flags().StringVarP(&deleteCard.Number, "number", "n", "",
		"Card number to delete.")
	deleteCardCmd.MarkFlagRequired("number")
}

var createCardCmd = &cobra.Command{
	Use:   commands[CreateCardCommand].Use,
	Short: commands[CreateCardCommand].Short,
	Long:  commands[CreateCardCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		card, ok := store.Card[createCard.GetNumber()]
		if ok {
			createCard.Version = card.Version + 1
		} else {
			createCard.Version = 1
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+(*jwt))

		res, err := cliCard.CreateCard(newCtx, &createCard)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to created card: %w", err), "code", st.Code(), "message", st.Message())
			return
		}

		store.Card[createCard.GetNumber()] = &model.Card{
			Number:  createCard.GetNumber(),
			Owner:   createCard.GetOwner(),
			ExpDate: createCard.GetExpDate(),
			Cvv:     createCard.GetCvv(),
			Meta:    createCard.GetMeta(),
			Version: createCard.GetVersion(),
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("created card with number %s", createCard.GetNumber()), "status", res.GetStatus())
		return
	},
}

var getCardCmd = &cobra.Command{
	Use:   commands[GetCardCommand].Use,
	Short: commands[GetCardCommand].Short,
	Long:  commands[GetCardCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		_, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		card, ok := store.Card[getCard.Number]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Local version for card data: "+
				"number: %s, owner: %s, exp-date: %s, cvv: %s, meta: %s, version: %d. "+
				"Make sure you have the latest version by synchronizing local storage",
				card.Number, card.Owner, card.ExpDate, card.Cvv, card.Meta, card.Version)
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliCard.GetCard(newCtx, &pb.GetCardRequest{
			Number: getCard.GetNumber(),
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to get card: %w", err), "code", st.Code())
			return
		}

		store.Card[getCard.Number] = &model.Card{
			Number:  res.GetNumber(),
			Owner:   res.GetOwner(),
			ExpDate: res.GetExpDate(),
			Cvv:     res.GetCvv(),
			Meta:    res.GetMeta(),
			Version: res.GetVersion(),
		}

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		msg := fmt.Sprintf("Server version for card: "+
			"number: %s, owner: %s, exp-date: %s, cvv: %s, meta: %s, version: %d. "+
			"Make sure you have the latest version by synchronizing local storage",
			res.GetNumber(), res.GetOwner(), res.GetExpDate(), res.GetCvv(), res.GetMeta(), res.GetVersion())
		log.Info(msg)
		return
	},
}

var deleteCardCmd = &cobra.Command{
	Use:   commands[DeleteCardCommand].Use,
	Short: commands[DeleteCardCommand].Short,
	Long:  commands[DeleteCardCommand].Long,
	Run: func(cmd *cobra.Command, args []string) {
		u, store, jwt, err := getUserInfo()
		if err != nil {
			log.Fatal(err)
		}

		_, ok := store.Card[deleteCard.GetNumber()]
		// local version doesn't exist: nothing to delete
		if !ok {
			msg := fmt.Sprintf("Nothing found for card number: %s. "+
				"Make sure you have the latest version by synchronizing your local storage.",
				deleteCard.GetNumber())
			log.Info(msg)
			return
		}

		// local version not found - search on server
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+*jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		res, err := cliCard.DeleteCard(newCtx, &pb.DeleteCardRequest{
			Number: deleteCard.GetNumber(),
		})
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(fmt.Sprintf("failed to delete card: %w", err), "code", st.Code())
			return
		}

		delete(storage.Objects[u.Username].Card, deleteCard.GetNumber())

		err = storage.UpdateFiles(cfg.UsersStoragePath, cfg.ObjectsStoragePath)
		if err != nil {
			log.Error("failed to update local storage files", err)
			return
		}

		log.Info(fmt.Sprintf("delete card with %s number", deleteCard.GetNumber()), "status", res.GetStatus())
		return
	},
}
