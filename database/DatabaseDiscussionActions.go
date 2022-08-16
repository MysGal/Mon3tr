package database

import (
	"errors"
	"github.com/MysGal/Mon3tr/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/xujiajun/nutsdb"
)

var discussionBucket = "discussion"
var didBucket = "did"

// Handler确保Topic、Title、Creator存在，Did在此生成
// 利用prefix完成release和discussion区分
func DiscussionCreate(discussion types.Discussion) (string, error) {
	did := DiscussionNode.Generate().String()

	discussion.Did = did

	discussionJson, err := jsoniter.Marshal(discussion)
	if err != nil {
		return "", err
	}

	var prefix string
	switch discussion.Type {
	case "release":
		prefix = "release"
	case "discussion":
		prefix = "discussion"
	default:
		err = errors.New("unknown discussion type")
		return "", err
	}

	err = GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.LPush(discussionBucket, []byte(prefix+discussion.Topic), discussionJson)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return "", err
	}

	// 将Discussion相关内容根据Did储存，便于查询
	err = GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.Put(didBucket, []byte(did), discussionJson, 0)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return "", err
	}

	// 用Did作为key存入索引，查询的时候直接得到did，检测key分隔符后为0即为主题
	err = GlobalDiscussionIndex.Index(did+"_0", discussion)
	if err != nil {
		return "", err
	}

	return did, nil
}

func DiscussionQueryByTopic(topic string, prefix string, start int, count int) ([]types.Discussion, error) {
	var discussions []types.Discussion
	err := GlobalDatabase.View(
		func(tx *nutsdb.Tx) error {
			items, err := tx.LRange(discussionBucket, []byte(prefix+topic), start, count)
			if errors.Is(err, nutsdb.ErrBucket) {
				return nil
			}
			if err != nil {
				return err
			}
			for _, item := range items {
				var discussion types.Discussion
				err := jsoniter.Unmarshal(item, &discussion)
				if err != nil {
					return err
				}
				discussions = append(discussions, discussion)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	if len(discussions) == 0 {
		discussions = []types.Discussion{}
	}
	return discussions, nil
}

func DiscussionQueryByDid(did string) (types.Discussion, error) {
	var discussion types.Discussion
	err := GlobalDatabase.View(
		func(tx *nutsdb.Tx) error {
			items, err := tx.Get(didBucket, []byte(did))
			if errors.Is(err, nutsdb.ErrBucket) || errors.Is(err, nutsdb.ErrBucketNotFound) || errors.Is(err, nutsdb.ErrKeyNotFound) {
				return nil
			}
			if err != nil {
				return err
			}
			err = jsoniter.Unmarshal(items.Value, &discussion)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return types.Discussion{}, err
	}

	return discussion, nil

}
