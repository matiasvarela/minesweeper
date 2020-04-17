package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matiasvarela/minesweeper/internal/game"
	"github.com/matiasvarela/minesweeper/internal/storage/localsto"
)

func main() {
	router := gin.New()

	conf := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Auth-Token"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(conf))

	routes(router)

	err := router.Run(":8080")
	if err != nil {
		panic("run server has fail")
	}
}

func routes(router *gin.Engine) {
	gameHttpHandler := game.NewHttpHandler(
		game.NewService(localsto.NewGameStorage()),
	)

	router.POST("/games", gameHttpHandler.Create)
	router.GET("/games/:id", gameHttpHandler.Get)
	router.PUT("/games/:id/play-square", gameHttpHandler.PlaySquare)
	router.PUT("/games/:id/mark-square", gameHttpHandler.MarkSquare)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}
