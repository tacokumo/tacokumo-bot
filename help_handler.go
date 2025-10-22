package main

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func helpHandler(
	ctx context.Context,
	apiClient *slack.Client,
	ev *slackevents.MessageEvent,
) error {
	helpText := "```\n" + `tacokumo-bot: TACOKUMOを管理するBotです
Usage:
  hello          - さあ、あなたもtacokumoに挨拶しよう!
  help           - ヘルプメッセージを表示します、今もう見ているよ!` + "\n```"

	_, _, err := apiClient.PostMessage(ev.Channel, slack.MsgOptionText(helpText, false))
	if err != nil {
		return err
	}

	return nil
}
