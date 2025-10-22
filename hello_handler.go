package main

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func helloHandler(
	ctx context.Context,
	apiClient *slack.Client,
	ev *slackevents.MessageEvent,
) error {
	_, _, err := apiClient.PostMessageContext(
		ctx,
		ev.Channel,
		slack.MsgOptionText("Hi there🐙☁️", false), slack.MsgOptionTS(ev.ThreadTimeStamp))
	return err
}
