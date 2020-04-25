package context

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//Context of update
type UpdateContext struct {
	Update       *tgbotapi.Update
	Bot          *tgbotapi.BotAPI
	ChatId       int64
	FromId       int
	Text         string
	CallbackData string
	Contact      *tgbotapi.Contact
	Location     *tgbotapi.Location
}

//Send text message in context of current update (to chat id)
func (ctx *UpdateContext) SendText(text string, markup ...tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(ctx.ChatId, text)
	if len(markup) == 1 {
		msg.ReplyMarkup = markup[0]
	}
	ctx.Bot.Send(msg)
}

//Sends location in current context
func (ctx *UpdateContext) SendLocation(lat, lng float64) {
	ctx.Bot.Send(tgbotapi.NewLocation(ctx.ChatId, lat, lng))
}
