package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"log"
	"messenger/graph/generated"
	graph "messenger/graph/resolvers"
	"messenger/internal/common"
)

func GraphqlHandler(c *gin.Context) {
	//TODO: Find out why the context becomes a pointer to pointer
	ctx, ok := c.Request.Context().Value("ctx").(**common.Context)
	if !ok {
		log.Fatal()
	}
	db := (*ctx).Database

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))
	h.ServeHTTP(c.Writer, c.Request)
}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graph/api")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
