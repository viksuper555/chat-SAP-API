package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var rm = NewRoom()

type Message struct {
	SenderId string `json:"sender_id,omitempty"`
	Message  string `json:"message,omitempty"`
}
type MessageBody struct {
	Message
	Recipients []string `json:"recipients,omitempty"`
}

type User struct {
	uuid string
	ch   chan Message
}

type Room struct {
	users map[string]*User
	mutex sync.Mutex
}

func NewRoom() *Room {
	r := new(Room)
	r.users = make(map[string]*User)
	return r
}

func (r *Room) AddUser(u *User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("Adding user: %s\n", u.uuid)
	r.users[u.uuid] = u
}

func (r *Room) DeleteUser(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("Deleting user: %s\n", id)
	delete(r.users, id)
}

func ChanStream(c *gin.Context) {
	chanStream := make(chan int, 10)
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			chanStream <- os.Getpid()
			time.Sleep(time.Second * 5)
		}
	}()
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

}

func SendMessage(c *gin.Context) {
	var msg MessageBody
	if err := c.BindJSON(&msg); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	for _, rec := range msg.Recipients {
		if u, ok := rm.users[rec]; ok {
			u.ch <- msg.Message
		} else {
			c.Status(http.StatusNotFound)
			return
		}
	}
	//println("Message = " + message.Message)
	//println("Sender = " + message.SenderId)
	//for i := range message.Recipients {
	//	println("Recipient = " + message.Recipients[i])
	//}
}

func main() {
	r := gin.Default()
	//go rm.InitBackgroundTask("Hey")

	r.POST("/message", SendMessage)
	r.GET("/stream", ChanStream)

	r.GET("/chat", func(c *gin.Context) {
		u := new(User)
		u.uuid = uuid.New().String()
		u.ch = make(chan Message)

		defer rm.DeleteUser(u.uuid)
		rm.AddUser(u)

		select {
		case <-c.Writer.CloseNotify():
			log.Printf("%s : disconnected\n", u.uuid)
		case out := <-u.ch:
			log.Printf("%s : received %+v\n", u.uuid, out)
			c.Stream(func(w io.Writer) bool {
				if msg, ok := <-u.ch; ok {
					c.SSEvent("message", msg)
					return true
				}
				return false
			})
		case <-time.After(time.Second * 60):
			log.Println("timed out")
		}
	})

	err := r.Run("localhost:5000")
	if err != nil {
		return
	}
}
