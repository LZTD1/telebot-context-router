package _examples

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	router "github.com/LZTD1/telebot-context-router"
	"gopkg.in/telebot.v4"
)

// --- Mock User Roles (Replace with your actual role logic) ---
var userRoles = map[int64]string{
	111111: "admin",
	222222: "moderator",
	333333: "user",
}

func getUserRole(userID int64) string {
	role, ok := userRoles[userID]
	if !ok {
		return "guest"
	}
	return role
}

// --- Middleware Definitions ---

// LoggerMiddleware logs basic information about incoming updates.
func LoggerMiddleware(next router.RouteHandler) router.RouteHandler {
	return router.HandlerFunc(func(c telebot.Context) error {
		start := time.Now()
		userId := c.Sender().ID
		log.Printf("[Log] --> Update %d received from User %d (%s)", c.Update().ID, userId, getUserRole(userId))

		err := next.ServeContext(c)

		log.Printf("[Log] <-- Update %d processed in %v. Error: %v", c.Update().ID, time.Since(start), err)
		return err
	})
}

// RequireRoleMiddleware checks if the user has at least the required role.
func RequireRoleMiddleware(requiredRole string) func(next router.RouteHandler) router.RouteHandler {
	roleHierarchy := map[string]int{
		"guest":     0,
		"user":      1,
		"moderator": 2,
		"admin":     3,
	}

	requiredLevel, ok := roleHierarchy[requiredRole]
	if !ok {
		log.Fatalf("FATAL: Invalid role specified in RequireRoleMiddleware: %s", requiredRole)
	}

	return func(next router.RouteHandler) router.RouteHandler {
		return router.HandlerFunc(func(c telebot.Context) error {
			userRole := getUserRole(c.Sender().ID)
			userLevel := roleHierarchy[userRole]

			if userLevel >= requiredLevel {
				log.Printf("[%s Auth] Access granted for user %d (Role: %s >= Required: %s)", requiredRole, c.Sender().ID, userRole, requiredRole)
				return next.ServeContext(c)
			}

			log.Printf("[%s Auth] Access DENIED for user %d (Role: %s < Required: %s)", requiredRole, c.Sender().ID, userRole, requiredRole)
			_ = c.Send(fmt.Sprintf("Access Denied. Required role: %s", requiredRole))
			return nil
		})
	}
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable not set")
	}

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
	r.Use(LoggerMiddleware)

	// --- Public Routes (No specific role required) ---
	r.HandleFuncText("/start", func(c telebot.Context) error {
		role := getUserRole(c.Sender().ID)
		return c.Send(fmt.Sprintf("Welcome! Your detected role is: %s", role))
	})
	r.HandleFuncText("/help", func(c telebot.Context) error {
		return c.Send("Public help message...")
	})

	// --- Group for Moderator Routes ---
	r.Group(func(modRouter router.Router) {
		// Apply middleware specific to this group
		modRouter.Use(RequireRoleMiddleware("moderator"))

		modRouter.HandleFuncText("/warn", func(c telebot.Context) error {
			return c.Send("Moderator command: Warn user...")
		})

		modRouter.HandleFuncText("/mute", func(c telebot.Context) error {
			return c.Send("Moderator command: Mute user...")
		})
	})

	// --- Group for Admin Routes ---
	r.Group(func(adminRouter router.Router) {
		adminRouter.Use(RequireRoleMiddleware("admin"))

		adminRouter.HandleFuncText("/ban", func(c telebot.Context) error {
			return c.Send("Admin command: Ban user...")
		})

		adminRouter.HandleFuncText("/config", func(c telebot.Context) error {
			return c.Send("Admin command: Show configuration...")
		})
	})

	// --- Route with Inline Middleware using r.With() ---
	r.With(RequireRoleMiddleware("user")).HandleFuncText("/profile", func(c telebot.Context) error {
		userIDstr := strconv.FormatInt(c.Sender().ID, 10)
		return c.Send("Showing profile for user " + userIDstr)
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
