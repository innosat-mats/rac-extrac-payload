package exports

import "testing"

func TestGetJSONFilename(t *testing.T) {
	templateName := "my/file.png"
	want := "my/file.json"
	if got := GetJSONFilename(templateName); got != want {
		t.Errorf("GetJSONFilename() = %v, want %v", got, want)
	}
}
