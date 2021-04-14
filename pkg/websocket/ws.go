package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Websocket Websocket utils collection
type Websocket struct {
	upgrader *websocket.Upgrader
}

type websocketHandler func(*websocket.Conn) error

const (
	writeWait    = 10 * time.Second
	pongWait     = 30 * time.Second
	pingInterval = pongWait * 9 / 10
)

// NewWebsocket create new Websocket
func NewWebsocket() *Websocket {
	return &Websocket{
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			HandshakeTimeout: 3 * time.Second,
		},
	}
}

// WithHeartbeat wrap handler with heartbeat probe
func (ws Websocket) WithHeartbeat(handler websocketHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := ws.upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// handle background heartbeat
		go heartbeatRoutine(conn)
		// handle business
		go handlerRoutine(conn, handler)
		return nil
	}
}

func heartbeatRoutine(conn *websocket.Conn) {
	ticker := time.NewTicker(pingInterval)
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer func() {
		ticker.Stop()
		conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

func handlerRoutine(conn *websocket.Conn, handler websocketHandler) {
	defer conn.Close()
	for {
		if err := handler(conn); err != nil {
			break
		}
	}
}
