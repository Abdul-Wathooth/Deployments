package main

import (
	"net/http"

	api "github.com/el10savio/TODO-Fullstack-App-Go-Gin-Postgres-React/backend/api"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

func indexView(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.JSON(http.StatusOK, gin.H{"message": "TODO APP"})
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	router.GET("/", indexView)
	router.GET("/items", api.TodoItems)
	router.POST("/items", api.CreateTodoItem)
	router.PUT("/items/:id", api.UpdateTodoItem)
	router.DELETE("/items/:id", api.DeleteTodoItem)

	return router
}

func main() {
	api.SetupPostgres()
	router := SetupRoutes()
	router.Run(":8081")
}
