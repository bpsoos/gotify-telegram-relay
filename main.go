package main

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
)

const (
	workSleepSeconds        = 15
	messageIterSleepSeconds = 1
	messageBatchCount       = 50
	gotifyTokenHeaderKey    = "X-Gotify-Key"
)

func main() {
	gotifyBaseUrl := os.Getenv("GTRELAY_GOTIFY_BASE_URL")
	gotifyClientToken := os.Getenv("GTRELAY_GOTIFY_CLIENT_TOKEN")
	telegramToken := os.Getenv("GTRELAY_TELEGRAM_TOKEN")
	telegramChatId := os.Getenv("GTRELAY_TELEGRAM_CHAT_ID")
	client := http.DefaultClient

	worker := Worker{
		telegramToken:     telegramToken,
		telegramChatId:    telegramChatId,
		gotifyClientToken: gotifyClientToken,
		gotifyBaseUrl:     gotifyBaseUrl,
		client:            client,
	}
	for {
		worker.work()
		println("sleeping")
		time.Sleep(workSleepSeconds * time.Second)
	}
}

type Worker struct {
	gotifyBaseUrl     string
	gotifyClientToken string
	telegramToken     string
	telegramChatId    string
	client            *http.Client
}

func (w Worker) work() {
	fetchUrl := w.gotifyBaseUrl + "/message?limit=" + strconv.Itoa(messageBatchCount)
	for {
		messagesRes := w.fetchMessages(fetchUrl)
		for i := range messagesRes.Messages {
			w.relayMessage(messagesRes.Messages[i].Message)
			w.deleteMessage(messagesRes.Messages[i].ID)
			time.Sleep(messageIterSleepSeconds * time.Second)
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
	headers.Add(gotifyTokenHeaderKey, w.gotifyClientToken)
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
		println(string(body))
		log.Fatal("bad status: ", res.Status)
	}
	messagesRes := new(MessagesResponse)
	err = json.NewDecoder(res.Body).Decode(messagesRes)
	if err != nil {
		log.Fatal(err)
	}
	return messagesRes
}

func (w Worker) relayMessage(message string) {
	body := fmt.Sprintf(
		`{"chat_id": "%s", "text": "%s", "disable_notification": false}`,
		w.telegramChatId,
		message,
	)
	println(body)
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
		println(string(body))
		log.Fatal("bad status: ", res.Status)
	}
}

func (w Worker) deleteMessage(id int) {
	println("deleting message", id)
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
	headers.Add(gotifyTokenHeaderKey, w.gotifyClientToken)
	req.Header = headers

	res, err := w.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		println(string(body))
		log.Fatal("bad status: ", res.Status)
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
