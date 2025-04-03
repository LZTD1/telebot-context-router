package router

import (
	"gopkg.in/telebot.v4"
	"regexp"
)

// TypeHandling defines the type of incoming update to handle.
type TypeHandling int

const (
	// CallbackHandle indicates a callback query update.
	CallbackHandle TypeHandling = iota
	// TextHandle indicates a text message update.
	TextHandle
)

// RouteHandler defines the interface for handlers.
type RouteHandler interface {
	// ServeContext processes the incoming telebot context.
	ServeContext(ctx telebot.Context) error
}

// Router describes the main interface for the telebot router.
type Router interface {
	// Use appends middleware to the router stack.
	Use(middlewares ...func(RouteHandler) RouteHandler)
	// With adds inline middleware for subsequent handlers.
	With(middlewares ...func(RouteHandler) RouteHandler) Router
	// Group creates a new router instance for route grouping.
	Group(fn func(r Router)) Router

	// Handle registers a handler for an exact pattern match.
	Handle(pattern string, h RouteHandler, t TypeHandling)
	// HandleFunc registers a handler function for an exact pattern match.
	HandleFunc(pattern string, fn telebot.HandlerFunc, t TypeHandling)

	// HandleText registers a handler for an exact text message match.
	HandleText(pattern string, h RouteHandler)
	// HandleFuncText registers a handler function for an exact text message match.
	HandleFuncText(pattern string, fn telebot.HandlerFunc)

	// HandleCallback registers a handler for an exact callback data match.
	HandleCallback(pattern string, h RouteHandler)
	// HandleFuncCallback registers a handler function for an exact callback data match.
	HandleFuncCallback(pattern string, fn telebot.HandlerFunc)

	// HandleRegexp registers a handler for a regular expression pattern match.
	HandleRegexp(pattern *regexp.Regexp, h RouteHandler, t TypeHandling)
	// HandleFuncRegexp registers a handler function for a regular expression pattern match.
	HandleFuncRegexp(pattern *regexp.Regexp, fn telebot.HandlerFunc, t TypeHandling)

	// HandleRegexpText registers a handler for a regex text message match.
	HandleRegexpText(pattern *regexp.Regexp, h RouteHandler)
	// HandleFuncRegexpText registers a handler function for a regex text message match.
	HandleFuncRegexpText(pattern *regexp.Regexp, fn telebot.HandlerFunc)

	// HandleRegexpCallback registers a handler for a regex callback data match.
	HandleRegexpCallback(pattern *regexp.Regexp, h RouteHandler)
	// HandleFuncRegexpCallback registers a handler function for a regex callback data match.
	HandleFuncRegexpCallback(pattern *regexp.Regexp, fn telebot.HandlerFunc)

	// NotFound sets the handler for routes not found.
	NotFound(h telebot.HandlerFunc)

	// ServeContext processes an incoming telebot update.
	ServeContext(ctx telebot.Context) error
}
