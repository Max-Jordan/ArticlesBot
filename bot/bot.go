package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Max-Jordan/ArticlesBot/lib/e"
)

const (
	API                       = "api.telegram.org"
	sendMessageMthod          = "sendMessage"
	getUpdatesMethod          = "getUpdates"
	answerCallBackQueryMethod = "answerCallbackQuery"
)

type Bot struct {
	token    string
	api      string
	basePath string
	client   *http.Client
	botStop  chan any
}

func NewBot(token string) Bot {
	return Bot{
		token:    token,
		api:      API,
		basePath: "bot" + token,
		client:   http.DefaultClient,
		botStop:  make(chan any),
	}
}

func (b *Bot) StartPolling(config UpdateConfig) UpdateChan {
	ch := make(chan Update, 100)
	go func() {
		for {
			select {
			case <-b.botStop:
				close(ch)
				return
			default:
			}
			updates, err := b.getUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get update. Try again in 3 seconds")
				time.Sleep(3 * time.Second)
				continue
			}

			for _, update := range updates {
				if update.UpdateId >= config.Offset {
					config.Offset = update.UpdateId + 1
					ch <- update
				}
			}
		}
	}()
	return ch
}

func (b *Bot) StopPolling() {
	close(b.botStop)
}

func (b Bot) SendMessage(chatId int, text string) (Message, error) {
	val := url.Values{}
	val.Add("chat_id", strconv.Itoa(chatId))
	val.Add("text", text)
	resp := b.doRequest(http.MethodPost, sendMessageMthod, val)
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

func (b Bot) ReplyToMessage(message Message, text string) (Message, error) {
	repParam := ReplyParameters{
		MessageId: message.MessageId,
		ChatId:    message.Chat.ChatId,
	}
	data, err := json.Marshal(repParam)
	if err != nil {
		log.Println(e.Wrap("Marshalling json error", err))
		return Message{}, err
	}
	val := url.Values{}
	val.Add("chat_id", strconv.Itoa(message.Chat.ChatId))
	val.Add("text", text)
	val.Add("reply_parameters", string(data))
	resp := b.doRequest(http.MethodPost, sendMessageMthod, val)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(e.Wrap("Repling messag error", err))
		return Message{}, nil
	}
	var mess Message
	if err = json.Unmarshal(data, &mess); err != nil {
		log.Println("Error unmarshalling message", err)
		return Message{}, nil
	}
	return mess, nil
}

func (b Bot) SendInlineKeyBoard(message Message) (Message, error) {
	btn := KeyboardButton{Text: "Save article", CallBack: "save_article"}
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]KeyboardButton{{btn}},
	}
	data, err := json.Marshal(keyboard)
	if err != nil {
		log.Println(e.Wrap("Marshalling data erro ", err))
		return Message{}, err
	}
	val := url.Values{}
	val.Add("chat_id", strconv.Itoa(message.Chat.ChatId))
	val.Add("text", "Your keyboard")
	val.Add("reply_markup", string(data))
	resp := b.doRequest(http.MethodPost, sendMessageMthod, val)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	var mess Message
	if err = json.Unmarshal(data, &mess); err != nil {
		log.Println(e.Wrap("Unmarshalling message error", err))
		return Message{}, err
	}
	return mess, nil
}

func (b Bot) AnswerCallBackQuery(query CallbackQuery, text string) {
	val := url.Values{}
	val.Add("callback_query_id", query.Id)
	val.Add("text", text)
	resp := b.doRequest(http.MethodPost, answerCallBackQueryMethod, val)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
	}
}

func (b Bot) getUpdates(config UpdateConfig) ([]Update, error) {
	val := url.Values{}
	val.Add("offset", strconv.Itoa(config.Offset))
	val.Add("limit", strconv.Itoa(config.Limit))
	val.Add("timeout", strconv.Itoa(config.Timeout))

	resp := b.doRequest(http.MethodGet, getUpdatesMethod, val)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	fmt.Println(string(data))
	if err != nil {
		return nil, e.Wrap("Getting updates error", err)
	}
	var respUp ResponseUpdate
	if err = json.Unmarshal(data, &respUp); err != nil {
		return nil, e.Wrap("Unpacking response error", err)
	}
	return respUp.Results, nil
}

func (b Bot) doRequest(httpMethod string, tgMethod string, values url.Values) *http.Response {
	u := url.URL{
		Scheme: "https",
		Host:   API,
		Path:   b.basePath,
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest(httpMethod, u.JoinPath(tgMethod).String(), nil)
	if err != nil {
		log.Println(e.Wrap("Requesting error", err))
		return nil
	}
	resp, err := b.client.Do(req)
	if err != nil {
		log.Println(e.Wrap("Gettin response error", err))
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		return nil
	}
	return resp
}
