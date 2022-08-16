package database

import (
	"errors"
	"github.com/MysGal/Mon3tr/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/xujiajun/nutsdb"
)

var floorBucket = "floor"

func FloorCreate(did string, floor *types.Floor) error {

	floorJson, err := jsoniter.Marshal(floor)
	if err != nil {
		return err
	}

	err = GlobalDatabase.Update(
		func(tx *nutsdb.Tx) error {
			err := tx.RPush(floorBucket, []byte(did), floorJson)
			if err != nil {
				return err
			}
			return nil
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
			items, err := tx.LRange(floorBucket, []byte(did), start, count)
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
		floors = []types.Floor{}
	}
	return floors, nil
}
