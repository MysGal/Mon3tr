package handlers

import (
	"github.com/MysGal/Mon3tr/utils"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendMessage(ctx *fiber.Ctx, code int, message string) {
	returnMessageRaw := Message{Code: code, Message: message}
	returnMessage, err := jsoniter.Marshal(returnMessageRaw)
	if err != nil {
		utils.GlobalLogger.Error(string(returnMessage))
		return
	}
	ctx.Send(returnMessage)
}
