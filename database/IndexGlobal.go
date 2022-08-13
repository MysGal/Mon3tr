package database

import (
	"github.com/blevesearch/bleve/v2"
	gsebleve "github.com/leopku/bleve-gse-tokenizer/v2"
	"log"
)

var GlobalIndex bleve.Index

func InitIndex() {
	index, err := bleve.Open("./data/database/index")
	if err != nil {
		mapping := bleve.NewIndexMapping()
		if err := mapping.AddCustomTokenizer("gse", map[string]interface{}{
			"type":       gsebleve.Name,
			"user_dicts": "./data/tokenizer/dict.txt",
		}); err != nil {
			panic(err)
		}
		if err := mapping.AddCustomAnalyzer("gse", map[string]interface{}{
			"type":      "gse",
			"tokenizer": "gse",
		}); err != nil {
			panic(err)
		}
		mapping.DefaultAnalyzer = "gse"

		index, err = bleve.New("./data/database/index", mapping)
		if err != nil {
			log.Fatal(err)
		}
	}

	GlobalIndex = index

}

//func Test() {
//	article := types.Article{1, "test", "testcontent", types.GalGame{1, "命123sa运asd1u29之casca夜", []string{"testauthor1", "testauthor2"}, []string{"testpublisher1", "testpublisher2"}, []string{"tag1", "tag2"}}}
//
//	//message := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工"
//	if err := GlobalIndex.Index(strconv.Itoa(article.Aid), article); err != nil {
//		panic(err)
//	}
//	query := "+related_galgame.name:命运之夜"
//	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
//	req.Highlight = bleve.NewHighlight()
//	res, err := GlobalIndex.Search(req)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Result of: '%s': %d matches\n", query, res.Total)
//	for i, hit := range res.Hits {
//		rv := fmt.Sprintf("%d. %s, (%f)\n", i+res.Request.From+1, hit.ID, hit.Score)
//		for fragmentField, fragments := range hit.Fragments {
//			rv += fmt.Sprintf("%s: ", fragmentField)
//			for _, fragment := range fragments {
//				rv += fmt.Sprintf("%s", fragment)
//			}
//		}
//		fmt.Printf("%s\n", rv)
//	}
//}
