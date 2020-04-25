package tgbot

import (
	"github.com/cadmean-ru/tgbot/context"
	"github.com/cadmean-ru/tgbot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Function to handle message or callbacks. 
//If it returns an error error handler will be called.
type UpdateHandler func(ctx *context.UpdateContext) error

//Function to be called when update handler returns sn error.
type ErrorHandler func(ctx *context.UpdateContext, err error)

type Bot struct {
	TgBot          *tgbotapi.BotAPI
	handlers       map[string]UpdateHandler
	defaultHandler UpdateHandler
	scenarios      []*Scenario
	errorHandler   ErrorHandler
	StateProvider  context.StateProvider
}

//Created a new bot with specified key.
//If debug is true bot will output a lot of logs bruuuu.
func NewBot(key string, debug bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	return &Bot{
		TgBot:     bot,
		handlers:  make(map[string]UpdateHandler),
		scenarios: make([]*Scenario, 0),
	}, nil
}

//Created a new bot with specified key. Also sets state provider.
//If debug is true bot will output a lot of logs bruuuu.
func NewBotWithProvider(key string, debug bool, provider context.StateProvider) (*Bot, error) {
	bot, err := NewBot(key, debug)
	if err != nil {
		return nil, err
	}
	bot.StateProvider = provider
	return bot, nil
}

//Register new handler for specified command
func (b *Bot) Handle(command string, handler UpdateHandler) {
	b.handlers[command] = handler
}

//Set the default handler if no command or scenarios match
func (b *Bot) HandleDefault(handler UpdateHandler) {
	b.defaultHandler = handler
}

//Register new scenario
func (b *Bot) HandleScenario(s *Scenario) {
	b.scenarios = append(b.scenarios, s)
}

//Register error handler that will be called if update handler returns an error
func (b *Bot) HandleError(handler ErrorHandler) {
	b.errorHandler = handler
}

//Start receiving updates
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.TgBot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		ctx := b.resolveUpdate(&update)

		if h := b.resolveCommand(ctx.Text); h != nil {
			b.handleError(h(&ctx))
			continue
		}

		if b.StateProvider != nil {
			if s := b.resolveScenario(ctx.Text); s != nil {
				state := context.State{ Scenario: s.Name }
				b.handleScenario(*s, &ctx, &state)
				err := b.StateProvider.Save(ctx, state)
				b.handleError(err)
				continue
			}

			state, err := b.StateProvider.Load(ctx)
			if err != nil {
				b.handleError(err)
				continue
			}

			if state.Scenario != "" && state.Step != "" {
				var s Scenario
				for _, s1 := range b.scenarios {
					if s1.Name == state.Scenario {
						s = *s1
						break
					}
				}
				b.handleScenario(s, &ctx, &state)
				err := b.StateProvider.Save(ctx, state)
				b.handleError(err)
				continue
			}
		}

		b.handleDefault(ctx)
	}

	return nil
}


func (b *Bot) resolveUpdate(update *tgbotapi.Update) context.UpdateContext {
	chatId, fromId := utils.GetChatFromId(update)
	var text string
	var contact *tgbotapi.Contact
	var location *tgbotapi.Location
	if update.Message != nil {
		text = update.Message.Text
		contact = update.Message.Contact
		location = update.Message.Location
	}
	var data string
	if update.CallbackQuery != nil {
		data = update.CallbackQuery.Data
	}

	ctx := context.UpdateContext{
		Update:       update,
		Bot:          b.TgBot,
		ChatId:       chatId,
		FromId:       fromId,
		Text:         text,
		CallbackData: data,
		Contact:      contact,
		Location:     location,
	}

	return ctx
}



func (b *Bot) resolveCommand(command string) UpdateHandler {
	if h, ok := b.handlers[command]; ok {
		return h
	}

	return nil
}

func (b *Bot) resolveScenario(trigger string) *Scenario {
	for _, s := range b.scenarios {
		for _, t := range s.Triggers {
			if t == trigger {
				return s
			}
		}
	}

	return nil
}

func (b *Bot) handleScenario(s Scenario, ctx *context.UpdateContext, state *context.State) {
	var step Step
	if state.Step == "" {
		step = s.Steps[0]
	} else {
		for _, s := range s.Steps {
			if s.Name == state.Step {
				step = s
				break
			}
		}
	}

	if step.Name == "" {
		return
	}

	var next = ""

	next, err := step.Handler(ctx, state)
	if err != nil {
		b.handleError(err)
		return
	}

	state.Step = next
	if state.Step == "" {
		state.Scenario = ""
		state.Data = nil
	}
}


func (b *Bot) handleError(err error) {
	if err != nil && b.errorHandler != nil {
		b.errorHandler(err)
	}
}

func (b *Bot) handleDefault(ctx context.UpdateContext) {
	if b.defaultHandler != nil {
		b.handleError(b.defaultHandler(&ctx))
	}
}
