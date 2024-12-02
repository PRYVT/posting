package websocket

import (
	"fmt"

	"github.com/PRYVT/utils/pkg/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebsocketConnection struct {
	connection      *websocket.Conn
	IsConnected     bool
	IsAuthenticated bool
	userUuid        uuid.UUID
}

func NewWebsocketConnection(conn *websocket.Conn) *WebsocketConnection {
	wC := &WebsocketConnection{connection: conn}
	go wC.ReadForDisconnect()
	return wC
}

func (wC *WebsocketConnection) WriteJSON(v interface{}) error {
	if !wC.IsAuthenticated {
		return fmt.Errorf("WebsocketConnection is not connected or authenticated")
	}
	err := wC.connection.WriteJSON(v)
	if err != nil {
		log.Warn().Err(err).Msg("Error while writing WriteJSON")
	}
	return err
}

func (wC *WebsocketConnection) ReadForDisconnect() {
	wC.IsConnected = true
	for {
		authRequest := AuthRequest{}
		err := wC.connection.ReadJSON(&authRequest)
		if err != nil {
			log.Debug().Err(err).Msg("Error while reading from websocket connection")
			wC.IsAuthenticated = false
			wC.connection.Close()
			wC.IsConnected = false
			break
		} else {
			log.Debug().Interface("authReq", authRequest).Msg("Received auth request")
			_, err = auth.VerifyToken(authRequest.Token)
			if err != nil {
				log.Debug().Err(err).Msg("Error while verifying token")
				wC.IsAuthenticated = false
				wC.connection.Close()
				wC.IsConnected = false
				break
			}
			uuid, err := auth.GetUserUuidFromToken(authRequest.Token)
			if err != nil {
				log.Debug().Err(err).Msg("Error while getting user uuid from token")
				wC.IsAuthenticated = false
				wC.connection.Close()
				wC.IsConnected = false
				break
			}
			wC.userUuid = uuid
			wC.IsAuthenticated = true
		}
	}
}
