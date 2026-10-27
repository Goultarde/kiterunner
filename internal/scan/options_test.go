package scan

import (
	"bytes"
	"testing"

	"github.com/assetnote/kiterunner/pkg/http"
)

func TestScanOptionsValidateAppliesForceMethod(t *testing.T) {
	route := &http.Route{Method: http.GET}
	opts := NewDefaultScanOptions()
	opts.Routes = []*http.Route{route}
	opts.SuccessStatusCodes = map[int]interface{}{200: struct{}{}}
	opts.ForceMethod = "POST"

	if err := opts.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !bytes.Equal(route.Method, http.POST) {
		t.Fatalf("route method = %v, want POST", route.Method)
	}
}
