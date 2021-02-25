package util

import (
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	if err := CopyDir(filepath.Join(GetBlogTemplatePath(), "assets"), filepath.Join(GetAbsPath(), "hyakkei")); err != nil {
		t.Fatal(err.Error())
	}
	t.Log("success")
}
