package main

import (
	"log"
	"os"

	"go.violettedev.com/eecs4222/shared/auth"
	"go.violettedev.com/eecs4222/shared/database"

	historicalMessagesPresentation "go.violettedev.com/eecs4222/historical_messaging/presentation"
	liveChatPresentation "go.violettedev.com/eecs4222/livechat/presentation"
	loginPresentation "go.violettedev.com/eecs4222/login/presentation"
	registerPresentation "go.violettedev.com/eecs4222/registration/presentation"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

/*
Load dev environment from .env if os.Getenv('GO_ENV') != 'prod'
*/
func loadDevEnv() {
	// Get dev env if not prod
	if os.Getenv("GO_ENV") != "prod" {
		// Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

/*
Exposes the rest endpoints for the various modules to the given app.
*/
func exposeEndpoints(app *fiber.App) {
	// Create auth provider
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	// Define endpoints
	app.Post("/api/login", loginPresentation.LoginEndpoint)
	app.Post("/api/logout", authProvider.IsAuthenticatedFiberMiddleware, loginPresentation.LogoutEndpoint)
	app.Post("/api/register", registerPresentation.RegisterEndpoint)
	app.Post("/api/change_password", authProvider.IsAuthenticatedFiberMiddleware, registerPresentation.ChangePasswordEndpoint)
	app.Get("/api/refresh_credentials", authProvider.IsAuthenticatedFiberMiddleware, loginPresentation.RefreshEndpoint)
	app.Get("/api/messages", historicalMessagesPresentation.GetMessagesEndpoint)
	// Chat websocket endpoint (populate w/ auth info if possible)
	app.Get("/ws/chat", liveChatPresentation.CanUpgradeToWebSocket,
		websocket.New(liveChatPresentation.LiveChatWebSocket))
	// Health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

}

func main() {
	// Load dev environment if needed
	loadDevEnv()
	// Initialize db connection
	database.InitializeDBConnection(os.Getenv("MONGODB_URL"))
	// Initialize distributed messaging
	liveChatPresentation.InitializeDistributedMessaging()
	// Initialize fiber (REST framework)
	app := fiber.New()
	exposeEndpoints(app)
	// Listen
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
