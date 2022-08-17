package handlers

import (
	"github.com/MysGal/Mon3tr/database"
	"github.com/MysGal/Mon3tr/types"
	"github.com/MysGal/Mon3tr/utils"
	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strconv"
)

type returnStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Topics []struct {
			Topic types.Topic `json:"topic"`
			Score float64     `json:"score"`
		} `json:"topics,omitempty"`
		Discussions []struct {
			Discussion types.Discussion `json:"discussion"`
			Score      float64          `json:"score"`
		} `json:"discussions,omitempty"`
	} `json:"data"`
}

func SearchHandler(ctx *fiber.Ctx) error {
	searchType := ctx.Params("type")

	keyWord, err := url.PathUnescape(ctx.Params("keyword"))

	if err != nil || searchType == ":type" || searchType == "" || keyWord == "" || keyWord == ":keyword" {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	from, err := ctx.ParamsInt("from")

	if err != nil {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	var returnData returnStruct

	switch searchType {
	case "global":
		// 全局搜索，分主题和帖子分别搜索

		queryTypes := []string{"discussion", "topic"}
		returnDataQueried, err := queryBySection(queryTypes, keyWord, from)
		returnData = returnDataQueried
		if err != nil {
			SendMessage(ctx, 500, "server error")
			return nil
		}
	default:
		SendMessage(ctx, 403, "unknow search type")
		return nil
	}

	returnBody, err := jsoniter.Marshal(returnData)
	if err != nil {
		SendMessage(ctx, 500, "server marshal error")
		return nil
	}

	ctx.Send(returnBody)

	return nil
}

func queryBySection(queryTypes []string, keyWord string, from int) (returnStruct, error) {
	returnData := returnStruct{
		Code:    200,
		Message: "success",
		Data: struct {
			Topics []struct {
				Topic types.Topic `json:"topic"`
				Score float64     `json:"score"`
			} `json:"topics,omitempty"`
			Discussions []struct {
				Discussion types.Discussion `json:"discussion"`
				Score      float64          `json:"score"`
			} `json:"discussions,omitempty"`
		}{Topics: make([]struct {
			Topic types.Topic `json:"topic"`
			Score float64     `json:"score"`
		}, 0), Discussions: make([]struct {
			Discussion types.Discussion `json:"discussion"`
			Score      float64          `json:"score"`
		}, 0)},
	}

	for _, v := range queryTypes {
		query := keyWord + " +type:" + v

		queryReq := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
		queryReq.Size = 10
		queryReq.From = from
		// 获取type在进行处理
		//queryReq.Fields = []string{"type"}
		queryReq.SortBy([]string{"_score"})
		queryRes, err := database.GlobalIndex.Search(queryReq)
		if err != nil {
			return returnStruct{}, err
		}
		if queryRes.Total != 0 {
			for _, hit := range queryRes.Hits {
				switch v {
				case "topic":
					// topic_
					topic, err := database.TopicQueryDetail(hit.ID[6:])
					if err != nil {
						utils.GlobalLogger.Error(err)
						continue
					}
					score := hit.Score
					returnTopic := struct {
						Topic types.Topic `json:"topic"`
						Score float64     `json:"score"`
					}{Topic: topic, Score: score}
					returnData.Data.Topics = append(returnData.Data.Topics, returnTopic)
				case "discussion":
					discussion, err := database.DiscussionQueryByDid(hit.ID[11:30])
					floorNumber, err := strconv.Atoi(hit.ID[31:])
					if err != nil {
						utils.GlobalLogger.Error(err)
						continue
					}

					var floor types.Floor
					if floorNumber != 0 {
						floors, err := database.FloorQuery(discussion.Did, floorNumber-1, 1)
						if err != nil {
							utils.GlobalLogger.Error(err)
							utils.GlobalLogger.Info(floorNumber)
							continue
						}
						floor = floors[0]
					} else {
						continue
					}
					// 补充数据
					discussion.Floor = &floor
					score := hit.Score
					returnDiscussion := struct {
						Discussion types.Discussion `json:"discussion"`
						Score      float64          `json:"score"`
					}{Discussion: discussion, Score: score}
					returnData.Data.Discussions = append(returnData.Data.Discussions, returnDiscussion)
				}
			}
		}
	}

	return returnData, nil
}
