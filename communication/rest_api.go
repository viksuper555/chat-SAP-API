package communication

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"messenger/db"
	"messenger/dto_model"
	"messenger/graphql"
	"messenger/services"
	"net/http"
	"time"
)

func SendMessage(c *gin.Context) {
	var msgBody dto_model.MessageBody
	if err := c.BindJSON(&msgBody); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	msg := graphql.Message{
		ID:     msgBody.Id,
		User:   msgBody.User,
		Type:   msgBody.Type,
		Sender: msgBody.Sender,
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
	var ub dto_model.UserBody
	if err := c.BindJSON(&ub); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	if _, err := db.GetUser(ub.Name); err != redis.Nil {
		c.JSON(http.StatusForbidden, "Name already exists.")
		return
	}

	u := graphql.User{
		Id:   uuid.New().String(),
		Name: ub.Name,
	}
	if err := db.SetUser(u.Name, u); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, u)
}
