package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"messenger/config"
	"messenger/handlers"
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
		ctx := context.WithValue(c.Request.Context(), "ctx", &ctx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	r.GET("/ws", func(c *gin.Context) {
		handlers.HandleWebsocket(c.Writer, c.Request)
	})
	// REST
	r.POST("/api/message", handlers.SendMessage)
	r.POST("/api/register", handlers.Register)

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

//func dbConnect(appConfig config.Config) (*gorm.DB, error) {
//	appConfigJSON, err := json.MarshalIndent(appConfig, "", "  ")
//	if err != nil {
//		log.Printf("Error on writing config as json: %s", err.Error())
//	}
//
//	log.Printf("appconfig: %s\n", string(appConfigJSON))
//	log.Printf("os env: %s", os.Environ())
//
//	return common.InitDb(appConfig)
//}

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
