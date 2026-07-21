package bot

import "strings"

func NewInlineKeyBoard() *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{}
}

func NewKeyboardBtn(text, callback_data string) KeyboadBtn {
	return KeyboadBtn{
		Text:     text,
		CallBack: callback_data,
	}
}

type CommandHandler func(up Update) error

func EscapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		"`", "\\`",
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}
