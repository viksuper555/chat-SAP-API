package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type MessageBody struct {
	SenderId   string   `json:"sender_id,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Message    string   `json:"message,omitempty"`
}

func ChanStream(c *gin.Context) {
	chanStream := make(chan int, 10)
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			chanStream <- i
			time.Sleep(time.Second * 1)
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
	var message MessageBody
	if err := c.BindJSON(&message); err != nil {
		// DO SOMETHING WITH THE ERROR
	}
	println("Message = " + message.Message)
	println("Sender = " + message.SenderId)
	for i := range message.Recipients {
		println("Recipient = " + message.Recipients[i])
	}
}

func Test(c *gin.Context) {
	ExampleClient()
	c.JSON(http.StatusOK, "success")
}

func main() {
	r := gin.New()
	r.POST("/message", SendMessage)
	r.GET("/test", Test)
	r.GET("/stream", ChanStream)

	err := r.Run("localhost:5000")
	if err != nil {
		return
	}
}
