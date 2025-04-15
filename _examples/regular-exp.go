package _examples

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	router "github.com/LZTD1/telebot-context-router"
	"gopkg.in/telebot.v4"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable not set")
	}

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

	// --- Pre-compile Regular Expressions ---
	// This regex matches "/user " followed by one or more digits, capturing the digits.
	userCommandRegex := regexp.MustCompile(`^/user (\d+)$`)
	// This regex matches "view_item:" followed by one or more characters, capturing the part after the colon.
	viewItemCallbackRegex := regexp.MustCompile(`^\fview_item:(.+)$`)

	// --- Register Regex Text Handler ---
	r.HandleFuncRegexpText(userCommandRegex, func(c telebot.Context) error {
		log.Printf("Handler: Received regex match for user command: %q", c.Text())

		matches := userCommandRegex.FindStringSubmatch(c.Text())

		userID := "unknown"
		if len(matches) > 1 {
			userID = matches[1]
		}

		return c.Send(fmt.Sprintf("Processing request for user ID: %s", userID))
	})

	// --- Register Regex Callback Handler ---
	r.HandleFuncRegexpCallback(viewItemCallbackRegex, func(c telebot.Context) error {
		callbackData := c.Callback().Data
		log.Printf("Handler: Received regex match for view item callback: %q", callbackData)

		err := c.Respond(&telebot.CallbackResponse{Text: "Fetching item details..."})
		if err != nil {
			log.Printf("Error responding to callback %q: %v", callbackData, err)
		}

		matches := viewItemCallbackRegex.FindStringSubmatch(callbackData)
		itemID := "unknown"
		if len(matches) > 1 {
			itemID = matches[1]
		}

		return c.Send(fmt.Sprintf("Displaying details for item: %s", itemID))
	})

	r.HandleFuncText("/start", func(c telebot.Context) error {
		log.Printf("Handler: Received /start from %s", c.Sender().Username)
		kbd := &telebot.ReplyMarkup{}
		kbd.Row(
			kbd.Data("View Item ABC", "view_item:ABC"),
			kbd.Data("View Item 123", "view_item:123"),
		)

		return c.Send(
			"Welcome! Try commands like `/user 123` or press a button below.",
			kbd,
			telebot.ModeMarkdown,
		)
	})

	// --- Connect Router to Telebot ---
	bot.Handle(telebot.OnText, r.ServeContext)
	bot.Handle(telebot.OnCallback, r.ServeContext)

	// --- Start the Bot ---
	log.Println("Bot starting...")
	bot.Start()
}
