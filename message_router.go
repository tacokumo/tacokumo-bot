package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type MessageRouter struct {
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{}
}

type RouterHandler func(
	ctx context.Context,
	apiClient *slack.Client,
	ev *slackevents.MessageEvent,
) error

func (mr *MessageRouter) Determine(text string) (RouterHandler, error) {
	parts := strings.Split(text, " ")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	// メンションがついている場合はすべて取り外す
	for i := 0; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "<@") && strings.HasSuffix(parts[i], ">") {
			parts[i] = ""
		}
	}
	parts = strings.Split(strings.TrimSpace(strings.Join(parts, " ")), " ")

	switch parts[0] {
	case "hello":
		return helloHandler, nil
	case "help":
		return helpHandler, nil
	}

	return nil, fmt.Errorf("no handler found for message: %s", text)
}
