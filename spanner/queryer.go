package spanner

import (
	"context"

	"cloud.google.com/go/spanner"
)

type Queryer interface {
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
}
