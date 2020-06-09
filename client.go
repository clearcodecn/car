package car

import "github.com/gorilla/websocket"

func NewClient(wsconn *websocket.Conn) *Agent {
	agent := newAgent(wsconn)
	go agent.Run()
	return agent
}
