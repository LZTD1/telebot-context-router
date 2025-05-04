# Telebot Context Router
[![Go Reference](https://pkg.go.dev/badge/github.com/LZTD1/telebot-context-router.svg)](https://pkg.go.dev/github.com/LZTD1/telebot-context-router)
![ChatGPT Image 3 апр  2025 г , 20_53_40](https://github.com/user-attachments/assets/0ec2d68b-6291-4f78-840e-43b00e46ed16)

A flexible router for [telebot v4](https://github.com/tucnak/telebot) inspired by `go-chi`. This router provides easy routing for text messages and callback queries in your Telegram bots, supporting middleware and route grouping.

## Installation

```bash
go get github.com/LZTD1/telebot-context-router@v1.3.0
```
## How It Works (Core Principle)

This router simplifies handling Telegram updates by providing two main ways to match incoming messages or callback queries:

1.  **Exact Match:** You can define handlers that trigger only when the incoming text or callback data *exactly* matches a specific string you provide (e.g., the command `/start` or the callback data `confirm_order`). This is very efficient for predefined commands and button actions.
2.  **Pattern Match (Regular Expressions):** For more complex scenarios, you can define handlers using Go's regular expressions. This allows you to match commands with arguments (like `/user 123`), callback data with variable parts (`item_view_*`), or any text conforming to a specific pattern.

The router first checks for an exact match. If none is found, it then checks the input against your registered regular expression patterns one by one until a match occurs.

Additionally, the router supports **Middleware**. Think of middleware as processing steps that can run *before* your main handler logic executes.

To ensure proper handling, the router wraps the `telebot.Context` to track if the context has already been processed (e.g., a message has been sent or edited). This prevents fallback to the "not found" handler when the context has been handled earlier in the routing process.

## Simple Example Usage

This example shows the most basic usage for handling a simple text command.

```go
func main() {
	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, _ := telebot.NewBot(pref)

	// Initialize the router
	r := router.NewRouter()

	// --- Register a simple text command handler ---
	r.HandleFuncText("/start", func(c telebot.Context) error {
		log.Printf("Handling /start command from %s", c.Sender().Username)
		return c.Send("Hello there!")
	})

	// --- Set a handler for unmatched routes ---
	r.NotFound(func(c telebot.Context) error {
		log.Printf("Unknown command from %s: %q", c.Sender().Username, c.Text())
		return c.Send("Sorry, I didn't understand that command.")
	})

	// --- Register the router's ServeContext as the main handler for Telebot ---
	bot.Handle(telebot.OnText, r.ServeContext)

	log.Println("Bot starting...")
	bot.Start()
}
```

## Examples

For more detailed examples covering specific features, please see the _examples directory:

- [Basic Text and Callbacks: Handling simple commands and button presses.](./_examples/basic-text-and-callbacks.go)
- [Using Regular Expressions: Matching commands or callback data with patterns (e.g., /user_(\d+), view_item:(.*)).](./_examples/regular-exp.go)
- [Middleware Usage: Applying global middleware (Use) and scoped middleware (Group, With) for logging, authentication, etc.](./_examples/middleware-usage.go)

## UPDATES
- **v1.3.0**: Context now ensures that NotFound is called unless an exact match marks the context as processed
- **v1.2.1**: Changing `_examples` to be more concise
- **v1.2.0**: Added context wrapping functionality to track whether the context has been processed.

## License
Licensed under [MIT License](./LICENSE)
