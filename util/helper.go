package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var httpUrlRegexp = regexp.MustCompile(`^(?:http|https)://[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]+(?::\d+)?(?:[/#?][-a-zA-Z0-9@:%_+.~#$!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`)

func GetAbsPath() string {
	v, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return v
}

func GetBlogDataPath() string {
	return filepath.Join(GetAbsPath(), "hyakkei", "data")
}

func GetBlogConfigPath() string {
	return filepath.Join(GetAbsPath(), "config")
}

func GetBlogTemplatePath() string {
	return filepath.Join(GetAbsPath(), "templates")
}

func Str2Int(s string) int {
	v, _ := strconv.ParseInt(s, 10, 64)
	return int(v)
}

func Str2Int32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return int32(v)
}

func Str2Int64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func CopyDir(src, dest string) error {
	if srcInfo, err := os.Stat(src); err != nil {
		return errors.New("源路径错误：" + err.Error())
	} else if !srcInfo.IsDir() {
		return errors.New("源路径不是一个文件夹")
	}
	if srcInfo, err := os.Stat(dest); err != nil {
		return errors.New("目标路径错误：" + err.Error())
	} else if !srcInfo.IsDir() {
		return errors.New("目标路径不是一个文件夹")
	}
	// 获取要复制的文件夹名字
	arr := strings.Split(src, string(filepath.Separator))
	baseDirName := arr[len(arr)-1]

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			log.Printf("[%s]路径出现错误: %s", path, err.Error())
			return filepath.SkipDir
		}
		fileName := getDirName(src, path)
		p := filepath.Join(dest, baseDirName, fileName)
		fmt.Printf("正在创建[%s]\n", p)
		if info.IsDir() {
			if !FileIsExist(p) {
				if err := os.Mkdir(p, os.ModePerm); err != nil {
					log.Println(err.Error())
				}
			}
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.New("读取文件失败: " + err.Error())
			}
			if err := WriteFile(data, p); err != nil {
				return errors.New("写入文件失败: " + err.Error())
			}
		}
		fmt.Println("创建成功")
		return nil
	})
	return err
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	return false
}

func getDirName(basePath, s string) string {
	arr := strings.Split(s, basePath)
	return arr[len(arr)-1]
}

func IsUrl(val string) bool {
	return httpUrlRegexp.MatchString(val)
}

var idLock sync.Mutex

func GetId() int64 {
	idLock.Lock()
	defer idLock.Unlock()
	id := time.Now().UnixNano() / 1000
	time.Sleep(time.Nanosecond * 2)
	return id
}

// 生成文件并写入数据
func WriteFile(data []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}
