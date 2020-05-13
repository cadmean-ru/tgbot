package context

type TgBot interface {
	TriggerScenario(ctx *UpdateContext, name string)
}
