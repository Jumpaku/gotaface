package schema_test

import (
	"context"
	"github.com/Jumpaku/gotaface/spanner/schema"
	"github.com/Jumpaku/gotaface/spanner/test"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLister(t *testing.T) {
	test.SkipWithoutDataSource(t)

	admin, client := test.Setup(t)
	test.InitDDLs(t, admin, client.DatabaseName(),
		lo.Flatten(lo.MapToSlice(TestDDLs, func(_ string, ddl []string) []string { return ddl })))

	want := []string{"A", "B_1", "B_2", "B_3", "B_4", "C_1", "C_2", "C_3", "C_4", "C_5", "D_1", "E_1", "E_2", "F_1", "F_2", "F_3", "G", "H", "I", "J", "K"}

	sut := schema.NewLister(client.ReadOnlyTransaction())
	got, err := sut.List(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
