package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/MysGal/Mon3tr/utils"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"regexp"
)

func TopicCreateHandler(ctx *fiber.Ctx) error {
	rawBody := ctx.Body()
	var topic types.Topic
	err := jsoniter.Unmarshal(rawBody, &topic)
	if err != nil {
		SendMessage(ctx, 403, "broken body")
		return nil
	}
	// 数据校验
	// TODO: 校验topic格式，不含有空格，全小写，空格用-替代
	legal, _ := regexp.MatchString("^[-a-z0-9]+$", topic.Topic)
	if topic.Name == "" || topic.Tags == nil || !legal {
		SendMessage(ctx, 403, "missing topic field")
		return nil
	}
	// 数据写入
	err = database.TopicCreate(topic)
	if err != nil {
		utils.GlobalLogger.Info(err)
		SendMessage(ctx, 500, "topic creation failed, maybe a server error")
		return nil
	}

	SendMessage(ctx, 200, "success")
	return nil
}

func TopicUpdateHandler(ctx *fiber.Ctx) error {
	topicParam := ctx.Params("topic")
	rawBody := ctx.Body()
	var topic types.Topic
	jsoniter.Unmarshal(rawBody, &topic)
	if topicParam != topic.Topic {
		SendMessage(ctx, 403, "missing topic or broken body")
		return nil
	}

	err := database.TopicUpdate(topic)
	if err != nil {
		utils.GlobalLogger.Error(err)
		SendMessage(ctx, 500, "server error")
		return nil
	}
	SendMessage(ctx, 200, "success")
	return nil
}

func TopicQueryAllHandler(ctx *fiber.Ctx) error {

	allTopic, err := database.TopicQueryAll()
	if err != nil {
		utils.GlobalLogger.Error(err)
		SendMessage(ctx, 500, "server error")
		return nil
	}

	type returnStruct struct {
		Code    int64         `json:"code"`
		Message string        `json:"message"`
		Data    []types.Topic `json:"data"`
	}

	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data:    allTopic,
	}
	returnJson, err := jsoniter.Marshal(returnData)
	if err != nil {
		utils.GlobalLogger.Error(err)
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnJson)
	return nil
}

func TopicDiscussionQueryHandler(ctx *fiber.Ctx) error {
	topic := ctx.Params("topic")

	discussionType := ctx.Params("type")
	if discussionType == ":type" || topic == ":start" {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	// 与DatabaseDiscussionActions中相同
	var prefix string
	switch discussionType {
	case "release":
		prefix = "release"
	case "discussion":
		prefix = "discussion"
	default:
		SendMessage(ctx, 403, "unknown discussion type")
		return nil
	}

	start, err := ctx.ParamsInt("start")
	if err != nil {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	count, err := ctx.ParamsInt("count")
	if err != nil {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	if start < 0 {
		start = 0
	}

	if count > 10 {
		count = 10
	}

	// 查询数据库
	discussions, err := database.DiscussionQueryByTopic(topic, prefix, start, count)
	if err != nil {
		SendMessage(ctx, 500, "database error")
		return nil
	}

	type returnStruct struct {
		Code    int64              `json:"code"`
		Message string             `json:"message"`
		Data    []types.Discussion `json:"data"`
	}

	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data:    discussions,
	}
	returnJson, err := jsoniter.Marshal(returnData)
	if err != nil {
		utils.GlobalLogger.Error(err)
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnJson)
	return nil
}
