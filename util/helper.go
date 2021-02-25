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

func ReplaceBasePath(d []byte) string {
	str := string(d)
	str = strings.ReplaceAll(str, "#{{img_path}}", GetInitConfig().ImgPath)
	str = strings.ReplaceAll(str, "#{{js_path}}", GetInitConfig().JsPath)
	str = strings.ReplaceAll(str, "#{{css_path}}", GetInitConfig().CssPath)
	str = strings.ReplaceAll(str, "#{{font_path}}", GetInitConfig().FontPath)
	str = strings.Replace(str, "#{{github_name}}", GetSysConfig().GithubName, 1)
	str = strings.Replace(str, "#{{blog_name}}", GetSysConfig().BlogName, 1)
	str = strings.Replace(str, "#{{pswp}}", pswpStr, 1)
	return str
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

const pswpStr = `<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">
        <!-- Background of PhotoSwipe.
         It's a separate element as animating opacity is faster than rgba(). -->
        <div class="pswp__bg"></div>
        <div class="pswp__scroll-wrap">
            <div class="pswp__container">
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
                <div class="pswp__item"></div>
            </div>
            <!-- Default (PhotoSwipeUI_Default) interface on top of sliding area. Can be changed. -->
            <div class="pswp__ui pswp__ui--hidden">
                <div class="pswp__top-bar">
                    <div class="pswp__counter"></div>
                    <button class="pswp__button pswp__button--close" title="Close (Esc)"></button>
                    <!-- <button class="pswp__button pswp__button--share" title="Share"></button> -->
                    <button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>
                    <button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>
                    <!-- Preloader demo https://codepen.io/dimsemenov/pen/yyBWoR -->
                    <!-- element will get class pswp__preloader--active when preloader is running -->
                    <div class="pswp__preloader">
                        <div class="pswp__preloader__icn">
                            <div class="pswp__preloader__cut">
                                <div class="pswp__preloader__donut"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
                    <div class="pswp__share-tooltip"></div>
                </div>

                <button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)">
                </button>

                <button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)">
                </button>

                <div class="pswp__caption">
                    <div class="pswp__caption__center"></div>
                </div>
            </div>
        </div>
    </div>`
