package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"messenger/dto_model"
	"messenger/internal/common"
	"messenger/model"
	"messenger/services"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWebsocket(c *gin.Context) {
	w := c.Writer
	r := c.Request
	ctx := c.Request.Context().Value("ctx").(*common.Context)
	db := ctx.Database
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
	var user model.User
	err = db.Where("name = ? AND password = ?", ub.Username, ub.Password).First(&user).Error
	//u, err := cache.GetUser(ub.Username)
	if err != nil {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "User not found", Type: "error"})
		if err := conn.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}

	user.Ws = conn

	services.Rm.LoginUser(&user)
	SendLoginInfo(&user)
	BroadcastOnlineUsers()
	defer CleanUp(conn, user.ID)

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
		fmt.Printf("%s, %s\n", msg.Message, user.ID)
		msg.Sender = user.ID
		err = db.Create(&model.Message{
			Text: msg.Message, UserID: msg.Sender,
			Date: time.Unix(msg.Timestamp, 0),
		}).Error
		if msg.Message != "" {
			msg.Type = "message"
			msg.Timestamp = time.Now().Unix()
			Broadcast(msg)
		}
	}
}
