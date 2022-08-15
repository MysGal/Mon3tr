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

	//DiscussionQuery("测试主题", 0, 1)
	//return
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
	app.Get("/discussion/:did/:start/:count", handlers.FloorQueryHandler)

	app.Listen(":2333")

	utils.GlobalLogger.Info("Running cleanup progress")
	// Clean up
	defer database.GlobalTopicIndex.Close()
	defer database.GlobalDiscussionIndex.Close()
	defer database.GlobalDatabase.Close()
}

//func Test() {
//	if err := database.GlobalDatabase.Update(
//		func(tx *nutsdb.Tx) error {
//			bucket := "bucketForList"
//			key := []byte("myList")
//			for i := 0; i < 10; i++ {
//				val := []byte("val" + strconv.Itoa(i))
//				tx.LPush(bucket, key, val)
//			}
//			return nil
//		}); err != nil {
//		log.Fatal(err)
//	}
//
//	if err := database.GlobalDatabase.View(
//		func(tx *nutsdb.Tx) error {
//			bucket := "bucketForList"
//			key := []byte("myList")
//			if items, err := tx.LRange(bucket, key, 0, -2); err != nil {
//				return err
//			} else {
//				//fmt.Println(items)
//				for _, item := range items {
//					fmt.Println(string(item))
//				}
//			}
//			return nil
//		}); err != nil {
//		log.Fatal(err)
//	}
//}

//func DiscussionQuery(topic string, start int, count int) ([]types.Discussion, error) {
//	var discussions []types.Discussion
//	err := database.GlobalDatabase.View(
//		func(tx *nutsdb.Tx) error {
//			items, err := tx.LRange("discussion", []byte(topic), start, count)
//			if err != nil {
//				return err
//			}
//			for _, item := range items {
//				var discussion types.Discussion
//				err := jsoniter.Unmarshal(item, &discussion)
//				if err != nil {
//					return err
//				}
//				discussions = append(discussions, discussion)
//			}
//			return nil
//		})
//	if err != nil {
//		return nil, err
//	}
//	return discussions, nil
//}
