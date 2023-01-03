package main

import (
	"encoding/json"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
	"log"
	"messenger/communication"
	"messenger/config"
	"messenger/graph/generated"
	graph "messenger/graph/resolvers"
	"messenger/internal/common"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	appConfig := config.NewFromEnv()

	log.Print("Connecting to DB...")
	db, err := common.InitDb(appConfig)
	if err != nil {
		log.Fatalf("Failed to init database: %s", err)
	}
	customCtx := &common.CustomContext{
		Database: db,
	}

	defer func(dbConn *gorm.DB) {
		err := dbDisconnect(dbConn)
		if err != nil {
			log.Fatal("Failed to close database connection: %s", err)
		}
	}(db)
	log.Print("Connected to DB!")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	// REST
	mux.HandleFunc("/api/message", communication.SendMessage)
	mux.HandleFunc("/api/register", communication.Register)

	// WS
	mux.Handle("/ws", websocket.Handler(communication.WebSocketHandler))

	// GQL
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	mux.Handle("/graph/api", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", common.CreateContext(customCtx, srv))

	//go func() {
	//	r := gin.Default()
	//
	//	api := r.Group("/api")
	//	{
	//		api.POST("/message", communication.SendMessage)
	//		api.POST("/register", communication.Register)
	//	}
	//
	//	if err = r.Run("0.0.0.0:5000"); err != nil {
	//		return
	//	}
	//}()
	//
	//go func() {
	//	mux.Handle("/ws", websocket.Handler(communication.WebSocketHandler))
	//
	//	if err := http.ListenAndServe("0.0.0.0:9000", nil); err != nil {
	//		log.Fatal("ListenAndServe:", err)
	//	}
	//}()
	//
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	//
	//mux.Handle("/api", playground.Handler("GraphQL playground", "/query"))
	//mux.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

//func handleREST() {
//	myRouter := mux.NewRouter().StrictSlash(true)
//	myRouter.HandleFunc("/api/message", communication.SendMessage)
//}

func dbConnect(appConfig config.Config) (*gorm.DB, error) {
	appConfigJSON, err := json.MarshalIndent(appConfig, "", "  ")
	if err != nil {
		log.Printf("Error on writing config as json: %s", err.Error())
	}

	log.Printf("appconfig: %s\n", string(appConfigJSON))
	log.Printf("os env: %s", os.Environ())

	return common.InitDb(appConfig)
}

func dbDisconnect(dbConn *gorm.DB) error {
	database, err := dbConn.DB()
	if err != nil {
		return err
	}

	return database.Close()
}
