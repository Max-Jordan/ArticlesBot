package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Max-Jordan/ArticlesBot/lib/e"
)

const (
	API              = "api.telegram.org"
	getUpdatesMethod = "getUpdates"
	sendMessageMthod = "sendMessage"
)

type Bot struct {
	token    string
	api      string
	basePath string
	client   *http.Client
}

func NewBot(token, api string) Bot {
	return Bot{
		token:    token,
		api:      api,
		basePath: "bot" + token,
		client:   http.DefaultClient,
	}
}

type ResponseUpdate struct {
	Ok      bool     `json:"ok"`
	Results []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	From      User   `json:"from"`
	SenderChat Chat `json:"sender_chat"`
}

type User struct {
	UserId int    `json:"id"`
	IsBot  bool   `json:"is_bot"`
	Name   string `json:"first_name"`
}

type Chat struct {
	ChatId int `json:"id"`
}

func (b Bot) getUpdates(offset int) ([]Update, error) {
	u := url.URL{
		Scheme: "https",
		Host:   API,
		Path:   b.basePath,
	}
	q := u.Query()
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.JoinPath(getUpdatesMethod).String(), nil)
	if err != nil {
		return nil, e.Wrap("Getting updates error", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, e.Wrap("Getting upate error", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("Getting updates error", err)
	}
	var respUp ResponseUpdate
	if err = json.Unmarshal(data, &respUp); err != nil {
		return nil, e.Wrap("Unpacking response error", err)
	}
	return respUp.Results, nil
}

func (b Bot) startPolling() {
	offset := 0
	for {
		updates, err := b.getUpdates(offset)
		if err != nil {
			log.Print(e.Wrap("Getting updates error", err))
		}
		for _, update := range updates {
			offset = update.UpdateId + 1
			fmt.Println(update)
			b.sendMessage(update.Message.From.UserId, update.Message.Text)
		}
	}
}

func (b Bot) sendMessage(chatId int, text string) (Message, error) {
	u := url.URL{
		Scheme: "https",
		Host:   API,
		Path:   b.basePath,
	}
	q := u.Query()
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodPost, u.JoinPath(sendMessageMthod).String(), nil)
	if err != nil {
		return Message{}, e.Wrap("Sending message error", err)
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return Message{}, e.Wrap("Getting response sending message error", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, e.Wrap("Reading response error", err)
	}
	var mess Message
	if err = json.Unmarshal(body, &mess); err != nil {
		return Message{}, e.Wrap("Unmarshalling response error", err)
	}
	return mess, nil
}

func main() {
	token := os.Getenv("bot_token")
	if token == "" {
		log.Fatal("Token isn't set")
	}
	bot := NewBot(token, API)
	bot.startPolling()
}
