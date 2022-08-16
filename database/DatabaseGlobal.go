package database

import (
	"github.com/MysGal/Mon3tr/utils"
	"github.com/bwmarrin/snowflake"
	"github.com/xujiajun/nutsdb"
)

var GlobalDatabase *nutsdb.DB

var DiscussionNode, UserNode *snowflake.Node

func InitDatabase() {
	// 数据库初始化
	databaseOptions := nutsdb.Options{
		EntryIdxMode:         nutsdb.HintKeyValAndRAMIdxMode,
		SegmentSize:          8 * 1024 * 1024,
		NodeNum:              1,
		RWMode:               nutsdb.FileIO,
		SyncEnable:           true,
		StartFileLoadingMode: nutsdb.MMap,
	}
	database, err := nutsdb.Open(
		databaseOptions,
		nutsdb.WithDir("./data/database/nuts"),
	)
	if err != nil {
		utils.GlobalLogger.Panic(err)
	}
	GlobalDatabase = database

	// ID生成器初始化
	// 用户ID生成器
	userNode, err := snowflake.NewNode(1)
	if err != nil {
		utils.GlobalLogger.Panic(err)
	}
	UserNode = userNode

	// 帖子ID生成器
	discussionNode, err := snowflake.NewNode(2)
	if err != nil {
		utils.GlobalLogger.Panic(err)
	}
	DiscussionNode = discussionNode

}
