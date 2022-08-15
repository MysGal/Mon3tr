package database

import (
	"errors"
	"github.com/MysGal/Mon3tr/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/xujiajun/nutsdb"
)

var topicBucket = "topic"

func TopicCreate(topic types.Topic) error {
	// 检测该主题是否已存在
	_, err := TopicQueryDetail(topic.Topic)
	if !errors.Is(err, nutsdb.ErrBucketNotFound) && !errors.Is(err, nutsdb.ErrKeyNotFound) {
		err = errors.New("topic already exist")
		return err
	}

	topicJson, err := jsoniter.Marshal(topic)
	if err != nil {
		return err
	}

	err = GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.Put(topicBucket, []byte(topic.Topic), topicJson, 0)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}

	// 加入索引
	err = GlobalTopicIndex.Index(topic.Topic, topic)
	if err != nil {
		return err
	}
	return nil
}

// TODO:考虑和TopicCreate合并
func TopicUpdate(topic types.Topic) error {
	topicJson, err := jsoniter.Marshal(topic)
	if err != nil {
		return err
	}
	err = GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.Put(topicBucket, []byte(topic.Topic), topicJson, 0)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	// 重新索引
	err = GlobalTopicIndex.Index(topic.Topic, topic)
	if err != nil {
		return err
	}
	return nil
}

// 查询一个主题的名字、Tags
func TopicQueryDetail(topic string) (types.Topic, error) {
	// 根据Topic的Topic来查询，无需其他部分
	var queriedTopic types.Topic
	err := GlobalDatabase.View(
		func(tx *nutsdb.Tx) error {
			entry, err := tx.Get(topicBucket, []byte(topic))
			if err != nil {
				return err
			}

			err = jsoniter.Unmarshal(entry.Value, &queriedTopic)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return types.Topic{}, err
	}
	return queriedTopic, nil
}

// 查询所有的主题并返回
func TopicQueryAll() ([]types.Topic, error) {
	var queriedTopic []types.Topic
	err := GlobalDatabase.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(topicBucket)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				var tempTopic types.Topic
				err = jsoniter.Unmarshal(entry.Value, &tempTopic)
				if err != nil {
					return err
				}

				queriedTopic = append(queriedTopic, tempTopic)
			}

			return nil
		})

	if errors.Is(err, nutsdb.ErrBucketEmpty) {
		return []types.Topic{}, nil
	}

	if err != nil {
		return nil, err
	}
	return queriedTopic, nil
}
