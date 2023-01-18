package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"messenger/config"
	"messenger/handlers"
	"messenger/hub"
	"messenger/internal/common"
	"os"
)

const defaultPort = "8080"

func main() {
	appConfig := config.NewFromEnv()

	db, err := common.InitDb(appConfig)
	if err != nil {
		log.Fatalf("Failed to init database: %s", err)
	}
	defer dbDisconnect(db)

	r := gin.Default()

	ctx := &common.Context{
		Database: db,
	}
	r.Use(func(c *gin.Context) {
		con := context.WithValue(c.Request.Context(), "ctx", ctx)
		c.Request = c.Request.WithContext(con)
		c.Next()
	})
	go hub.MainHub.Run()

	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWs(c, hub.MainHub)
	})
	// REST
	r.POST("/api/message", handlers.SendMessage)
	r.POST("/api/register", handlers.Register)
	r.POST("/api/join", handlers.Join)
	r.POST("/api/leave", handlers.Leave)

	// GQL
	r.GET("", handlers.PlaygroundHandler())
	r.POST("/graph/api", handlers.GraphqlHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(r.Run(":" + port))
}

func dbDisconnect(dbConn *gorm.DB) {
	database, err := dbConn.DB()
	if err != nil {
		log.Fatalf("Failed to close database connection: %s", err)
	}
	err = database.Close()
	if err != nil {
		log.Fatalf("Failed to close database connection: %s", err)
	}
}
