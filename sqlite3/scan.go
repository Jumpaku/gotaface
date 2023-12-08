package sqlite3

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ScanRowsStruct[RowStruct any](itr *sqlx.Rows) ([]RowStruct, error) {
	rowsStruct := []RowStruct{}
	for itr.Next() {
		var row RowStruct
		if err := itr.StructScan(&row); err != nil {
			return nil, fmt.Errorf(`fail to scan row: %v`, err)
		}

		rowsStruct = append(rowsStruct, row)
	}
	return rowsStruct, nil
}
