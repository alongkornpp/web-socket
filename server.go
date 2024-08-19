package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler
func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return err
	}
	defer ws.Close()

	for {
		// Read message from client
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		log.Printf("Received: %s", msg)

		// Write message back to client
		err = ws.WriteMessage(mt, msg)
		if err != nil {
			log.Println("Write Error:", err)
			break
		}
	}
	return nil
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// WebSocket route
	e.GET("/ws", handleWebSocket)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}