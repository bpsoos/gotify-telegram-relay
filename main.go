package main

import (
	"time"

	configpkg "github.com/bpsoos/gotify-telegram-relay/config"
	"github.com/bpsoos/gotify-telegram-relay/logging"
	workerkpkg "github.com/bpsoos/gotify-telegram-relay/worker"
)

var logger = logging.GetLogger()

func main() {
	config := configpkg.LoadConfig()

	worker := workerkpkg.NewWorker(
		&workerkpkg.Config{
			TelegramToken:                config.TelegramToken,
			TelegramChatId:               config.TelegramChatId,
			GotifyClientToken:            config.GotifyClientToken,
			GotifyBaseUrl:                config.GotifyBaseUrl,
			MessageBatchCount:            config.MessageBatchCount,
			PostRelayMessageSleepSeconds: config.PostRelayMessageSleepSeconds,
			Client:                       config.Client,
			GotifyTokenHeaderKey:         config.GotifyTokenHeaderKey,
		},
	)

	for {
		worker.Work()
		logger.Info("sleeping")
		time.Sleep(time.Duration(config.WorkSleepSeconds) * time.Second)
	}
}
