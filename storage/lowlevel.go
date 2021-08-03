package storage

import "fmt"

func SelectOne(tx ReadTransaction, tbl *table, pred Predicate) (tup *Tuple, err error) {
	if scanErr := tbl.Scan(tx, func(t *Tuple) bool {
		passes, evalErr := pred.Test(t)
		if evalErr != nil {
			err = fmt.Errorf("evaluating predicate: %w", err)
			return false
		}

		if passes {
			tup = t
			return false
		}

		return true
	}); scanErr != nil {
		return nil, scanErr
	}
	return
}
