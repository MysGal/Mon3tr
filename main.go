package main

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/handlers"
	"github.com/MysGal/Mon3tr/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

func main() {

	utils.InitLogger()
	database.InitIndex()
	database.InitDatabase()
	handlers.InitSessionHandler()

	app := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		utils.GlobalLogger.Info("Shutdown in progress")
		_ = app.Shutdown()
	}()

	// 主题处理器
	app.Post("/topic/create", handlers.TopicCreateHandler)
	app.Get("/topic/all", handlers.TopicQueryAllHandler)
	app.Post("/topic/:topic/update", handlers.TopicUpdateHandler)
	// 帖子处理器
	app.Get("/topic/:topic/discussion/:type/:start/:count", handlers.TopicDiscussionQueryHandler)
	app.Post("/discussion/create", handlers.DiscussionCreateHandler)
	// 楼处理器
	app.Post("/discussion/:did/comment", handlers.FloorCreateHandler)
	app.Get("/discussion/:did/detail", handlers.DiscussionDetailQueryHandler)
	app.Get("/discussion/:did/:start/:count", handlers.FloorQueryHandler)
	// 搜索处理器
	app.Get("/search/:query/:from", handlers.SearchHandler)

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	err := app.Listen(":2333")
	if err != nil {
		utils.GlobalLogger.Fatal(err)
	}

	utils.GlobalLogger.Info("Running cleanup progress")
	// Clean up
	err = database.GlobalIndex.Close()
	if err != nil {
		utils.GlobalLogger.Fatal(err)
	}
	err = database.GlobalDatabase.Close()
	if err != nil {
		utils.GlobalLogger.Fatal(err)
	}
	utils.GlobalLogger.Info("Cleanup progress finished")
}
