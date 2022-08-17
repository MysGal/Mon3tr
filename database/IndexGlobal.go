package database

import (
	"github.com/MysGal/Mon3tr/utils"
	"github.com/blevesearch/bleve/v2"
	gse "github.com/vcaesar/gse-bleve"
)

/*
全都以Index结构入索引
*/

var (
	GlobalIndex bleve.Index
)

func InitIndex() {

	index, err := bleve.Open("./data/database/index/discussion")
	if err != nil {
		opt := gse.Option{
			Index: "./data/database/index/discussion",
			Dicts: "./data/tokenizer/dict.small.txt",
			Stop:  "./data/tokenizer/stop.txt",
			Opt:   "search-hmm",
			Trim:  "trim",
		}

		index, err := gse.New(opt)
		if err != nil {
			utils.GlobalLogger.Panic(err)
		}
		GlobalIndex = index
	} else {
		GlobalIndex = index
	}
}

//func Test() {
//	query := "testgal"
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
