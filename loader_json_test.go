package openapi

import "testing"

func TestLoadFromDataJSON_Error(t *testing.T) {
	want := `unexpected end of JSON input`
	if _, err := LoadFromDataJSON([]byte(`{`)); err == nil && err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
