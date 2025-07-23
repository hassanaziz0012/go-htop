package opts

import "github.com/hassanaziz0012/go-htop/ui/state"

type ButtonHandler struct {
	Label   string
	Handler func(state *state.AppState)
}

func NewButtonHandler(label string) *ButtonHandler {
	return &ButtonHandler{Label: label}
}

func (b *ButtonHandler) SetHandler(handler func(state *state.AppState)) *ButtonHandler {
	b.Handler = handler
	return b
}
