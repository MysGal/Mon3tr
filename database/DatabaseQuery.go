package database

// 全局搜索，搜索一切匹配的
//func GlobalQuery(str string) error {
//	req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(str))
//	res, err := GlobalIndex.Search(req)
//	if err != nil {
//		utils.GlobalLogger.Error(err)
//		return err
//	}
//
//	for _, hit := range res.Hits {
//		if strings.HasPrefix(hit.ID, "aid") {
//			// article主体，即搜索到标题/tags/相关GalGame
//		} else if strings.HasPrefix(hit.ID, "acid") {
//			// 帖子中楼，获取完Aid后返回
//		}
//	}
//}
