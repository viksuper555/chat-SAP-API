package handlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"messenger/dto_model"
	"messenger/model"
	"messenger/services"
)

func Broadcast(msg interface{}) {
	updateJson, _ := json.Marshal(msg)
	for _, u := range services.Rm.UMap {
		if err := u.Ws.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
			return
		}
	}
}

func BroadcastOnlineUsers() {
	names := make([]string, len(services.Rm.UMap))
	i := 0
	for k := range services.Rm.UMap {
		names[i] = services.Rm.UMap[k].Name
		i++
	}

	body := dto_model.BroadcastOnline{
		Users: names,
		Type:  "online",
	}
	updateJson, _ := json.Marshal(body)
	for _, u := range services.Rm.UMap {
		if err := u.Ws.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
			return
		}
	}
}

func SendLoginInfo(u *model.User) {
	updateJson, err := json.Marshal(dto_model.UserBody{User: *u, Type: "login"})
	if err != nil {
		log.Println(err)
		return
	}
	if err := u.Ws.WriteMessage(websocket.TextMessage, updateJson); err != nil {
		log.Println(err)
		return
	}
}

func CleanUp(ws *websocket.Conn, id int) {
	ws.Close()
	services.Rm.LogoutUser(id)
	BroadcastOnlineUsers()
}
