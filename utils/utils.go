package utils


import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func GetChatFromId(update *tgbotapi.Update) (int64, int) {
	if update.Message != nil {
		return update.Message.Chat.ID, update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.ID
	} else {
		return 0, 0
	}
}

func GetChatId(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	} else {
		return 0
	}
}

func GtFromId(update tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	} else {
		return 0
	}
}

