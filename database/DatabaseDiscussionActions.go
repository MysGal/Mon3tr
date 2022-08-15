package database

import (
	"errors"
	"github.com/MysGal/Mon3tr/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/xujiajun/nutsdb"
)

var discussionBucket = "discussion"

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

	err = GlobalDiscussionIndex.Index(did+discussion.Topic, discussion)
	if err != nil {
		return "", err
	}

	return did, nil
}

func DiscussionQuery(topic string, prefix string, start int, count int) ([]types.Discussion, error) {
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
