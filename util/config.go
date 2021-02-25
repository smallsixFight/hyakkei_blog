package util

import (
	"encoding/json"
	"errors"
	"github.com/smallsixFight/hyakkei_blog/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	SysConfigName  = `sys_config.json`
	InitConfigName = `init_config.json`
)

var lock sync.Mutex

type SysSetting struct {
	model.BaseSysSetting
	LogoName string `json:"logo_name"`
	InitTime int64  `json:"init_time"`
	Salt     string `json:"salt,omitempty"`
}

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

var sysConfig *SysSetting

func GetSysConfig() SysSetting {
	if sysConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		f, err := ioutil.ReadFile(filepath.Join(GetBlogConfigPath(), SysConfigName))
		if err != nil {
			log.Fatalln("read config.json failed: " + err.Error())
		}
		if err := json.Unmarshal(f, &sysConfig); err != nil {
			log.Fatalln("parse config failed: " + err.Error())
		}
	}
	return *sysConfig
}

func InitSysConfig(config *SysSetting) error {
	sysConfig = config
	sysConfig.Password, _ = MD5Encrypt(sysConfig.Password, sysConfig.Salt)
	filePath := filepath.Join(GetBlogConfigPath(), SysConfigName)
	d, _ := json.Marshal(&sysConfig)
	if err := WriteFile(d, filePath); err != nil {
		return errors.New("创建系统配置文件失败: " + err.Error())
	}
	return nil
}

func SaveSysConfig(config *SysSetting) error {
	lock.Lock()
	defer lock.Unlock()
	filePath := filepath.Join(GetBlogConfigPath(), SysConfigName)
	if v := strings.TrimSpace(config.Salt); v != "" {
		sysConfig.Salt = v
	}
	if v := strings.TrimSpace(config.BlogName); v != "" {
		sysConfig.BlogName = v
	}
	if v := strings.TrimSpace(config.Username); v != "" {
		sysConfig.Username = v
	}
	if v := strings.TrimSpace(config.Password); v != "" {
		sysConfig.Password, _ = MD5Encrypt(config.Password, sysConfig.Salt)
	}
	if v := strings.TrimSpace(config.GithubName); v != "" {
		sysConfig.GithubName = v
	}
	if v := strings.TrimSpace(config.SavePath); v != "" {
		sysConfig.SavePath = v
	}
	if v := strings.TrimSpace(config.LogoName); v != "" {
		sysConfig.LogoName = v
	}
	if config.InitTime != 0 {
		sysConfig.InitTime = config.InitTime
	}
	if fileInfo, err := os.Stat(config.SavePath); err != nil {
		return errors.New("存储路径错误: " + err.Error())
	} else if !fileInfo.IsDir() {
		return errors.New("存储路径应是一个目录路径")
	}
	sysConfig.IsShowGithub = config.IsShowGithub
	sysConfig.IsShowBook = config.IsShowBook
	d, _ := json.Marshal(&sysConfig)
	if err := WriteFile(d, filePath); err != nil {
		return errors.New("更新系统配置文件失败: " + err.Error())
	}
	return nil
}
