package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(c *gin.Context) {
	ExampleClient()
	c.JSON(http.StatusOK, "success")
}
func main() {
	r := gin.New()

	r.GET("/test", Test)

	err := r.Run("localhost:5000")
	if err != nil {
		return
	}
}
