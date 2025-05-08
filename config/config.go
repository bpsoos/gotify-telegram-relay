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
	TelegramToken                string
	TelegramChatId               string
	Client                       *http.Client
	WorkSleepSeconds             int
	MessageBatchCount            int
	PostRelayMessageSleepSeconds int
	GotifyTokenHeaderKey         string
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
		TelegramToken:                k.MustString("TELEGRAM_TOKEN"),
		TelegramChatId:               k.MustString("TELEGRAM_CHAT_ID"),
		Client:                       http.DefaultClient,
		WorkSleepSeconds:             int64(defaultInt(k, "WORK_SLEEP_SECONDS", defaultWorkSleepSeconds)),
		MessageBatchCount:            defaultInt(k, "MESSAGE_BATCH_COUNT", defaultMessageBatchCount),
		PostRelayMessageSleepSeconds: int64(defaultInt(k, "POST_RELAY_MESSAGE_SLEEP_SECONDS", defaultPostRelayMessageSleepSeconds)),
		GotifyTokenHeaderKey:         gotifyTokenHeaderKey,
	}
}
