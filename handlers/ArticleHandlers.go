package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func ArticleCreateHandler(ctx *fiber.Ctx) error {
	rawBody := ctx.Body()
	rawarticle := types.Article{}
	err := jsoniter.Unmarshal(rawBody, &rawarticle)
	if err != nil {
		SendMessage(ctx, 403, "Unable to unmarshal, maybe broken body")
		return nil
	}

	// 校验合法性并清理无效数据，确保数据库干净
	if rawarticle.ArticleType == "" {
		SendMessage(ctx, 403, "Missing article type")
		return nil
	}

	if len(rawarticle.ArticleContent) != 1 {
		SendMessage(ctx, 403, "Wrong formate, too many article content")
		return nil
	}

	article := types.Article{
		ArticleType: rawarticle.ArticleType,
		ArticleContent: []types.ArticleContent{{
			ContentFloor:   1,
			ArticleContent: rawarticle.ArticleContent[0].ArticleContent,
			Author: types.User{
				Uid: rawarticle.ArticleContent[0].Author.Uid,
			},
			Likes: 0},
		},
		ArticleTags: rawarticle.ArticleTags,
	}

	if rawarticle.RelatedGalGame != nil {
		article.RelatedGalGame = rawarticle.RelatedGalGame
	} else if rawarticle.ArticleType == "galgame" {
		SendMessage(ctx, 403, "Broken Body, is galgame but has no content")
		return nil
	}

	// 数据库写入部分
	err = database.ArticleCreate(article)
	if err != nil {
		SendMessage(ctx, 500, "Database input error, this may be a server error")
		return nil
	}

	SendMessage(ctx, 200, "success")

	return nil
}
