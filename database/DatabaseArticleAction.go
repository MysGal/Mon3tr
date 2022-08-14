package database

import (
	"github.com/MysGal/Mon3tr/types"
	"github.com/MysGal/Mon3tr/utils"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

func ArticleCreate(articleraw types.Article) error {

	content := articleraw.ArticleContent // 此处为新帖子的第1楼
	articleraw.ArticleContent = nil
	statement, err := GlobalDatabase.Prepare("INSERT INTO article (data) VALUES (?)")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	article, err := jsoniter.Marshal(articleraw)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	result, err := statement.Exec(string(article))
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}
	aid, err := result.LastInsertId()
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	defer statement.Close()

	err = insertArticleConent(aid, content)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	// 写入索引,文章本体均为aid起始
	err = GlobalIndex.Index("aid"+strconv.FormatInt(aid, 10), articleraw)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	return nil
}

func insertArticleConent(aid int64, contentraw []types.ArticleContent) error {
	statement, err := GlobalDatabase.Prepare("INSERT INTO article_content (aid,data) VALUES (?,?)")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}
	for _, v := range contentraw {
		// 进行过滤, 防止储存其他信息
		tempArticleContent := types.ArticleContent{
			ContentFloor:   v.ContentFloor,
			ArticleContent: v.ArticleContent,
			Author: types.User{
				Uid: v.Author.Uid,
			},
			Likes: v.Likes,
		}
		articleContent, err := jsoniter.Marshal(tempArticleContent)
		if err != nil {
			utils.GlobalLogger.Error(err)
			return err
		}
		result, err := statement.Exec(aid, string(articleContent))
		if err != nil {
			utils.GlobalLogger.Error(err)
			return err
		}
		acid, err := result.LastInsertId()
		if err != nil {
			utils.GlobalLogger.Error(err)
			return err
		}

		// 写入索引,文章内容均为acid起始
		err = GlobalIndex.Index("acid"+strconv.FormatInt(acid, 10), tempArticleContent)
		if err != nil {
			utils.GlobalLogger.Error(err)
			return err
		}
	}
	defer statement.Close()
	return nil
}
