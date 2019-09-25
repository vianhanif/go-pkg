package helper_test

import (
	"testing"

	"github.com/vianhanif/go-pkg/sql/helper"
)

func TestGetOrder(t *testing.T) {
	where, _ := helper.BuildFilter([]helper.QueryFilter{
		helper.QueryFilter{Key: "order", Operation: "order", Column: `"created_at"`, Value: `DESC`},
		helper.QueryFilter{Key: "page", Operation: "offset", Value: "1"},
		helper.QueryFilter{Key: "perPage", Operation: "limit", Value: "10"},
	}...)
	if where.OrderBy() == "" {
		t.Fatal("column ordered-by is empty")
	}
}
