package router

import (
	"errors"
	"gopkg.in/telebot.v4"
	"regexp"
)

var (
	// ErrNotFound is returned when no matching route is found.
	ErrNotFound = errors.New("router: not found")
)

// regexEntry holds a compiled regular expression and its associated handler.
// Used for routing based on regex patterns.
type regexEntry struct {
	regex   *regexp.Regexp
	handler RouteHandler
}

// Mux implements the Router interface. It matches incoming telebot updates
// against registered routes and executes the corresponding handler.
// It supports exact string matching (using maps for O(1) lookup) and
// regular expression matching (using slices for O(N) lookup).
// Middleware can be applied globally or scoped using groups.
type Mux struct {
	parent              *Mux
	middlewares         []func(RouteHandler) RouteHandler
	exactTextRoutes     map[string]RouteHandler
	exactCallbackRoutes map[string]RouteHandler
	regexTextRoutes     []regexEntry
	regexCallbackRoutes []regexEntry
	notFoundHandler     telebot.HandlerFunc
}

// NewRouter returns a new, initialized Mux ready to configure.
// It initializes internal maps and slices to avoid nil pointers.
func NewRouter() *Mux {
	return &Mux{
		exactTextRoutes:     make(map[string]RouteHandler),
		exactCallbackRoutes: make(map[string]RouteHandler),
		regexTextRoutes:     make([]regexEntry, 0),
		regexCallbackRoutes: make([]regexEntry, 0),
	}
}

// collectMiddlewares walks up the Mux hierarchy (via parent pointers)
// and gathers all middleware functions, starting from the root Mux down
// to the current one. This ensures middlewares are applied in the correct
// order (root first).
func (m *Mux) collectMiddlewares() []func(RouteHandler) RouteHandler {
	var middlewares []func(RouteHandler) RouteHandler
	current := m
	for current != nil {
		middlewares = append(current.middlewares, middlewares...)
		current = current.parent
	}
	return middlewares
}

// chain builds a middleware chain by wrapping the final endpoint RouteHandler
// with the provided middleware functions. Execution happens in reverse order
// of the slice (onion-style).
func chain(middlewares []func(RouteHandler) RouteHandler, endpoint RouteHandler) RouteHandler {
	innerChain := endpoint
	for i := len(middlewares) - 1; i >= 0; i-- {
		innerChain = middlewares[i](innerChain)
	}
	return innerChain
}

// Handle registers a handler for an exact match of the pattern string.
// It applies the middleware stack collected from the Mux hierarchy to the handler
// before storing it. If the Mux is part of a group, the route is also copied
// to the parent Mux's corresponding map.
func (m *Mux) Handle(pattern string, h RouteHandler, t TypeHandling) {
	allMiddlewares := m.collectMiddlewares()
	finalHandler := chain(allMiddlewares, h)

	switch t {
	case TextHandle:
		m.exactTextRoutes[pattern] = finalHandler
		if p := m.parent; p != nil {
			p.exactTextRoutes[pattern] = finalHandler
		}
	case CallbackHandle:
		m.exactCallbackRoutes[pattern] = finalHandler
		if p := m.parent; p != nil {
			p.exactCallbackRoutes[pattern] = finalHandler
		}
	}
}

// HandleFunc is a convenience method for registering a telebot.HandlerFunc
// for an exact match of the pattern string. It adapts the function to the
// RouteHandler interface.
func (m *Mux) HandleFunc(pattern string, fn telebot.HandlerFunc, t TypeHandling) {
	m.Handle(pattern, HandlerFunc(fn), t)
}

// HandleRegexp registers a handler for a pattern defined by a compiled regular expression.
// It applies the middleware stack collected from the Mux hierarchy to the handler.
// If the Mux is part of a group, the route is also copied to the parent Mux's
// corresponding slice. The pattern must not be nil.
func (m *Mux) HandleRegexp(pattern *regexp.Regexp, h RouteHandler, t TypeHandling) {
	if pattern == nil {
		panic("router: HandleRegexp called with nil pattern")
	}
	allMiddlewares := m.collectMiddlewares()
	finalHandler := chain(allMiddlewares, h)

	entry := regexEntry{
		regex:   pattern,
		handler: finalHandler,
	}

	switch t {
	case TextHandle:
		m.regexTextRoutes = append(m.regexTextRoutes, entry)
		if p := m.parent; p != nil {
			p.regexTextRoutes = append(p.regexTextRoutes, entry)
		}
	case CallbackHandle:
		m.regexCallbackRoutes = append(m.regexCallbackRoutes, entry)
		if p := m.parent; p != nil {
			p.regexCallbackRoutes = append(p.regexCallbackRoutes, entry)
		}
	}
}

// HandleFuncRegexp is a convenience method for registering a telebot.HandlerFunc
// for a pattern defined by a compiled regular expression. It adapts the function
// to the RouteHandler interface.
func (m *Mux) HandleFuncRegexp(pattern *regexp.Regexp, fn telebot.HandlerFunc, t TypeHandling) {
	m.HandleRegexp(pattern, HandlerFunc(fn), t)
}

// HandleText is a convenience method for Handle with TypeHandling set to TextHandle.
// Registers a handler for an exact text message match.
func (m *Mux) HandleText(pattern string, h RouteHandler) { m.Handle(pattern, h, TextHandle) }

// HandleFuncText is a convenience method for HandleFunc with TypeHandling set to TextHandle.
// Registers a handler function for an exact text message match.
func (m *Mux) HandleFuncText(pattern string, fn telebot.HandlerFunc) {
	m.HandleFunc(pattern, fn, TextHandle)
}

// HandleCallback is a convenience method for Handle with TypeHandling set to CallbackHandle.
// Registers a handler for an exact callback data match.
func (m *Mux) HandleCallback(pattern string, h RouteHandler) { m.Handle(pattern, h, CallbackHandle) }

// HandleFuncCallback is a convenience method for HandleFunc with TypeHandling set to CallbackHandle.
// Registers a handler function for an exact callback data match.
func (m *Mux) HandleFuncCallback(pattern string, fn telebot.HandlerFunc) {
	m.HandleFunc(pattern, fn, CallbackHandle)
}

// HandleRegexpText is a convenience method for HandleRegexp with TypeHandling set to TextHandle.
// Registers a handler for a text message match based on a regular expression.
func (m *Mux) HandleRegexpText(pattern *regexp.Regexp, h RouteHandler) {
	m.HandleRegexp(pattern, h, TextHandle)
}

// HandleFuncRegexpText is a convenience method for HandleFuncRegexp with TypeHandling set to TextHandle.
// Registers a handler function for a text message match based on a regular expression.
func (m *Mux) HandleFuncRegexpText(pattern *regexp.Regexp, fn telebot.HandlerFunc) {
	m.HandleFuncRegexp(pattern, fn, TextHandle)
}

// HandleRegexpCallback is a convenience method for HandleRegexp with TypeHandling set to CallbackHandle.
// Registers a handler for a callback data match based on a regular expression.
func (m *Mux) HandleRegexpCallback(pattern *regexp.Regexp, h RouteHandler) {
	m.HandleRegexp(pattern, h, CallbackHandle)
}

// HandleFuncRegexpCallback is a convenience method for HandleFuncRegexp with TypeHandling set to CallbackHandle.
// Registers a handler function for a callback data match based on a regular expression.
func (m *Mux) HandleFuncRegexpCallback(pattern *regexp.Regexp, fn telebot.HandlerFunc) {
	m.HandleFuncRegexp(pattern, fn, CallbackHandle)
}

// Use adds one or more middleware handlers to the Mux's middleware stack.
// Middleware added via Use are applied before middleware added via With or Group
// during the handler chaining process in Handle/HandleRegexp.
func (m *Mux) Use(middlewares ...func(RouteHandler) RouteHandler) {
	m.middlewares = append(m.middlewares, middlewares...)
}

// With creates a new Mux instance configured as a sub-router (inline group).
// It inherits the parent's NotFound handler and gains a pointer to the parent,
// allowing it to collect the parent's middleware when its own Handle/HandleRegexp
// methods are called. Middlewares passed to With are added to the new Mux's stack.
func (m *Mux) With(middlewares ...func(RouteHandler) RouteHandler) Router {
	nm := &Mux{
		parent:              m,
		middlewares:         middlewares,
		exactTextRoutes:     make(map[string]RouteHandler),
		exactCallbackRoutes: make(map[string]RouteHandler),
		regexTextRoutes:     make([]regexEntry, 0),
		regexCallbackRoutes: make([]regexEntry, 0),
		notFoundHandler:     m.notFoundHandler,
	}
	return nm
}

// Group creates a new Mux sub-router (inline group) similar to With.
// It executes the provided function `fn` with the new sub-router, allowing
// for convenient route definition within the group's scope. Middlewares applied
// within the group will be collected when registering handlers inside `fn`.
func (m *Mux) Group(fn func(r Router)) Router {
	im := m.With()
	if fn != nil {
		fn(im)
	}
	return im
}

// NotFound sets the handler function to be called when no route matches.
// The handler is stored on the current Mux instance.
func (m *Mux) NotFound(h telebot.HandlerFunc) {
	m.notFoundHandler = h
}

// findNotFoundHandler searches for a configured NotFound handler by walking
// up the Mux hierarchy via the parent pointer. If no custom handler is found,
// it returns a default handler that simply returns ErrNotFound.
func (m *Mux) findNotFoundHandler() telebot.HandlerFunc {
	current := m
	for current != nil {
		if current.notFoundHandler != nil {
			return current.notFoundHandler
		}
		current = current.parent
	}
	return func(ctx telebot.Context) error {
		return ErrNotFound
	}
}

// NotFoundHandler returns the appropriate NotFound handler for this Mux,
// searching up the hierarchy if necessary.
func (m *Mux) NotFoundHandler() telebot.HandlerFunc {
	return m.findNotFoundHandler()
}

// ServeContext is the main entry point for processing telebot updates.
// It determines the type of update (Text or Callback), finds a matching
// handler (first checking exact matches, then regular expressions), executes
// the handler (which includes the pre-applied middleware chain), and returns
// the result. If no handler is found, it calls the NotFound handler.
func (m *Mux) ServeContext(ctx telebot.Context) error {
	var input string
	var exactMap map[string]RouteHandler
	var regexSlice []regexEntry

	cb := ctx.Callback()
	msg := ctx.Message()

	if cb != nil {
		input = cb.Data
		exactMap = m.exactCallbackRoutes
		regexSlice = m.regexCallbackRoutes
	} else if msg != nil && msg.Text != "" {
		input = msg.Text
		exactMap = m.exactTextRoutes
		regexSlice = m.regexTextRoutes
	} else {
		return m.NotFoundHandler()(ctx)
	}

	if handler, ok := exactMap[input]; ok {
		return handler.ServeContext(ctx)
	}

	isHandeled := false
	for _, entry := range regexSlice {
		if entry.regex.MatchString(input) {
			err := entry.handler.ServeContext(ctx)
			if err != nil {
				return err
			}
			isHandeled = true
		}
	}

	if !isHandeled {
		return m.NotFoundHandler()(ctx)
	}
	return nil
}

// HandlerFunc is an adapter type that allows a regular telebot.HandlerFunc
// to be used as a RouteHandler.
type HandlerFunc telebot.HandlerFunc

// ServeContext implements the RouteHandler interface for HandlerFunc.
func (h HandlerFunc) ServeContext(ctx telebot.Context) error {
	return h(ctx)
}
