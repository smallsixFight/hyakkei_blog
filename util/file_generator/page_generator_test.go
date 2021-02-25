package file_generator

import "testing"

func TestGenerateHeaderFile(t *testing.T) {
	if err := GenerateHeaderFile(); err != nil {
		t.Fatal(err.Error())
	}
	t.Log("success")
}
