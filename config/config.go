package config

import (
	"net/http"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	GotifyBaseUrl                string
	GotifyClientToken            string
	GotifyAppId                  string
	TelegramToken                string
	TelegramChatId               string
	Client                       *http.Client
	WorkSleepSeconds             int
	MessageBatchCount            int
	PostRelayMessageSleepSeconds int
	GotifyTokenHeaderKey         string
	DisableTelegramNotification  bool
}

const (
	defaultWorkSleepSeconds             = 15
	defaultMessageBatchCount            = 50
	defaultPostRelayMessageSleepSeconds = 1
	gotifyTokenHeaderKey                = "X-Gotify-Key"
)

func defaultInt(k *koanf.Koanf, key string, defVal int) int {
	if !k.Exists(key) {
		return defVal
	}

	return k.Int(key)
}

func LoadConfig() *Config {
	var k = koanf.New(".")
	k.Load(env.Provider(
		"GTRELAY_",
		"__",
		func(s string) string {
			return strings.TrimPrefix(s, "GTRELAY_")
		},
	), nil)

	return &Config{
		GotifyBaseUrl:                k.MustString("GOTIFY_BASE_URL"),
		GotifyClientToken:            k.MustString("GOTIFY_CLIENT_TOKEN"),
		GotifyAppId:                  k.MustString("GOTIFY_APP_ID"),
		TelegramToken:                k.MustString("TELEGRAM_TOKEN"),
		TelegramChatId:               k.MustString("TELEGRAM_CHAT_ID"),
		Client:                       http.DefaultClient,
		WorkSleepSeconds:             defaultInt(k, "WORK_SLEEP_SECONDS", defaultWorkSleepSeconds),
		MessageBatchCount:            defaultInt(k, "MESSAGE_BATCH_COUNT", defaultMessageBatchCount),
		PostRelayMessageSleepSeconds: defaultInt(k, "POST_RELAY_MESSAGE_SLEEP_SECONDS", defaultPostRelayMessageSleepSeconds),
		GotifyTokenHeaderKey:         gotifyTokenHeaderKey,
		DisableTelegramNotification:  k.Bool("DISABLE_TELEGRAM_NOTIFICATION"),
	}
}
