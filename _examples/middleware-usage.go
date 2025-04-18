package _examples

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	router "github.com/LZTD1/telebot-context-router"
	"gopkg.in/telebot.v4"
)

var (
	ADMIN_IDS = map[int64]interface{}{
		111: "",
		222: "",
	}
)

// --- Middleware Definitions ---

// RandomAccessMw Provides random access to the command
func RandomAccessMw(next router.RouteHandler) router.RouteHandler {
	return router.HandlerFunc(func(c telebot.Context) error {
		if rand.Intn(100) > 50 {
			return next.ServeContext(c)
		}
		log.Printf("[RandomAccessMw] access denined")
		return nil
	})
}

// LoggerMw logs user requests
func LoggerMw(next router.RouteHandler) router.RouteHandler {
	return router.HandlerFunc(func(c telebot.Context) error {
		log.Printf("[LoggerMw] recivied context from %s", c.Sender().Username)
		return next.ServeContext(c)
	})
}

// AdminFilterMw filters admin
func AdminFilterMw(next router.RouteHandler) router.RouteHandler {
	return router.HandlerFunc(func(c telebot.Context) error {
		if _, ok := ADMIN_IDS[c.Sender().ID]; ok {
			return next.ServeContext(c)
		}

		log.Printf("[AdminFilterMw] access denined")
		return nil
	})
}

func main() {
	botToken := "INSERT_TOKEN"

	// --- Bot Setup ---
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// --- Router Setup ---
	r := router.NewRouter()

	// --- Apply Global Middleware using r.Use() ---
	r.Use(LoggerMw)

	// --- Public Routes (No specific role required) ---
	r.HandleFuncText("/start", func(c telebot.Context) error {
		return c.Send(fmt.Sprintf("Welcome!\nTry:\n/rnd - get a random chance to access data\n/ap - admin panel"))
	})

	// -- Route with middleware
	r.With(RandomAccessMw).HandleFuncText("/rnd", func(ctx telebot.Context) error {
		log.Printf("user %s accessed to /rnd endpoint", ctx.Sender().Username)
		return ctx.Send("Happy hacking!")
	})

	// --- Group for admin Routes ---
	r.Group(func(r router.Router) {
		// Apply middleware specific to this group
		r.Use(AdminFilterMw)

		r.HandleFuncText("/ap", func(c telebot.Context) error {
			return c.Send("You can use this commands:\n\n/ap - admin panel\n/restart - restart server")
		})

		r.HandleFuncText("/restart", func(c telebot.Context) error {
			return c.Send("Server restarting...")
		})
	})

	// --- Register a NotFound Handler ---
	r.NotFound(func(c telebot.Context) error {
		log.Printf("NotFound: No route matched for input from %d", c.Sender().ID)
		return c.Send("Unknown command.")
	})

	// --- Connect Router to Telebot ---
	bot.Handle(telebot.OnText, r.ServeContext)

	// --- Start the Bot ---
	log.Println("Bot starting...")
	bot.Start()
}
