package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/api"
	"github.com/smallsixFight/hyakkei_blog/blog_install"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/middleware"
	"github.com/smallsixFight/hyakkei_blog/task/ticker"
	"log"
)

var install = flag.Bool("install", false, "true|false")
var mode = flag.String("mode", "release", "release|debug")
var port = flag.Int("listen_port", 9900, "1024-65535")

func main() {
	flag.Parse()
	// 日志配置
	logger.Init(logger.LowestLevel(logger.InfoLevel))

	if *install {
		fmt.Println("该操作将重新生成配置信息及静态资源文件，不会对已生成的数据做任何操作。")
		fmt.Println("注意：如修改了文件存储路径，那么新生成的文件将发生变化，请注意对先前生成的文件进行处理。")
	}
	if *install || blog_install.IsNeedInstall() {
		if err := blog_install.Init(); err != nil {
			log.Fatal(err.Error())
		}
	}
	if err := blog_install.SetRunCache(); err != nil {
		log.Fatal("设置初始化缓存数据失败: " + err.Error())
	}
	ticker.RunTickerTask()
	fmt.Println("初始化成功，启动服务...")
	// 启动服务
	gin.SetMode(*mode)
	engine := gin.Default()
	engine.Use(middleware.Cors) // 解决跨域
	api.RegisterRoute(engine)
	if err := engine.Run(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatal("服务启动失败: ", err.Error())
	}
}
