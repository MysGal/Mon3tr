package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func ArticleCreateHandler(ctx *fiber.Ctx) error {
	rawBody := ctx.Body()
	rawArticle := types.Article{}
	err := jsoniter.Unmarshal(rawBody, &rawArticle)
	if err != nil {
		SendMessage(ctx, 403, "Unable to unmarshal, maybe broken body")
		return nil
	}

	// 校验合法性并清理无效数据，确保数据库干净
	if rawArticle.ArticleType == "" {
		SendMessage(ctx, 403, "Missing article type")
		return nil
	}

	if rawArticle.ArticleTitle == "" {
		SendMessage(ctx, 403, "Missing article title")
		return nil
	}

	if len(rawArticle.ArticleContent) != 1 {
		SendMessage(ctx, 403, "Wrong format, too many article content")
		return nil
	}

	article := types.Article{
		ArticleType:  rawArticle.ArticleType,
		ArticleTitle: rawArticle.ArticleTitle,
		ArticleContent: []types.ArticleContent{{
			ArticleContent: rawArticle.ArticleContent[0].ArticleContent,
			Author: types.User{
				Uid: rawArticle.ArticleContent[0].Author.Uid,
			},
			Likes: 0},
		},
		ArticleTags: rawArticle.ArticleTags,
	}

	if rawArticle.RelatedGalGame != nil {
		article.RelatedGalGame = rawArticle.RelatedGalGame
	} else if rawArticle.ArticleType == "galgame" {
		SendMessage(ctx, 403, "Broken Body, is galgame but has no content")
		return nil
	}

	// 数据库写入部分
	aid, err := database.ArticleCreate(article)
	if err != nil {
		SendMessage(ctx, 500, "Database input error, this may be a server error")
		return nil
	}

	successResponse, err := jsoniter.Marshal(struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Aid int64 `json:"aid"`
		} `json:"data"`
	}{Code: 200, Message: "success", Data: struct {
		Aid int64 `json:"aid"`
	}(struct{ Aid int64 }{Aid: aid})})

	if err != nil {
		SendMessage(ctx, 500, "Marshal response error, this is a server error")
		return nil
	}

	ctx.Send(successResponse)

	return nil
}

func ArticleUpdateHandler(ctx *fiber.Ctx) error {
	aid, err := ctx.ParamsInt("aid")
	if err != nil {
		SendMessage(ctx, 403, "Unable to get aid, maybe broken path")
		return nil
	}
	rawBody := ctx.Body()
	rawArticle := types.Article{}
	err = jsoniter.Unmarshal(rawBody, &rawArticle)
	if err != nil {
		SendMessage(ctx, 403, "Unable to unmarshal, maybe broken body")
		return nil
	}
	if rawArticle.ArticleContent != nil {
		SendMessage(ctx, 403, "You should not update article content with this API")
		return nil
	}

	// 数据校验和清理
	if rawArticle.Aid == 0 || int64(aid) != rawArticle.Aid {
		SendMessage(ctx, 403, "Broken body, missing Aid")
		return nil
	}

	err = database.ArticleUpdate(rawArticle)
	if err != nil {
		SendMessage(ctx, 403, "Database input error, this may be a server error")
		return nil
	}

	SendMessage(ctx, 200, "success")
	return nil
}

func ArticleQueryHandler(ctx *fiber.Ctx) error {
	aid, err := ctx.ParamsInt("aid")
	if err != nil {
		SendMessage(ctx, 403, "Unable to get aid, maybe broken path")
		return nil
	}
	article, err := database.ArticleQuery(int64(aid))
	if err != nil {
		SendMessage(ctx, 500, "Database input error, this may be a server error")
		return nil
	}

	successResponse, err := jsoniter.Marshal(struct {
		Code    int64         `json:"code"`
		Message string        `json:"message"`
		Data    types.Article `json:"data"`
	}{Code: 200, Message: "success", Data: article})

	if err != nil {
		SendMessage(ctx, 500, "Marshal response error, this is a server error")
		return nil
	}

	ctx.Send(successResponse)

	return nil
}

func ArticleContentCreateHandler(ctx *fiber.Ctx) error {
	aid, err := ctx.ParamsInt("aid")
	if err != nil {
		SendMessage(ctx, 403, "Unable to get aid, maybe broken path")
		return nil
	}
	rawBody := ctx.Body()
	articleContent := types.ArticleContent{}
	err = jsoniter.Unmarshal(rawBody, &articleContent)
	if err != nil {
		SendMessage(ctx, 403, "Unable to unmarshal, maybe broken body")
		return nil
	}

	if articleContent.Aid == 0 || int64(aid) != articleContent.Aid {
		SendMessage(ctx, 403, "Broken body, missing Aid")
		return nil
	}

	// 数据校验
	if articleContent.Author.Uid == 0 {
		SendMessage(ctx, 403, "Missing uid")
		return nil
	}
	if articleContent.ArticleContent == "" {
		SendMessage(ctx, 403, "Empty content")
		return nil
	}
	articleContent.Likes = 0

	_, err = database.ArticleContentCreate(int64(aid), articleContent)
	if err != nil {
		SendMessage(ctx, 403, "Database input error, this may be a server error")
		return nil
	}

	SendMessage(ctx, 200, "success")
	return nil
}

func ArticleContentQueryHandler(ctx *fiber.Ctx) error {
	aid, err := ctx.ParamsInt("aid")
	if err != nil {
		SendMessage(ctx, 403, "Unable to get aid, maybe broken path")
		return nil
	}
	startAcid, err := ctx.ParamsInt("startAcid")
	if err != nil {
		SendMessage(ctx, 403, "Unable to get acid, maybe broken path")
		return nil
	}
	articleContent, err := database.ArticleContentQuery(int64(aid), int64(startAcid))
	if err != nil {
		SendMessage(ctx, 500, "Database input error, this may be a server error")
		return nil
	}

	successResponse, err := jsoniter.Marshal(struct {
		Code    int64                  `json:"code"`
		Message string                 `json:"message"`
		Data    []types.ArticleContent `json:"data"`
	}{Code: 200, Message: "success", Data: articleContent})

	if err != nil {
		SendMessage(ctx, 500, "Marshal response error, this is a server error")
		return nil
	}

	ctx.Send(successResponse)

	return nil
}
