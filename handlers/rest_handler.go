package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"messenger/dto_model"
	"messenger/hub"
	"messenger/internal/common"
	"messenger/model"
	"net/http"
	"time"
)

func SendMessage(c *gin.Context) {
	var msgBody dto_model.MessageBody
	err := c.BindJSON(&msgBody)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	msg := model.Message{
		ID:     msgBody.Id,
		UserID: msgBody.SenderId,
		Text:   msgBody.Message,
		RoomID: msgBody.RoomId,
		Date:   time.Unix(msgBody.Timestamp, 0),
	}
	common.Db.Create(&msg)
	bytes, err := json.Marshal(msg)
	hub.MainHub.Rooms[msgBody.RoomId].Broadcast <- bytes
	//hub.MainHub.Broadcast <- bytes
}

func Register(c *gin.Context) {
	ctx := c.Request.Context().Value("ctx").(*common.Context)
	db := ctx.Database

	var ub dto_model.UserBody
	err := c.BindJSON(&ub)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
	var user model.User
	if err = db.Where("username = ?", ub.Username).First(&user).Error; err == nil {
		c.Status(http.StatusForbidden)
		log.Println(err)
		return
	}
	var gr model.Room
	if err = db.Where("id = ?", "global").First(&gr).Error; err != nil {
		c.Status(http.StatusForbidden)
		log.Println(err)
		return
	}

	u := &model.User{
		Username: ub.Username,
		Password: ub.Password,
		Rooms:    []model.Room{gr},
	}

	err = db.Create(&u).Error
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func Join(c *gin.Context) {
	var jb dto_model.JoinRoomBody
	err := c.BindJSON(&jb)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
	client, ok := hub.MainHub.Clients[jb.UserId]
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}
	room, ok := hub.MainHub.Rooms[jb.RoomId]
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}

	room.Clients[client] = true
	err = common.AddUserToRoom(common.Db, jb.UserId, jb.RoomId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
}

func Leave(c *gin.Context) {
	var jb dto_model.JoinRoomBody
	err := c.BindJSON(&jb)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
	client, ok := hub.MainHub.Clients[jb.UserId]
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}
	room, ok := hub.MainHub.Rooms[jb.RoomId]
	if !ok {
		c.Status(http.StatusForbidden)
		return
	}

	delete(room.Clients, client)
	err = common.RemoveUserFromRoom(common.Db, jb.UserId, jb.RoomId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Println(err)
		return
	}
}
