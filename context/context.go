package context

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
	TgBot        TgBot
}

//Send text message in context of current update (to chat id)
func (ctx *UpdateContext) SendText(text string, markup ...interface{}) {
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

//Sends message as html
func (ctx *UpdateContext) SendHTML(text string, markup ...interface{}) {
	msg := tgbotapi.NewMessage(ctx.ChatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if len(markup) == 1 {
		msg.ReplyMarkup = markup[0]
	}
	ctx.Bot.Send(msg)
}

//Sends photo upload
func (ctx *UpdateContext) SendPhoto(path string) {
	ctx.Bot.Send(tgbotapi.NewPhotoUpload(ctx.ChatId, path))
}

//Answer callback query in the current context
func (ctx *UpdateContext) AnswerCallbackQuery(text string) {
	if ctx.Update.CallbackQuery != nil {
		_, _ = ctx.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(ctx.Update.CallbackQuery.ID, text))
	}
}

//Triggers specified scenario for current user
func (ctx *UpdateContext) TriggerScenario(name string) {
	ctx.TgBot.TriggerScenario(ctx, name)
}