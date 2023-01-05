package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"messenger/cache"
	"messenger/dto_model"
	"messenger/services"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		log.Printf("sad %s\n", err)
	}
	_, bytes, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	var ub dto_model.UserBody
	if err = json.Unmarshal(bytes, &ub); err != nil {
		log.Println(err)
		return
	}

	u, err := cache.GetUser(ub.Name)
	if err != nil {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "User not found", Type: "error"})
		if err := conn.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}
	if u.ID != ub.ID {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "Wrong token.", Type: "error"})
		if err = conn.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
			return
		}
	}
	u.Ws = conn

	services.Rm.LoginUser(u)
	SendLoginInfo(u)
	BroadcastOnlineUsers()
	defer CleanUp(conn, u.ID)

	// read in a message
	var msg dto_model.MessageBody
	for {
		err = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			log.Printf("sad %s\n", err)
		}
		_, bytes, err = conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err = json.Unmarshal(bytes, &msg); err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Printf("%s, %s\n", msg.Message, u.ID)
		msg.Sender = u.ID
		if msg.Message != "" {
			msg.Type = "message"
			msg.Timestamp = time.Now().Unix()
			Broadcast(msg)
		}
	}
}
