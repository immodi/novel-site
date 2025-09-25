package handlers_test

import (
	"github.com/a-h/templ"
	"immodi/novel-site/internal/app/handlers"
	indexdtostructs "immodi/novel-site/internal/http/structs/index"
	"net/http"
	"testing"
)

func TestGenericHandler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w          http.ResponseWriter
		r          *http.Request
		data       *indexdtostructs.MetaDataStruct
		cmp        templ.Component
		statusCode []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.GenericHandler(tt.w, tt.r, tt.data, tt.cmp, tt.statusCode)
		})
	}
}
