package handler

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	b "github.com/Max-Jordan/ArticlesBot/bot"
	"github.com/Max-Jordan/ArticlesBot/lib/e"
	"github.com/Max-Jordan/ArticlesBot/services"
)

const (
	waitingURL         = "waitingURL"
	waitingDescription = "waitingDesc"
)

var Storage services.PostgreStorage

type handler struct {
	bot      *b.Bot
	commands map[string]b.CommandHandler
	state    map[int64]string
	draft    map[int64]*services.Article
	storage  *services.PostgreStorage
}

func NewHandler(bot *b.Bot, storage *services.PostgreStorage) *handler {
	h := &handler{bot: bot, storage: storage}

	h.commands = map[string]b.CommandHandler{
		"/start":    h.startCommandHandler,
		"/keyboard": h.keyboardCommandHandler,
	}
	h.state = make(map[int64]string)
	h.draft = make(map[int64]*services.Article)

	return h
}

func (h *handler) StartHandling() {
	updates := h.bot.GetUpdates()
	for update := range updates {
		if update.Message != nil {
			if h.state[update.Message.From.UserId] == waitingURL || h.state[update.Message.From.UserId] == waitingDescription {
				h.createArticle(update)
				continue
			}
			f, ok := h.commands[update.Message.Text]
			if !ok {
				h.bot.SendMessage(*update.Message, "Command not found")
				continue
			} else {
				f(update)
			}
		}
		if update.CallBack != nil {
			h.bot.AnswerCallBackQuery(*update.CallBack, "Loading...")
			switch update.CallBack.Data {
			case "save_article":
				h.state[update.CallBack.From.UserId] = waitingURL
				h.draft[update.CallBack.From.UserId] = services.NewEptyArticle()
				h.bot.SendMessage(update.CallBack.Message, "Write URL of article: ")
			case "list":
				list, err := h.storage.List(update.CallBack.From.UserId)
				if err != nil {
					log.Println(e.Wrap("Getting article list error", err))
				}
				fmt.Println(articlesToString(list))
				h.bot.SendMessage(update.CallBack.Message, articlesToString(list))
			case "del":
				h.bot.SendMessage(update.CallBack.Message, "This command will delete article")
			case "complete":
				
			}
		}

	}
}

func articlesToString(list []services.Article) string {
	var result strings.Builder
	for _, article := range list {
		_, err := result.WriteString(article.ToString())
		if err != nil {
			log.Println(err)
		}
	}
	return result.String()
}

func (h *handler) createArticle(update b.Update) {
	id := update.Message.From.UserId

	a, ok := h.draft[id]
	if !ok {
		h.bot.SendMessage(*update.Message, "Draft not found.")
		delete(h.state, id)
	}
	switch h.state[id] {
	case waitingURL:
		if !isValidURL(strings.ToLower(update.Message.Text)) {
			h.bot.SendMessage(*update.Message, "Invalid URL. Try again")
			return
		}
		a.URL = strings.ToLower(update.Message.Text)
		h.state[id] = waitingDescription
		h.draft[id] = a
		h.bot.SendMessage(*update.Message, b.EscapeMarkdownV2("Write article decription"))
	case waitingDescription:
		a.Description = update.Message.Text
		a.UserId = id
		a.Readed = false
		if err := h.storage.Save(a); err != nil {
			log.Println(err)
			h.bot.SendMessage(*update.Message, b.EscapeMarkdownV2("Something went wrong. Try again"))
			return
		}
		fmt.Println(*update.Message)
		h.bot.SendMessage(*update.Message, b.EscapeMarkdownV2("You article was saved."))
		delete(h.state, id)
		delete(h.draft, id)
	}
}

func (h handler) startCommandHandler(update b.Update) error {
	if _, err := h.bot.SendMessage(*update.Message, "Hello"); err != nil {
		return err
	}
	return nil
}

func (h handler) keyboardCommandHandler(update b.Update) error {
	btnSave := b.NewKeyboardBtn("Save article", "save_article")
	btnList := b.NewKeyboardBtn("List articles", "list")
	btnDel := b.NewKeyboardBtn("Delete article", "del")
	btnComplete := b.NewKeyboardBtn("Complete article", "complete")
	keyboard := b.NewInlineKeyBoard()
	keyboard.InlineKeyboard = [][]b.KeyboadBtn{{btnSave, btnComplete}, {btnList, btnDel}}
	if _, err := h.bot.SendInlineKeyBoard(*update.Message, *keyboard); err != nil {
		return err
	}
	return nil
}

func isValidURL(URL string) bool {
	u, err := url.ParseRequestURI(URL)
	if err != nil {
		log.Println(e.Wrap("Parse URL error", err))
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}
