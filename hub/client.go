package hub

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"messenger/dto_model"
	"messenger/internal/common"
	"messenger/model"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var DB *gorm.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn
	user *model.User
	// Buffered channel of outbound messages.
	send chan interface{}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		c.hub.Broadcast <- dto_model.ActiveUsersUpdate{Type: "online", UserId: c.user.ID, Connected: false}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	var msg dto_model.MessageBody
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Println(err)
			return
		}
		msg.SenderId = c.user.ID
		msg.SenderName = c.user.Username
		msg.Timestamp = time.Now().Unix()
		msg.Type = "message"
		c.hub.Broadcast <- msg
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case obj, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			bytes, _ := json.Marshal(obj)
			w.Write(bytes)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				bytes, err = json.Marshal(<-c.send)
				w.Write(bytes)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(c *gin.Context, hub *Hub) {
	w := c.Writer
	r := c.Request
	ctx := c.Request.Context().Value("ctx").(*common.Context)
	DB = ctx.Database
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan interface{}, 256)}
	hub.register <- client

	//region Initial Login
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
	err = DB.Where("username = ? AND password = ?", ub.Username, ub.Password).First(&user).Error
	//u, err := cache.GetUser(ub.Username)
	if err != nil {
		updateJson, _ := json.Marshal(dto_model.MessageBody{Message: "User not found", Type: "error"})
		if err := conn.WriteMessage(websocket.TextMessage, updateJson); err != nil {
			log.Println(err)
		}
		conn.Close()
		return
	}
	client.user = &user

	var onlineUsers = getOnlineUsers(Hub1)
	updateJson, err := json.Marshal(dto_model.DataOnLogin{User: user, Type: "login", OnlineUsers: onlineUsers})
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, updateJson); err != nil {
		log.Println(err)
	}
	//client.send <- updateJson
	// endregion

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	hub.Broadcast <- dto_model.ActiveUsersUpdate{Type: "online", UserId: client.user.ID, Connected: true}
}
