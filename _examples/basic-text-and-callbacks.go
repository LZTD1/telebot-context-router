package _examples

import (
	"fmt"
	"log"
	"time"

	router "github.com/LZTD1/telebot-context-router"
	"gopkg.in/telebot.v4"
)

func main() {
	botToken := "INSERT_TOKEN"

	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := router.NewRouter()

	// Handle the "/start" text command
	r.HandleFuncText("/start", func(c telebot.Context) error {
		log.Printf("Handler: Received /start from %s", c.Sender().Username)

		kbd := &telebot.ReplyMarkup{}

		kbd.Inline(
			kbd.Row(
				kbd.URL("Visit Repo", "https://github.com/LZTD1/telebot-context-router"),
				kbd.Data("Show Help", "show_help_callback"),
			),
		)

		return c.Send(
			fmt.Sprintf("Hello, %s!\n\n Welcome to the basic router example.", c.Sender().FirstName),
			kbd,
		)
	})

	// Handle the "Help" text command (case-sensitive exact match)
	r.HandleFuncText("Help", func(c telebot.Context) error {
		log.Printf("Handler: Received 'Help' text from %s", c.Sender().Username)
		return c.Send("Available commands:\n/start - Show welcome message and buttons\nHelp - Show this help message")
	})

	// --- Register Exact Callback Query Handlers ---

	// Handle the callback data "show_help_callback" sent by the button
	r.HandleFuncCallback("\fshow_help_callback", func(c telebot.Context) error {
		log.Printf("Handler: Received 'show_help_callback' callback from %s", c.Sender().Username)

		err := c.Respond(&telebot.CallbackResponse{
			Text: "Okay, showing help!",
		})
		if err != nil {
			log.Printf("Error responding to callback 'show_help_callback': %v", err)
		}

		return c.Send("Available commands:\n/start - Show welcome message and buttons\nHelp - Show this help message")
	})

	r.NotFound(func(ctx telebot.Context) error {
		log.Printf("Recivied unknown command %s", ctx.Message().Text)
		return ctx.Send("Hmm... i dont know this command")
	})

	// --- Connect Router to Telebot ---
	bot.Handle(telebot.OnText, r.ServeContext)
	bot.Handle(telebot.OnCallback, r.ServeContext)

	// --- Start the Bot ---
	log.Println("Bot starting...")
	bot.Start()
}
