package utils

func ExpandStruct(n int, get func(i int) ([]uint64, error), set func(i int, v []interface{}) error, batchGet func(ids []uint64) (map[uint64]interface{}, error)) error {
	idArr := []uint64{}
	idMap := map[uint64]struct{}{}

	for i := 0; i < n; i++ {
		ids, err := get(i)
		if err != nil {
			return err
		}
		for _, id := range ids {
			if _, ok := idMap[id]; !ok {
				idArr = append(idArr, id)
				idMap[id] = struct{}{}
			}
		}
	}

	idObjects, err := batchGet(idArr)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		ids, err := get(i)
		if err != nil {
			return err
		}
		values := []interface{}{}
		for _, id := range ids {
			if vv, ok := idObjects[id]; ok {
				values = append(values, vv)
			}
		}
		if err := set(i, values); err != nil {
			return err
		}
	}
	return nil
}
