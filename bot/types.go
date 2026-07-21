package bot

type ResponseUpdate struct {
	Ok      bool     `json:"ok"`
	Results []Update `json:"result"`
}

type Update struct {
	UpdateId int            `json:"update_id"`
	Message  *Message       `json:"message"`
	CallBack *CallbackQuery `json:"callback_query"`
}

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

type CallbackQuery struct {
	Id      string  `json:"id"`
	From    User    `json:"from"`
	Data    string  `json:"data"`
	Message Message `json:"message"`
}

type Message struct {
	MessageId      int                   `json:"message_id"`
	Text           string                `json:"text"`
	From           *User                 `json:"from"`
	Chat           *Chat                 `json:"chat"`
	ReplyToMessage *Message              `json:"reply_to_message"`
	ReplyMarkUp    *InlineKeyboardMarkup `json:"reply_markup"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]KeyboadBtn `json:"inline_keyboard"`
}

type KeyboadBtn struct {
	Text     string `json:"text"`
	CallBack string `json:"callback_data"`
}

type ReplyParameters struct {
	MessageId int `json:"message_id"`
	ChatId    int `json:"chat_id"`
}

type User struct {
	UserId   int64  `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Name     string `json:"first_name"`
	UserName string `json:"username"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type Params map[string]string

type UpdateChan chan Update
