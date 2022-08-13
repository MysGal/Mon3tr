package main

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

func main() {

	database.InitIndex()
	database.InitDatabase()

	//database.Test()

	app := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Shutdown in progress")
		_ = app.Shutdown()
	}()

	//app.ListenTLS(":8888", "./data/tls/cert.pem", "./data/tls/key")
	app.Listen(":8888")

	log.Println("Running cleanup progress")
	// Clean up
	database.GlobalIndex.Close()
	defer database.GlobalDatabase.Close()
}
