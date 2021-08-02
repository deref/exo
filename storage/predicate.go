package storage

type Predicate interface {
	Test(t *Tuple) (bool, error)
}

func ColumnByIndexEquals(idx int, testVal interface{}) *colEqConst {
	return &colEqConst{
		colIdx:  idx,
		testVal: testVal,
	}
}

func ColumnByNameEquals(schema *Schema, name string, testVal interface{}) *colEqConst {
	idx := -1
	for elemIdx, elem := range schema.Elements {
		if elem.Name == name {
			idx = elemIdx
			break
		}
	}

	if idx == -1 {
		// Should this return an error?
		return nil
	}

	return &colEqConst{
		colIdx:  idx,
		testVal: testVal,
	}
}

type colEqConst struct {
	colIdx  int
	testVal interface{}
}

func (f *colEqConst) Test(t *Tuple) (bool, error) {
	val, err := t.GetDynamic(f.colIdx)
	if err != nil {
		return false, err
	}
	return val == f.testVal, nil
}
