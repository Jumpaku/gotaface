package schema

import "context"

type Lister interface {
	List(ctx context.Context) (tables []string, err error)
}
