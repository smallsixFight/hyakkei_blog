package api

import (
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/middleware"
)

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/api/v1")
	group.POST("/login", Login)              // 登录
	group.GET("/visitor/add", VisitorAdd)    // 访客记录
	group.GET("/articles", FetchArticleList) // 查询发布文章列表
	group.GET("/books", GetBooks)            // 查询书籍数据
	group.GET("/friends", GetFriends)        // 查询友链数据

	authGroup := group.Use(middleware.JWTVerify)
	authGroup.GET("/dashboard/info", GetDashboardInfo)          // 获取仪表盘数据
	authGroup.GET("/article/list", GetArticleList)              // 获取文章列表
	authGroup.GET("/article/detail/:id", GetArticleDetail)      // 获取文章明细信息
	authGroup.DELETE("/article/del", DeleteArticle)             // 删除文章
	authGroup.POST("/article/save", SaveArticle)                // 新增或更新文章
	authGroup.GET("/page/list", GetPageList)                    // 获取文章列表
	authGroup.GET("/page/detail/:id", GetPageDetail)            // 获取文章明细信息
	authGroup.DELETE("/page/del", DeletePage)                   // 删除文章
	authGroup.POST("/page/save", SavePage)                      // 新增或更新文章
	authGroup.GET("/tag/list", GetTagList)                      // 获取标签列表
	authGroup.POST("/tag/save", SaveTag)                        // 新增或更新标签
	authGroup.DELETE("/tag/del", DeleteTag)                     // 删除标签
	authGroup.GET("/friend/list", GetFriendLinkList)            // 友链列表
	authGroup.POST("/friend/save", SaveFriendLink)              // 新增或更新友链
	authGroup.DELETE("/friend/del", DeleteFriendLink)           // 删除友链
	authGroup.GET("/sys_setting/info", GetSysSettingInfo)       // 获取系统设置信息
	authGroup.POST("/sys_setting/update", UpdateSysSettingInfo) // 更新系统设置信息
	authGroup.GET("/book/list", GetBookList)                    // 获取书籍列表
	authGroup.POST("/book/save", SaveBookInfo)                  // 新增或更新阅读书籍
	authGroup.DELETE("/book/del", DeleteBookInfo)               // 删除阅读书籍
	authGroup.POST("/upload/file", UploadAttachment)            // 文件上传
	authGroup.POST("/post/markdown/preview", PreviewMarkdown)   // 预览 markdown 文件
}
