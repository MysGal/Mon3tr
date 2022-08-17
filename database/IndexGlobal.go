package database

import (
	"github.com/MysGal/Mon3tr/utils"
	"github.com/blevesearch/bleve/v2"
	gse "github.com/vcaesar/gse-bleve"
)

var GlobalTopicIndex, GlobalDiscussionIndex bleve.Index

func InitIndex() {

	topicIndex, err := bleve.Open("./data/database/index/topic")
	if err != nil {
		opt := gse.Option{
			Index: "./data/database/index/topic",
			Dicts: "./data/tokenizer/dict.small.txt",
			// Dicts: "embed, zh",
			Stop: "",
			Opt:  "search-hmm",
			Trim: "trim",
		}

		topicIndex, err := gse.New(opt)
		if err != nil {
			utils.GlobalLogger.Panic(err)
		}
		GlobalTopicIndex = topicIndex
	} else {
		GlobalTopicIndex = topicIndex
	}

	discussionIndex, err := bleve.Open("./data/database/index/discussion")
	if err != nil {
		opt := gse.Option{
			Index: "./data/database/index/discussion",
			Dicts: "./data/tokenizer/dict.small.txt",
			// Dicts: "embed, zh",
			Stop: "",
			Opt:  "search-hmm",
			Trim: "trim",
		}

		discussionIndex, err := gse.New(opt)
		if err != nil {
			utils.GlobalLogger.Panic(err)
		}
		GlobalDiscussionIndex = discussionIndex
	} else {
		GlobalDiscussionIndex = discussionIndex
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
