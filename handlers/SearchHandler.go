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
	"strings"
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
	rawQuery := ctx.Params("query")
	if rawQuery == ":query" {
		SendMessage(ctx, 403, "broken path")
		return nil
	}
	rawQuery, err := url.PathUnescape(rawQuery)
	if err != nil {
		SendMessage(ctx, 403, "broken path")
		return nil
	}
	from, err := ctx.ParamsInt("from")
	if err != nil {
		SendMessage(ctx, 403, "broken path")
		return nil
	}

	// 处理query语法
	// 检查是否有语法，没有的话直接按照全局搜索处理
	var query string
	queryTypes := []string{"discussion", "topic"}
	querySlice := strings.Split(rawQuery, " ")
	if len(querySlice) == 1 {
		query = rawQuery
	} else {
		queryGalGame := false
		for _, str := range querySlice {
			// 检查是否为类型限制语法
			if strings.HasPrefix(str, "type:") {
				// 比较类型是否为现有类型
				switch str[5:] {
				case "galgame":
					queryTypes = []string{"topic"}
					query += " +data.type:galgame"
					queryGalGame = true
				default:
					queryTypes = []string{"discussion", "topic"}
				}
				continue
			}

			// 检查galgame相关要求
			if queryGalGame {
				if strings.HasPrefix(str, "author:") {
					query += " +data.related_data.galgame_author:" + str[7:]
					continue
				}
				if strings.HasPrefix(str, "publisher:") {
					query += " +data.related_data.galgame_publisher:" + str[10:]
					continue
				}
			}

			query += " " + str
		}
	}

	utils.GlobalLogger.Info(query)
	utils.GlobalLogger.Info(queryTypes)

	var returnData returnStruct

	// 全局搜索，分主题和帖子分别搜索
	returnDataQueried, err := queryBySection(query, queryTypes, from)
	returnData = returnDataQueried
	if err != nil {
		SendMessage(ctx, 500, "server error")
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

func queryBySection(query string, queryTypes []string, from int) (returnStruct, error) {
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
		query := query + " +type:" + v

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
