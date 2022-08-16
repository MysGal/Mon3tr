package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func FloorCreateHandler(ctx *fiber.Ctx) error {
	did := ctx.Params("did")
	if did == "" {
		SendMessage(ctx, 403, "broken path")
	}
	rawBody := ctx.Body()
	var floor types.Floor
	err := jsoniter.Unmarshal(rawBody, &floor)
	if err != nil {
		SendMessage(ctx, 403, "broken body")
		return nil
	}
	// TODO: 用户身份校验

	// 校验Did是否存在
	discussion, err := database.DiscussionQueryByDid(did)
	if err != nil {
		SendMessage(ctx, 500, "discussion query error")
		return nil
	}
	if discussion.Did != did {
		SendMessage(ctx, 403, "discussion not found")
		return nil
	}

	err = database.FloorCreate(did, &floor)
	if err != nil {
		SendMessage(ctx, 500, "database error")
		return nil
	}
	SendMessage(ctx, 200, "success")
	return nil
}

func FloorQueryHandler(ctx *fiber.Ctx) error {
	did := ctx.Params("did")
	if did == "" {
		SendMessage(ctx, 403, "broken path")
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

	floors, err := database.FloorQuery(did, start, count)
	if err != nil {
		SendMessage(ctx, 500, "database error")
		return nil
	}

	type returnStruct struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    []types.Floor `json:"data"`
	}

	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data:    floors,
	}

	returnBody, err := jsoniter.Marshal(returnData)
	if err != nil {
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnBody)

	return nil
}
