package postgres

import (
	"github.com/jackc/pgx/v5"
)

func ScanRowsStruct[RowStruct any](rows pgx.Rows) ([]RowStruct, error) {
	return pgx.CollectRows[RowStruct](rows, pgx.RowToStructByNameLax[RowStruct])
}
