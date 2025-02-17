package main

import (
	"context"
	"log"
	"np-notify/pkg/aws_helpers"
	"np-notify/pkg/fetcher"
	"np-notify/pkg/models"
	"np-notify/pkg/utils"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
)

type GameChangeStatus int

const (
	Changed GameChangeStatus = iota
	Unchanged
	Zeroed
)

func handleRequest(ctx context.Context) error {
	var db aws_helpers.TableBasics
	var snsActions aws_helpers.SnsActions
	var err error

	// ctx := context.TODO()
	f := fetcher.Fetcher{ApiURL: "https://np.ironhelmet.com"}

	snsTopicArn := os.Getenv("SNS_ARN")
	if snsTopicArn == "" {
		utils.LogError("SNS_TOPIC_ARN is not set")
		panic("SNS_TOPIC_ARN is not set")
	}
	tableName := os.Getenv("DYNAMO_TABLE_NAME")
	if snsTopicArn == "" {
		utils.LogError("DYNAMO_TABLE_NAME is not set")
		panic("DYNAMO_TABLE_NAME is not set")
	}

	if os.Getenv("LOCAL_STACK") == "" {
		db, snsActions, err = aws_helpers.AWS(ctx, tableName)
	} else {
		db, snsActions, err = aws_helpers.LocalAWS(ctx, tableName)
	}
	if err != nil {
		utils.LogError("Failed to setup AWS")
		panic(err)
	}

	accounts, err := models.GetAllAccounts(ctx, db)
	if err != nil {
		utils.LogError("Failed to get Accounts")
		panic(err)
	}

	mailBody := ""

	for _, acc := range accounts {

		games, err := f.GetPlayerGames(acc)
		if err != nil {
			log.Fatalf("Failed to get player games: %v", err)
		}

		gameChanges := make([]GameChangeStatus, len(games))

		previousUnreads := acc.UnreadCounts
		if err != nil {
			log.Fatalf("Failed to get previous unreads: %v", err)
		}
		if previousUnreads == nil {
			previousUnreads = make(map[string]models.UnreadCount)
		}

		for i, game := range games {
			if !game.IsActive() {
				gameChanges[i] = GameChangeStatus(Unchanged)
				continue
			}
			unreads, err := f.GetUnreadCount(acc, game.ID)
			if err != nil {
				log.Fatalf("Failed to get unread count: %v", err)
			}
			game.Unreads = unreads

			if _, ok := previousUnreads[game.ID]; !ok {
				previousUnreads[game.ID] = models.UnreadCount{}
			}

			if models.CompareUnreadCounts(unreads, previousUnreads[game.ID]) {
				gameChanges[i] = GameChangeStatus(Unchanged)
				continue
			}
			gameChanges[i] = GameChangeStatus(Changed)

			previousUnreads[game.ID] = unreads
			if !unreads.AllZero() {
				events, err := f.GetGameEvents(acc, game)
				if err != nil {
					log.Fatalf("Failed to get game events: %v", err)
				}

				game.Events = events
				messages, err := f.GetGameDiplomacy(acc, game)
				if err != nil {
					log.Fatalf("Failed to get game messages: %v", err)
				}
				game.Diplomacy = messages

				// // Sadly fethcing the unread message "thread" comments, sets them to read status. :(
			}
			mailBody += game.HtmlString()
			mailBody += "\n\n"
		}

		for _, chagameChange := range gameChanges {
			if chagameChange == Changed {
				snsActions.Publish(ctx, snsTopicArn, mailBody, "", "", "", "")
				acc.UnreadCounts = previousUnreads
				acc.UpdateUnreads(ctx, db)
				break
			} else if chagameChange == Zeroed {
				acc.UnreadCounts = previousUnreads
				acc.UpdateUnreads(ctx, db)
			}
		}
	}
	return nil

}

func main() {

	if os.Getenv("LOCAL") == "" {
		lambda.Start(handleRequest)
	} else {
		handleRequest(context.Background())
	}
}
