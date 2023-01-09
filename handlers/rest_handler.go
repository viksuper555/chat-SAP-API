package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"messenger/dto_model"
	"messenger/internal/common"
	"messenger/model"
	"messenger/services"
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
		Type:   msgBody.Type,
		UserID: msgBody.Sender,
		Text:   msgBody.Message,
		Date:   time.Unix(msgBody.Timestamp, 0),
	}
	for _, rec := range msgBody.Recipients {
		if u, ok := services.Rm.UMap[rec]; ok {
			u.Ch <- msg
		} else {
			c.Status(http.StatusNotFound)
			return
		}
	}
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
	if err = db.Where("name = ?", ub.Username).First(&user).Error; err != nil {
		c.Status(http.StatusForbidden)
		log.Println(err)
		return
	}
	u := &model.User{
		Username: ub.Username,
		Password: ub.Password,
	}

	err = db.Create(&u).Error
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, u)
}
