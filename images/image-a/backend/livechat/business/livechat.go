package business

import (
	"os"

	"github.com/gofiber/websocket/v2"
	"go.violettedev.com/eecs4222/livechat/structs"
	"go.violettedev.com/eecs4222/shared/auth"
)

func HandleAuthMessage(c *websocket.Conn, authCtx *structs.AuthContext, content string) {
	authProvider := auth.Initialize(os.Getenv("AUTH_PRIVATE_KEY"),
		os.Getenv("REFRESH_PRIVATE_KEY"))
	name, email, err := authProvider.GetAuthContextWebSocket(c, content)
	if err != nil && err.Error() == "refresh" {
		// Request refresh
		c.WriteJSON(structs.Message{
			Type: "refresh",
		})
	} else {
		authCtx.Name = name
		authCtx.Email = email
		c.WriteJSON(structs.Message{
			Type: "signed_in",
		})
	}
	// Note: We don't send anything on no auth -- as it's default behavior
}

func HandlePing(c *websocket.Conn, content string) {
	c.WriteJSON(structs.Message{
		Type: "pong",
	})
}
