package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Bot struct {
	apiClient    *slack.Client
	SocketClient *socketmode.Client
	router       MessageRouter
	botUserID    string
}

func NewBot(apiClient *slack.Client, socketClient *socketmode.Client) *Bot {
	bot := &Bot{
		apiClient:    apiClient,
		SocketClient: socketClient,
	}

	// Fetch bot user ID so we can filter to direct mentions.
	if authTest, err := apiClient.AuthTest(); err != nil {
		slog.Error("auth.test failed", "error", err)
	} else {
		bot.botUserID = authTest.UserID
	}

	return bot
}

func (b *Bot) Run(ctx context.Context) error {
	go b.handleEvents(ctx)
	return b.SocketClient.RunContext(ctx)
}

func (b *Bot) handleEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-b.SocketClient.Events:
			switch event.Type {
			case socketmode.EventTypeConnected:
				slog.InfoContext(ctx, "Connected to Slack with Socket Mode.")
			case socketmode.EventTypeConnectionError:
				slog.ErrorContext(ctx, "Connection error occurred.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					slog.ErrorContext(ctx, "Could not type cast the event to EventsAPIEvent")
					continue
				}

				b.SocketClient.Ack(*event.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch innerEventData := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if innerEventData.BotID != "" {
							// Ignore messages from bots
							continue
						}

						if !b.isMentioned(innerEventData) {
							continue
						}

						handler, err := b.router.Determine(innerEventData.Text)
						if err != nil {
							_, _, err := b.apiClient.PostMessageContext(ctx, innerEventData.Channel, slack.MsgOptionText("未定義のコマンドです: `@tacokumo-bot help` で確認してください", false), slack.MsgOptionTS(innerEventData.ThreadTimeStamp))
							if err != nil {
								slog.ErrorContext(ctx, "Failed to post message", "error", err)
							}
							continue
						}
						if err := handler(ctx, b.apiClient, innerEventData); err != nil {
							slog.ErrorContext(ctx, "Handler error", "error", err)
						}
					}
				default:
					slog.WarnContext(ctx, "Unsupported Events API event received", "type", eventsAPIEvent.Type)
				}
			}
		}
	}
}

func (b *Bot) isMentioned(ev *slackevents.MessageEvent) bool {
	if b.botUserID == "" {
		return false
	}

	mention := fmt.Sprintf("<@%s>", b.botUserID)

	// Check both the top-level text and the normalized message text.
	texts := []string{ev.Text}
	if ev.Message != nil {
		texts = append(texts, ev.Message.Text)
	}

	for _, text := range texts {
		if strings.Contains(text, mention) {
			return true
		}
	}
	return false
}
