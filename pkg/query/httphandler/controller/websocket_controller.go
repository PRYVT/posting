package controller

import (
	"net/http"

	"github.com/PRYVT/posting/pkg/query/eventhandling"
	ws "github.com/PRYVT/posting/pkg/query/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WSController struct {
	userEventH *eventhandling.PostEventHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWsController(userEventH *eventhandling.PostEventHandler) *WSController {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &WSController{userEventH: userEventH}
}

func (w *WSController) OnRequest(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Warn().Err(err).Msg("Error while upgrading connection")

	} else {
		w.userEventH.AddWebsocketConnection(ws.NewWebsocketConnection(conn))
	}
}
