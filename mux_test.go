package router

import (
	"github.com/stretchr/testify/assert"
	tb "gopkg.in/telebot.v4"
	"regexp"
	"testing"
)

func TestMux(t *testing.T) {
	t.Run("Exact Text Match", func(t *testing.T) {
		mux := NewRouter()
		mux.HandleFuncText("/ping", func(ctx tb.Context) error {
			return ctx.Send("pong")
		})

		ctx := &mockContext{text: "/ping"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "pong")
	})

	t.Run("Regex Text Match", func(t *testing.T) {
		mux := NewRouter()
		pattern := regexp.MustCompile(`^/user \d+$`)
		mux.HandleFuncRegexpText(pattern, func(ctx tb.Context) error {
			return ctx.Send("user matched")
		})

		ctx := &mockContext{text: "/user 42"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "user matched")
	})

	t.Run("Exact Callback Match", func(t *testing.T) {
		mux := NewRouter()
		mux.HandleFuncCallback("button_click", func(ctx tb.Context) error {
			return ctx.Send("clicked")
		})

		ctx := &mockContext{callback: "button_click"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "clicked")
	})

	t.Run("Regex Callback Match", func(t *testing.T) {
		mux := NewRouter()
		pattern := regexp.MustCompile(`^action:[a-z]+$`)
		mux.HandleFuncRegexpCallback(pattern, func(ctx tb.Context) error {
			return ctx.Send("regex callback")
		})

		ctx := &mockContext{callback: "action:edit"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "regex callback")
	})

	t.Run("NotFound Called", func(t *testing.T) {
		mux := NewRouter()
		mux.NotFound(func(ctx tb.Context) error {
			return ctx.Send("not found")
		})

		ctx := &mockContext{text: "/unknown"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "not found")
	})

	t.Run("Middleware Applied", func(t *testing.T) {
		mux := NewRouter()
		var steps []string

		mux.Use(func(next RouteHandler) RouteHandler {
			return HandlerFunc(func(ctx tb.Context) error {
				steps = append(steps, "mw")
				return next.ServeContext(ctx)
			})
		})

		mux.HandleFuncText("/hello", func(ctx tb.Context) error {
			steps = append(steps, "handler")
			return ctx.Send("ok")
		})

		ctx := &mockContext{text: "/hello"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"mw", "handler"}, steps)
		assert.Contains(t, ctx.sent, "ok")
	})

	t.Run("Wrapped Context: wasHandled via Send", func(t *testing.T) {
		mux := NewRouter()
		mux.HandleFuncText("/mark", func(ctx tb.Context) error {
			return ctx.Send("marked")
		})

		ctx := &mockContext{text: "/mark"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.True(t, ctx.wasHandled)
		assert.Contains(t, ctx.sent, "marked")
	})

	t.Run("Not access to mw and return notFound", func(t *testing.T) {
		mux := NewRouter()

		mux.With(func(routeHandler RouteHandler) RouteHandler {
			return HandlerFunc(func(c tb.Context) error {
				return nil
			})
		}).HandleFuncText("/ghost", func(ctx tb.Context) error {
			return ctx.Send("")
		})

		mux.NotFound(func(ctx tb.Context) error {
			return ctx.Send("not found from ghost")
		})

		ctx := &mockContext{text: "/ghost"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Contains(t, ctx.sent, "not found from ghost")
	})

	t.Run("Group Middleware Inheritance", func(t *testing.T) {
		mux := NewRouter()
		var trace []string

		mux.Use(func(next RouteHandler) RouteHandler {
			return HandlerFunc(func(ctx tb.Context) error {
				trace = append(trace, "root")
				return next.ServeContext(ctx)
			})
		})

		group := mux.With(func(next RouteHandler) RouteHandler {
			return HandlerFunc(func(ctx tb.Context) error {
				trace = append(trace, "group")
				return next.ServeContext(ctx)
			})
		})

		group.HandleFuncText("/greet", func(ctx tb.Context) error {
			trace = append(trace, "handler")
			return ctx.Send("hi")
		})

		ctx := &mockContext{text: "/greet"}
		err := mux.ServeContext(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []string{"root", "group", "handler"}, trace)
		assert.Contains(t, ctx.sent, "hi")
	})
}

type mockContext struct {
	tb.Context
	text       string
	callback   string
	sent       []string
	marked     bool
	wasHandled bool
}

func (m *mockContext) Message() *tb.Message {
	if m.text == "" {
		return nil
	}
	return &tb.Message{Text: m.text}
}

func (m *mockContext) Callback() *tb.Callback {
	if m.callback == "" {
		return nil
	}
	return &tb.Callback{Data: m.callback}
}

func (m *mockContext) Send(what interface{}, _ ...interface{}) error {
	m.sent = append(m.sent, what.(string))
	m.wasHandled = true
	return nil
}

func (m *mockContext) Bot() tb.API {
	return dummyBot{ctx: m}
}

type dummyBot struct {
	tb.API
	ctx *mockContext
}

func (b dummyBot) Send(to tb.Recipient, what interface{}, _ ...interface{}) (*tb.Message, error) {
	b.ctx.marked = true
	return nil, nil
}
