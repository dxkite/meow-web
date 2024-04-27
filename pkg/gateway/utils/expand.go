package utils

type ExpandHandler interface {
	Get(i int) ([]uint64, error)
	Set(i int, v []interface{}) error
	BatchGet(id []uint64) (map[uint64]interface{}, error)
}

func ExpandStruct(n int, h ExpandHandler) error {
	idArr := []uint64{}
	idMap := map[uint64]struct{}{}

	for i := 0; i < n; i++ {
		ids, err := h.Get(i)
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

	idObjects, err := h.BatchGet(idArr)
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		ids, err := h.Get(i)
		if err != nil {
			return err
		}
		values := []interface{}{}
		for _, id := range ids {
			if vv, ok := idObjects[id]; ok {
				values = append(values, vv)
			}
		}
		if err := h.Set(i, values); err != nil {
			return err
		}
	}
	return nil
}
