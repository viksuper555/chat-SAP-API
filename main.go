package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
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
	ws   *websocket.Conn
}

type RegisterBody struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
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
func ChatHandler(c *gin.Context) {
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
}
func Register(c *gin.Context) {
	var body RegisterBody
	if err := c.BindJSON(&body); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	id := uuid.New().String()
	if msg, ok := SetUser(body.Username, id); !ok {
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	c.JSON(http.StatusOK, id)
}

func Echo(ws *websocket.Conn) {
	var err error
	updateJson, err := json.Marshal(`Hi, the handshake it complete!`)

	msg := string(updateJson)
	for {
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
		} else {
			fmt.Println("Sending")
		}
		time.Sleep(time.Second * 1)
	}

}

func WebSocketHandler(ws *websocket.Conn) {
	u := new(User)
	u.uuid = uuid.New().String()
	u.ws = ws

	defer rm.DeleteUser(u.uuid)
	rm.AddUser(u)
	defer ws.Close()
	for {
		// read in a message
		var msg MessageBody
		err := ws.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Printf("sad %s\n", err)
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			log.Printf("sad2 %s\n", err)
		}
		// print out that message for clarity
		fmt.Printf("%s, %s\n", msg.Message, msg.SenderId)
		if msg.Message.Message != "" {
			Broadcast(msg.Message.Message)
		}
	}
}
func Broadcast(msg string) {
	updateJson, _ := json.Marshal(msg)

	sendMsg := string(updateJson)
	for _, u := range rm.users {
		if err := websocket.JSON.Send(u.ws, &sendMsg); err != nil {
			log.Println(err)
			return
		}
	}

}

func main() {
	//http.Handle("/ws", websocket.Handler(Echo))
	http.Handle("/ws", websocket.Handler(WebSocketHandler))

	if err := http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/message", SendMessage)
		api.GET("/stream", ChanStream)
		api.POST("/register", Register)
		api.GET("/chat", ChatHandler)
		api.GET("/test", func(c *gin.Context) {
			message := Message{SenderId: "Ivan", Message: "Hello"}
			c.JSON(http.StatusOK, message)
		})
	}

	if err := r.Run("0.0.0.0:5000"); err != nil {
		return
	}
}
