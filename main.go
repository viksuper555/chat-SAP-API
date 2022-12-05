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
	"time"
)

var rm = NewRoom()

func NewRoom() *Room {
	r := new(Room)
	r.uMap = make(map[string]*User)
	return r
}

func (r *Room) LoginUser(u *User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User logged in: %s\n", u.uuid)
	r.uMap[u.uuid] = u
}

func (r *Room) LogoutUser(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User offline: %s\n", id)
	delete(r.uMap, id)
}

//func ChanStream(c *gin.Context) {
//	chanStream := make(chan int, 10)
//	go func() {
//		defer close(chanStream)
//		for i := 0; i < 10; i++ {
//			chanStream <- os.Getpid()
//			time.Sleep(time.Second * 5)
//		}
//	}()
//	c.Stream(func(w io.Writer) bool {
//		if msg, ok := <-chanStream; ok {
//			c.SSEvent("message", msg)
//			return true
//		}
//		return false
//	})
//
//}
//func ChatHandler(c *gin.Context) {
//	u := new(User)
//	u.uuid = uuid.New().String()
//	u.ch = make(chan Message)
//
//	defer rm.LogoutUser(u.uuid)
//	rm.LoginUser(u)
//
//	select {
//	case <-c.Writer.CloseNotify():
//		log.Printf("%s : disconnected\n", u.uuid)
//	case out := <-u.ch:
//		log.Printf("%s : received %+v\n", u.uuid, out)
//		c.Stream(func(w io.Writer) bool {
//			if msg, ok := <-u.ch; ok {
//				c.SSEvent("message", msg)
//				return true
//			}
//			return false
//		})
//	case <-time.After(time.Second * 60):
//		log.Println("timed out")
//	}
//}

func SendMessage(c *gin.Context) {
	var msg MessageBody
	if err := c.BindJSON(&msg); err != nil {
		// DO SOMETHING WITH THE ERROR
	}

	for _, rec := range msg.Recipients {
		if u, ok := rm.uMap[rec]; ok {
			u.ch <- msg.Message
		} else {
			c.Status(http.StatusNotFound)
			return
		}
	}
}
func Register(c *gin.Context) {
	var rb RegisterBody
	if err := c.BindJSON(&rb); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	if exists := UsernameExists(rb.Username); exists {
		c.JSON(http.StatusForbidden, "Username already exists.")
		return
	}

	id := uuid.New().String()
	if msg, ok := SetUser(id, rb.Username); !ok {
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	rb.Uuid = id
	c.JSON(http.StatusOK, rb)
}

func WebSocketHandler(ws *websocket.Conn) {
	err := ws.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		log.Printf("sad %s\n", err)
	}
	var rb RegisterBody
	if err := websocket.JSON.Receive(ws, &rb); err != nil {
		log.Printf("rb %s\n", err)
	}

	username, ok := GetUser(rb.Uuid)
	if !ok {
		updateJson, _ := json.Marshal(Message{Message: "User not found", Type: "error"})
		str := string(updateJson)
		if err := websocket.JSON.Send(ws, &str); err != nil {
			log.Println(err)
			return
		}
		ws.Close()
		return
	}

	u := User{
		uuid: rb.Uuid,
		ws:   ws,
		name: username,
	}

	rm.LoginUser(&u)
	SendLoginInfo(ws, u.name)
	BroadcastOnlineUsers()
	defer CleanUp(ws, u.uuid)
	for {
		// read in a message
		var msg MessageBody
		err := ws.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Printf("sad %s\n", err)
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("sad2 %s\n", err)
			continue
		}
		// print out that message for clarity
		fmt.Printf("%s, %s\n", msg.Message, u.uuid)
		msg.Sender = u.name
		if msg.Message.Message != "" {
			msg.Type = "message"
			msg.Timestamp = time.Now().Unix()
			Broadcast(msg.Message)
		}
	}
}

func CleanUp(ws *websocket.Conn, uuid string) {
	ws.Close()
	rm.LogoutUser(uuid)
	BroadcastOnlineUsers()
}

func Broadcast(msg Message) {
	updateJson, _ := json.Marshal(msg)

	sendMsg := string(updateJson)
	for _, u := range rm.uMap {
		if err := websocket.JSON.Send(u.ws, &sendMsg); err != nil {
			log.Println(err)
			return
		}
	}
}

func BroadcastOnlineUsers() {
	names := make([]string, len(rm.uMap))
	i := 0
	for k := range rm.uMap {
		names[i] = rm.uMap[k].name
		i++
	}

	body := BroadcastOnline{
		Users: names,
		Type:  "online",
	}
	updateJson, _ := json.Marshal(body)

	online := string(updateJson)
	for _, u := range rm.uMap {
		if err := websocket.JSON.Send(u.ws, &online); err != nil {
			log.Println(err)
			return
		}
	}
}

func SendLoginInfo(ws *websocket.Conn, name string) {
	updateJson, _ := json.Marshal(RegisterBody{Username: name, Type: "login"})
	str := string(updateJson)
	if err := websocket.JSON.Send(ws, &str); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	go func() {
		r := gin.Default()

		api := r.Group("/api")
		{
			api.POST("/message", SendMessage)
			api.POST("/register", Register)
		}

		if err := r.Run("0.0.0.0:5000"); err != nil {
			return
		}
	}()
	http.Handle("/ws", websocket.Handler(WebSocketHandler))

	if err := http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
