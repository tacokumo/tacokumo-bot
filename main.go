package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		logger.ErrorContext(ctx, "SLACK_BOT_TOKEN is not set")
	}
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		logger.ErrorContext(ctx, "SLACK_APP_TOKEN is not set")
	}
	slackClient := slack.New(botToken, slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(slackClient, socketmode.OptionDebug(logLevel == slog.LevelDebug))
	b := NewBot(slackClient, socketClient)

	if err := b.Run(ctx); err != nil {
		slog.ErrorContext(ctx, "Bot encountered an error", "error", err)
	}
}
