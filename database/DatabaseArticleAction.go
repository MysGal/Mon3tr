package database

import (
	"github.com/MysGal/Mon3tr/types"
	"github.com/MysGal/Mon3tr/utils"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

func ArticleCreate(articleRaw types.Article) (int64, error) {

	content := articleRaw.ArticleContent // 此处为新帖子的第1楼
	articleRaw.ArticleContent = nil

	statement, err := GlobalDatabase.Prepare("INSERT INTO articles (data) VALUES (?)")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	article, err := jsoniter.Marshal(articleRaw)

	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	result, err := statement.Exec(string(article))
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}
	aid, err := result.LastInsertId()
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	defer statement.Close()

	_, err = ArticleContentCreate(aid, content[0])
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	// 写入索引,文章本体均为aid起始
	err = GlobalIndex.Index("aid"+strconv.FormatInt(aid, 10), articleRaw)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	return aid, nil
}

func ArticleUpdate(articleRaw types.Article) error {

	statement, err := GlobalDatabase.Prepare("UPDATE articles SET data=? WHERE aid=?")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}
	aid := articleRaw.Aid

	articleRaw.Aid = 0

	article, err := jsoniter.Marshal(articleRaw)

	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	_, err = statement.Exec(article, aid)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return err
	}

	// 重新索引
	err = GlobalIndex.Index("aid"+strconv.FormatInt(aid, 10), articleRaw)

	return nil
}

func ArticleQuery(aid int64) (types.Article, error) {
	statement, err := GlobalDatabase.Prepare("SELECT data FROM articles WHERE aid=?")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return types.Article{}, err
	}

	var articleJson []byte
	statement.QueryRow(aid).Scan(&articleJson)

	var article types.Article
	jsoniter.Unmarshal(articleJson, &article)

	defer statement.Close()
	return article, nil
}

// 给帖子添加楼
func ArticleContentCreate(aid int64, contentRaw types.ArticleContent) (int64, error) {
	statement, err := GlobalDatabase.Prepare("INSERT INTO article_contents (aid,data) VALUES (?,?)")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}
	// 进行过滤, 防止储存其他信息
	tempArticleContent := types.ArticleContent{
		ArticleContent: contentRaw.ArticleContent,
		Author: types.User{
			Uid: contentRaw.Author.Uid,
		},
		Likes: contentRaw.Likes,
	}
	articleContent, err := jsoniter.Marshal(tempArticleContent)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}
	result, err := statement.Exec(aid, string(articleContent))
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}
	acid, err := result.LastInsertId()
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	// 写入索引,文章内容均为acid起始
	err = GlobalIndex.Index("acid"+strconv.FormatInt(acid, 10), tempArticleContent)
	if err != nil {
		utils.GlobalLogger.Error(err)
		return 0, err
	}

	defer statement.Close()
	return acid, nil
}

func ArticleContentQuery(aid int64, acid int64) ([]types.ArticleContent, error) {
	statement, err := GlobalDatabase.Prepare("SELECT acid,data FROM article_contents WHERE (aid=? AND acid>?)")
	if err != nil {
		utils.GlobalLogger.Error(err)
		return nil, err
	}

	rows, err := statement.Query(aid, acid)

	defer statement.Close()

	var articleContent []types.ArticleContent
	var i int
	for rows.Next() {
		tempArticleContent := types.ArticleContent{}
		var acid int64
		rows.Scan(&acid, &tempArticleContent)
		tempArticleContent.Acid = acid
		tempArticleContent.Aid = aid
		articleContent = append(articleContent, tempArticleContent)
		i++
		if i >= 10 {
			break
		}
	}

	rows.Close()

	return articleContent, nil
}
