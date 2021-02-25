package util

import (
	"github.com/smallsixFight/hyakkei_blog/model"
	"testing"
)

func TestSaveSysConfig(t *testing.T) {
	cfg := SysSetting{
		BaseSysSetting: model.BaseSysSetting{
			BlogName:   "hyakkei",
			Username:   "hyakkei",
			GithubName: "https://github.com/smallsixFight",
		},
		Salt: "hyakkei",
	}
	if err := SaveSysConfig(&cfg); err != nil {
		t.Fatal(err.Error())
	}
	t.Log("success")
}
