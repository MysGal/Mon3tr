package main

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/handlers"
	"github.com/MysGal/Mon3tr/utils"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
)

func main() {

	database.InitIndex()
	database.InitDatabase()
	handlers.InitSessionHandler()
	utils.InitLogger()
	//database.Test()
	//return
	app := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		utils.GlobalLogger.Info("Shutdown in progress")
		_ = app.Shutdown()
	}()

	app.Post("/article/create", handlers.ArticleCreateHandler)
	app.Post("/article/:aid/update", handlers.ArticleUpdateHandler)
	app.Get("/article/:aid/detail", handlers.ArticleQueryHandler)
	app.Post("/article/:aid/comment/create", handlers.ArticleContentCreateHandler)
	app.Get("/article/:aid/comment/fetch/:startAcid", handlers.ArticleContentQueryHandler)

	app.ListenTLS(":8888", "./data/tls/cert.pem", "./data/tls/key")
	//app.Listen(":8888")

	utils.GlobalLogger.Info("Running cleanup progress")
	// Clean up
	database.GlobalIndex.Close()
	defer database.GlobalDatabase.Close()
}
