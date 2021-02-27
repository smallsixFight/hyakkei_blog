package service

import (
	"encoding/json"
	"errors"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var sysConfig *model.SysSetting

var settingLock sync.Mutex

const SysConfigName = `sys_config.json`

func InitSysConfig() error {
	f, err := ioutil.ReadFile(filepath.Join(util.GetBlogConfigPath(), SysConfigName))
	if err != nil {
		return errors.New("read config file failed: " + err.Error())
	}
	if err := json.Unmarshal(f, &sysConfig); err != nil {
		return errors.New("parse config params failed: " + err.Error())
	}
	return nil
}

func GetSysConfig() model.SysSetting {
	if sysConfig == nil {
		settingLock.Lock()
		defer settingLock.Unlock()
		if err := InitSysConfig(); err != nil {
			logger.Error(err.Error())
		}
	}
	return *sysConfig
}

func GenerateSysSetting(config *model.SysSetting) error {
	sysConfig = config
	sysConfig.Password, _ = util.MD5Encrypt(sysConfig.Password, sysConfig.Salt)
	filePath := filepath.Join(util.GetBlogConfigPath(), SysConfigName)
	d, _ := json.Marshal(&sysConfig)
	if err := util.WriteFile(d, filePath); err != nil {
		return errors.New("创建系统配置文件失败: " + err.Error())
	}
	return nil
}

func SaveSysConfig(config *model.SysSetting) error {
	settingLock.Lock()
	defer settingLock.Unlock()
	filePath := filepath.Join(util.GetBlogConfigPath(), SysConfigName)
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
		sysConfig.Password, _ = util.MD5Encrypt(config.Password, sysConfig.Salt)
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
	if err := util.WriteFile(d, filePath); err != nil {
		return errors.New("更新系统配置文件失败: " + err.Error())
	}
	return nil
}
