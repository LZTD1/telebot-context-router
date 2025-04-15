package router

import (
	"gopkg.in/telebot.v4"
)

// wrappedContext embeds the original telebot.Context and overrides all methods
// that perform an action (e.g., sending, replying, editing, etc.) to mark the
// context as "handled". This is used internally by the router to track whether
// a handler has responded to the update, allowing the router to avoid calling
// the NotFound handler unnecessarily.
type wrappedContext struct {
	telebot.Context
	api     *wrappedBot
	handled bool
}

func (w *wrappedContext) markHandled() {
	w.handled = true
}

func (w *wrappedContext) WasHandled() bool {
	return w.handled
}

func (w *wrappedContext) Send(what interface{}, opts ...interface{}) error {
	w.markHandled()
	return w.Context.Send(what, opts...)
}

func (w *wrappedContext) Bot() telebot.API {
	return w.api
}

func (w *wrappedContext) SendAlbum(a telebot.Album, opts ...interface{}) error {
	w.markHandled()
	return w.Context.SendAlbum(a, opts...)
}

func (w *wrappedContext) Reply(what interface{}, opts ...interface{}) error {
	w.markHandled()
	return w.Context.Reply(what, opts...)
}

func (w *wrappedContext) Forward(msg telebot.Editable, opts ...interface{}) error {
	w.markHandled()
	return w.Context.Forward(msg, opts...)
}

func (w *wrappedContext) ForwardTo(to telebot.Recipient, opts ...interface{}) error {
	w.markHandled()
	return w.Context.ForwardTo(to, opts...)
}

func (w *wrappedContext) Edit(what interface{}, opts ...interface{}) error {
	w.markHandled()
	return w.Context.Edit(what, opts...)
}

func (w *wrappedContext) EditCaption(caption string, opts ...interface{}) error {
	w.markHandled()
	return w.Context.EditCaption(caption, opts...)
}

func (w *wrappedContext) EditOrSend(what interface{}, opts ...interface{}) error {
	w.markHandled()
	return w.Context.EditOrSend(what, opts...)
}

func (w *wrappedContext) EditOrReply(what interface{}, opts ...interface{}) error {
	w.markHandled()
	return w.Context.EditOrReply(what, opts...)
}

func (w *wrappedContext) Delete() error {
	w.markHandled()
	return w.Context.Delete()
}

func (w *wrappedContext) Respond(resp ...*telebot.CallbackResponse) error {
	w.markHandled()
	return w.Context.Respond(resp...)
}

func (w *wrappedContext) RespondText(text string) error {
	w.markHandled()
	return w.Context.RespondText(text)
}

func (w *wrappedContext) RespondAlert(text string) error {
	w.markHandled()
	return w.Context.RespondAlert(text)
}

func (w *wrappedContext) Answer(resp *telebot.QueryResponse) error {
	w.markHandled()
	return w.Context.Answer(resp)
}

func (w *wrappedContext) Accept(errorMessage ...string) error {
	w.markHandled()
	return w.Context.Accept(errorMessage...)
}

func (w *wrappedContext) Ship(what ...interface{}) error {
	w.markHandled()
	return w.Context.Ship(what...)
}

func (w *wrappedContext) Notify(action telebot.ChatAction) error {
	w.markHandled()
	return w.Context.Notify(action)
}
