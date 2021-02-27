package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const InitConfigName = `init_config.json`

var lock sync.Mutex

type InitConfig struct {
	JsPath   string `json:"js_path"`
	ImgPath  string `json:"img_path"`
	CssPath  string `json:"css_path"`
	FontPath string `json:"font_path"`
}

var initConfig *InitConfig

func GetInitConfig() *InitConfig {
	if initConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		f, err := ioutil.ReadFile(filepath.Join(GetBlogConfigPath(), InitConfigName))
		if err != nil {
			log.Fatalln("read config.json failed: " + err.Error())
		}
		if err := json.Unmarshal(f, &initConfig); err != nil {
			log.Fatalln("parse config failed: " + err.Error())
		}
	}
	return initConfig
}

func SaveInitConfig(config *InitConfig) error {
	d, _ := json.Marshal(&config)
	filePath := filepath.Join(GetBlogConfigPath(), InitConfigName)
	f, err := os.Create(filePath)
	if err != nil {
		return errors.New("创建初始化配置文件失败: " + err.Error())
	}
	defer f.Close()
	if _, err = f.Write(d); err != nil {
		return errors.New("写入初始化配置文件失败: " + err.Error())
	}
	initConfig = config
	return nil
}
