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

	logger.Info(
		"starting worker with config",
		"gotify_app_id", config.GotifyAppId,
		"gotify_base_url", config.GotifyBaseUrl,
		"message_batch_count", config.MessageBatchCount,
		"post_relay_message_sleep_seconds", config.PostRelayMessageSleepSeconds,
		"disable_telegram_notification", config.DisableTelegramNotification,
	)

	worker := workerkpkg.NewWorker(
		&workerkpkg.Config{
			TelegramToken:                config.TelegramToken,
			TelegramChatId:               config.TelegramChatId,
			GotifyClientToken:            config.GotifyClientToken,
			GotifyBaseUrl:                config.GotifyBaseUrl,
			GotifyAppId:                  config.GotifyAppId,
			MessageBatchCount:            config.MessageBatchCount,
			PostRelayMessageSleepSeconds: config.PostRelayMessageSleepSeconds,
			Client:                       config.Client,
			GotifyTokenHeaderKey:         config.GotifyTokenHeaderKey,
			DisableTelegramNotification:  config.DisableTelegramNotification,
		},
	)

	for {
		worker.Work()
		logger.Info("sleeping")
		time.Sleep(time.Duration(config.WorkSleepSeconds) * time.Second)
	}
}
