package communication

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"io"
	"log"
	"messenger/db"
	"messenger/dto_model"
	"messenger/internal/common"
	"messenger/model"
	"messenger/services"
	"net/http"
	"time"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var msgBody dto_model.MessageBody
	err := json.Unmarshal(reqBody, &msgBody)
	if err != nil {
		return
	}
	msg := model.Message{
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
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
func Register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var ub dto_model.UserBody
	err := json.Unmarshal(reqBody, &ub)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if _, err := db.GetUser(ub.Name); err != redis.Nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
		return
	}
	u := &model.User{
		Name:     ub.Name,
		Password: ub.Password,
	}

	err = common.Db.Create(&u).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
