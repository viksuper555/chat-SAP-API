package communication

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"messenger/db"
	"messenger/dto_model"
	"messenger/services"
	"time"
)

func WebSocketHandler(ws *websocket.Conn) {
	err := ws.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		log.Printf("sad %s\n", err)
	}
	var ub dto_model.UserBody
	if err := websocket.JSON.Receive(ws, &ub); err != nil {
		log.Printf("ub %s\n", err)
	}

	u, err := db.GetUser(ub.Name)
	if err != nil {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "User not found", Type: "error"})
		str := string(updateJson)
		if err = websocket.JSON.Send(ws, &str); err != nil {
			log.Println(err)
			return
		}
		ws.Close()
		return
	}
	if u.Id != ub.Id {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "Wrong token.", Type: "error"})
		str := string(updateJson)
		if err = websocket.JSON.Send(ws, &str); err != nil {
			log.Println(err)
			return
		}
	}
	u.Ws = ws

	services.Rm.LoginUser(u)
	SendLoginInfo(u)
	BroadcastOnlineUsers()
	defer CleanUp(ws, u.Id)

	// read in a message
	var msg dto_model.MessageBody
	for {
		err = ws.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			log.Printf("sad %s\n", err)
		}

		if err = websocket.JSON.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("timeout %s\n", err)
			continue
		}
		// print out that message for clarity
		fmt.Printf("%s, %s\n", msg.Message, u.Id)
		msg.Sender = u.Id
		if msg.Message != "" {
			msg.Type = "message"
			msg.Timestamp = time.Now().Unix()
			Broadcast(msg)
		}
	}
}
