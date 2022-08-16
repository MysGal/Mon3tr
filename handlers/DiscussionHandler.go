package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func DiscussionCreateHandler(ctx *fiber.Ctx) error {
	rawBody := ctx.Body()
	var discussion types.Discussion
	err := jsoniter.Unmarshal(rawBody, &discussion)
	if err != nil {
		SendMessage(ctx, 403, "broken body")
		return nil
	}

	// 数据校验
	// TODO: 根据Session校验用户
	if discussion.Creator == 0 || discussion.Topic == "" || discussion.Title == "" || discussion.Floor == nil {
		SendMessage(ctx, 403, "broken body")
		return nil
	}

	// 检查是否存在该主题
	_, err = database.TopicQueryDetail(discussion.Topic)
	if err != nil {
		SendMessage(ctx, 403, "wrong topic")
		return nil
	}

	firstFloor := discussion.Floor
	discussion.Floor = nil

	did, err := database.DiscussionCreate(discussion)
	if err != nil {
		SendMessage(ctx, 500, "discussion database error")
		return nil
	}

	// 新建第一层楼
	err = database.FloorCreate(did, firstFloor)
	if err != nil {
		// TODO: 删除帖子
		SendMessage(ctx, 500, "floor database error")
		return nil
	}

	type returnStruct struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Did string `json:"did"`
		} `json:"data"`
	}

	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data: struct {
			Did string `json:"did"`
		}(struct{ Did string }{Did: did}),
	}

	returnBody, err := jsoniter.Marshal(returnData)
	if err != nil {
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnBody)
	return nil
}

func DiscussionDetailQueryHandler(ctx *fiber.Ctx) error {
	did := ctx.Params("did")
	if did == "" {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	discussion, err := database.DiscussionQueryByDid(did)
	if err != nil {
		SendMessage(ctx, 500, "database error")
		return nil
	}

	if discussion.Did == "" {
		SendMessage(ctx, 403, "did not found")
		return nil
	}

	type returnStruct struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    types.Discussion `json:"data"`
	}

	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data:    discussion,
	}

	returnBody, err := jsoniter.Marshal(returnData)

	if err != nil {
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnBody)

	return nil
}
