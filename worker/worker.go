package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bpsoos/gotify-telegram-relay/logging"
)

var logger = logging.GetLogger()

type Worker struct {
	gotifyBaseUrl                string
	gotifyClientToken            string
	gotifyAppId                  string
	telegramToken                string
	telegramChatId               string
	messageBatchCount            int
	postRelayMessageSleepSeconds int
	client                       *http.Client
	gotifyTokenHeaderKey         string
	disableTelegramNotification  bool
}

type Config struct {
	GotifyBaseUrl                string
	GotifyClientToken            string
	GotifyAppId                  string
	TelegramToken                string
	TelegramChatId               string
	MessageBatchCount            int
	PostRelayMessageSleepSeconds int
	Client                       *http.Client
	GotifyTokenHeaderKey         string
	DisableTelegramNotification  bool
}

func NewWorker(config *Config) *Worker {
	return &Worker{
		gotifyBaseUrl:                config.GotifyBaseUrl,
		gotifyClientToken:            config.GotifyClientToken,
		gotifyAppId:                  config.GotifyAppId,
		telegramToken:                config.TelegramToken,
		telegramChatId:               config.TelegramChatId,
		messageBatchCount:            config.MessageBatchCount,
		postRelayMessageSleepSeconds: config.PostRelayMessageSleepSeconds,
		client:                       config.Client,
		gotifyTokenHeaderKey:         config.GotifyTokenHeaderKey,
		disableTelegramNotification:  config.DisableTelegramNotification,
	}
}

func (w Worker) Work() {
	fetchUrl := w.gotifyBaseUrl + "/application/" + w.gotifyAppId + "/message?limit=" + strconv.Itoa(w.messageBatchCount)
	for {
		messagesRes := w.fetchMessages(fetchUrl)
		for i := range messagesRes.Messages {
			w.relayMessage(messagesRes.Messages[i].Message)
			w.deleteMessage(messagesRes.Messages[i].ID)
			time.Sleep(time.Duration(w.postRelayMessageSleepSeconds) * time.Second)
		}
		if messagesRes.Paging.Next == "" {
			break
		}
		fetchUrl = messagesRes.Paging.Next
	}
}
func (w Worker) fetchMessages(fetchUrl string) *MessagesResponse {
	req, err := http.NewRequest(http.MethodGet, fetchUrl, nil)
	headers := http.Header{}
	headers.Add(w.gotifyTokenHeaderKey, w.gotifyClientToken)
	req.Header = headers
	if err != nil {
		log.Fatal(err)
	}

	res, err := w.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		logger.Error("fetch response bad status", "response", string(body))
		os.Exit(1)
	}
	messagesRes := new(MessagesResponse)
	err = json.NewDecoder(res.Body).Decode(messagesRes)
	if err != nil {
		log.Fatal(err)
	}
	return messagesRes
}

func (w Worker) relayMessage(message string) {
	logger.Info("relaying message", "message", message)
	body := fmt.Sprintf(
		`{"chat_id": "%s", "text": "%s", "disable_notification": %t}`,
		w.telegramChatId,
		message,
		w.disableTelegramNotification,
	)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+w.telegramToken+"/sendMessage",
		strings.NewReader(body),
	)
	if err != nil {
		log.Fatal(err)
	}
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	req.Header = headers
	res, err := w.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		logger.Error("relay response bad status", "response", string(body))
		os.Exit(1)
	}
}

func (w Worker) deleteMessage(id int) {
	logger.Info("deleting message", "message_id", id)
	url := w.gotifyBaseUrl + "/message/" + strconv.Itoa(id)
	req, err := http.NewRequest(
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	headers := http.Header{}
	headers.Add(w.gotifyTokenHeaderKey, w.gotifyClientToken)
	req.Header = headers

	res, err := w.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		logger.Error("delete response bad status", "response", string(body))
		os.Exit(1)
	}
}

type Message struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}

type Paging struct {
	Limit string `json:"message"`
	Next  string `json:"next"`
	Since int    `json:"since"`
	Size  int    `json:"size"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
	Paging   Paging    `json:"paging"`
}
