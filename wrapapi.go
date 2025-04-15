package router

import "gopkg.in/telebot.v4"

// wrappedBot is a thin wrapper around telebot.API that intercepts all outgoing actions
// (like sending, editing, replying, etc.) and invokes the markHandled callback before
// performing the actual action. This is used internally by the router to track whether
// a context has already been handled, preventing fallback to the NotFound handler.
type wrappedBot struct {
	telebot.API
	markHandled func()
}

func (b *wrappedBot) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.Send(to, what, opts...)
}

func (b *wrappedBot) SendAlbum(to telebot.Recipient, a telebot.Album, opts ...interface{}) ([]telebot.Message, error) {
	b.markHandled()
	return b.API.SendAlbum(to, a, opts...)
}

func (b *wrappedBot) SendPaid(to telebot.Recipient, stars int, a telebot.PaidAlbum, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.SendPaid(to, stars, a, opts...)
}

func (b *wrappedBot) Reply(to *telebot.Message, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.Reply(to, what, opts...)
}

func (b *wrappedBot) Edit(msg telebot.Editable, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.Edit(msg, what, opts...)
}

func (b *wrappedBot) EditCaption(msg telebot.Editable, caption string, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.EditCaption(msg, caption, opts...)
}

func (b *wrappedBot) EditMedia(msg telebot.Editable, media telebot.Inputtable, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.EditMedia(msg, media, opts...)
}

func (b *wrappedBot) EditReplyMarkup(msg telebot.Editable, markup *telebot.ReplyMarkup) (*telebot.Message, error) {
	b.markHandled()
	return b.API.EditReplyMarkup(msg, markup)
}

func (b *wrappedBot) Delete(msg telebot.Editable) error {
	b.markHandled()
	return b.API.Delete(msg)
}

func (b *wrappedBot) Notify(to telebot.Recipient, action telebot.ChatAction, threadID ...int) error {
	b.markHandled()
	return b.API.Notify(to, action, threadID...)
}

func (b *wrappedBot) Respond(c *telebot.Callback, resp ...*telebot.CallbackResponse) error {
	b.markHandled()
	return b.API.Respond(c, resp...)
}

func (b *wrappedBot) Answer(q *telebot.Query, r *telebot.QueryResponse) error {
	b.markHandled()
	return b.API.Answer(q, r)
}

func (b *wrappedBot) Ship(q *telebot.ShippingQuery, what ...interface{}) error {
	b.markHandled()
	return b.API.Ship(q, what...)
}

func (b *wrappedBot) Accept(q *telebot.PreCheckoutQuery, errorMessage ...string) error {
	b.markHandled()
	return b.API.Accept(q, errorMessage...)
}

func (b *wrappedBot) Forward(to telebot.Recipient, msg telebot.Editable, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.Forward(to, msg, opts...)
}

func (b *wrappedBot) ForwardMany(to telebot.Recipient, msgs []telebot.Editable, opts ...*telebot.SendOptions) ([]telebot.Message, error) {
	b.markHandled()
	return b.API.ForwardMany(to, msgs, opts...)
}

func (b *wrappedBot) Copy(to telebot.Recipient, msg telebot.Editable, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.Copy(to, msg, opts...)
}

func (b *wrappedBot) CopyMany(to telebot.Recipient, msgs []telebot.Editable, opts ...*telebot.SendOptions) ([]telebot.Message, error) {
	b.markHandled()
	return b.API.CopyMany(to, msgs, opts...)
}

func (b *wrappedBot) React(to telebot.Recipient, msg telebot.Editable, r telebot.Reactions) error {
	b.markHandled()
	return b.API.React(to, msg, r)
}

func (b *wrappedBot) Pin(msg telebot.Editable, opts ...interface{}) error {
	b.markHandled()
	return b.API.Pin(msg, opts...)
}

func (b *wrappedBot) Unpin(chat telebot.Recipient, messageID ...int) error {
	b.markHandled()
	return b.API.Unpin(chat, messageID...)
}

func (b *wrappedBot) UnpinAll(chat telebot.Recipient) error {
	b.markHandled()
	return b.API.UnpinAll(chat)
}

func (b *wrappedBot) StopPoll(msg telebot.Editable, opts ...interface{}) (*telebot.Poll, error) {
	b.markHandled()
	return b.API.StopPoll(msg, opts...)
}

func (b *wrappedBot) StopLiveLocation(msg telebot.Editable, opts ...interface{}) (*telebot.Message, error) {
	b.markHandled()
	return b.API.StopLiveLocation(msg, opts...)
}
