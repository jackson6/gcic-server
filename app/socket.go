package app

import (
	"github.com/googollee/go-socket.io"
	"log"
)

func(a App)InitializeSocket(){
	a.Socket.On("connection", func(so socketio.Socket) {
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println(msg)
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	a.Socket.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
}