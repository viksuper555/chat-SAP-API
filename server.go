package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"log"
	"messenger/communication"
	"messenger/config"
	"messenger/db"
	"messenger/graphql"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	appConfig := config.NewFromEnv()

	log.Print("Connecting to DB...")
	database, err := dbConnect(appConfig)
	if err != nil {
		log.Fatalf("Failed to init database: %s", err)
	}
	defer dbDisconnect(database)
	log.Print("Connected to DB!")

	go func() {
		r := gin.Default()

		api := r.Group("/api")
		{
			api.POST("/message", communication.SendMessage)
			api.POST("/register", communication.Register)
		}

		if err = r.Run("0.0.0.0:5000"); err != nil {
			return
		}
	}()
	go func() {
		http.Handle("/ws", websocket.Handler(communication.WebSocketHandler))

		if err := http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func dbConnect(appConfig config.Config) (db.Database, error) {
	appConfigJSON, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		log.Printf("Error on writing config as json: %s", err.Error())
	}

	log.Printf("appconfig: %s\n", string(appConfigJSON))
	log.Printf("os env: %s", os.Environ())

	return db.Init(appConfig)
}

func dbDisconnect(database db.Database) {
	if database == nil {
		return
	}
	err := database.Close()
	if err != nil {
		log.Printf("Error closing database connection: %s", err.Error())
	}
	log.Print("Successfully closed connection to database")
}
