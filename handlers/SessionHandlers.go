package handlers

import "github.com/gofiber/fiber/v2/middleware/session"

var GlobalSessionHandler *session.Store

func InitSessionHandler() {
	store := session.New()
	GlobalSessionHandler = store
}
