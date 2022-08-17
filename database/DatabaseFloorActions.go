package database

import (
	"errors"
	"github.com/MysGal/Mon3tr/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/xujiajun/nutsdb"
	"strconv"
)

var floorBucket = "floor"

func FloorCreate(did string, floor *types.Floor) error {

	var currentFloor int
	err := GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			size, err := tx.LSize(floorBucket, []byte(did))
			if errors.Is(err, nutsdb.ErrBucket) || errors.Is(err, nutsdb.ErrBucketEmpty) || errors.Is(err, nutsdb.ErrBucketNotFound) || errors.Is(err, nutsdb.ErrBucketNotFound) || size == 0 {
				currentFloor = 1
			} else {
				if err != nil {
					return err
				}
				currentFloor = size + 1
			}

			floor.Floor = int64(currentFloor)

			floorJson, err := jsoniter.Marshal(floor)
			if err != nil {
				return err
			}

			err = tx.RPush(floorBucket, []byte(did), floorJson)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}

	err = GlobalIndex.Index("discussion_"+did+"_"+strconv.Itoa(currentFloor), types.Index{
		Type: "discussion",
		Data: floor,
	})
	if err != nil {
		return err
	}

	return nil
}

func FloorQuery(did string, start int, count int) ([]types.Floor, error) {
	var floors []types.Floor
	err := GlobalDatabase.View(
		func(tx *nutsdb.Tx) error {
			items, err := tx.LRange(floorBucket, []byte(did), start, start+count-1)
			if errors.Is(err, nutsdb.ErrBucket) {
				return nil
			}
			if err != nil {
				return err
			}
			for i, item := range items {
				var floor types.Floor
				err := jsoniter.Unmarshal(item, &floor)
				if err != nil {
					return err
				}
				floor.Floor = int64(start + i + 1)
				floors = append(floors, floor)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	if len(floors) == 0 {
		return nil, err
	}
	return floors, nil
}
