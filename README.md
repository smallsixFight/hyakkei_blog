# Hyakkei

--- 

### 概述
hyakkei 是一个基于 Golang 构建的博客系统，使用 markdown 作为书写格式，并自定义了一些 markdown 标签来给内容添加音乐以及图片预览方式。

在 hyakkei 中的文章和自定义页会在发布时基于模板自动生成静态页面（.html），列表数据的获取基于 AJAX 进行异步查询后插入。

本博客还有许多不完善的地方，欢迎提交 PR，如果你有问题和建议，请在 issue 留言。

### 为什么写这个系统
我在从大学出来实习时，由于工作导致个人的时间变得碎片化，并且工作上带来的琐碎事也挺多的，然后就想有个固定的地方写东西记东西，大学舍友说博客啊，我就跟着弄博客了。

我一开始用了 WorkPress，后面自己写了一个叫 colorful 的博客系统，写得很臃肿，也不好用。

我还是想写一个简洁的，能够持续优化，不需要在从头来过的博客系统。我还是想写一个简洁的，能够持续优化，不需要在从头来过的博客系统。我接触的博客系统不多，自己也在思考，一个个人博客系统需要的是什么。纯静态？动态？需要服务器？需要数据库？

思考后有了下面几点模糊的观点：

- 我有个小服务器，所以我写个可以放在自己服务器，不用托管的，好掌控。
- 文章基本上变动不大，可以直接生成静态页面。
- 后台管理数据操作较为频繁，所以使用动态获取插入就好，这样就不用频繁重新生成静态文件。
- 不用因一点点使用就用一个大的外部依赖，如 colorful 使用的 Vue 和 elementUI，其实都可以不要。 
- 并不要数据库，数据用一种好用的格式存在文件里就好了，然后进行缓存，这样也就不会产生大量 I/O。
- 分类和标签其实没啥大的区别，至少对于个人博客来讲，所以不需要同时具备标签和分类。
- 要有一个跟系统符合，并且简单的 markdown 编辑器。

机缘巧合的是，我发现我一直关注的一位博主 [imalan](https://blog.imalan.cn/) 重新改了博客主题，很符合我的心水，于是我的重写博客系统之路又开始了。

关于 hyakkei 这个词的由来，这个词应该是`百景`的日文的罗马音，而百景这个词取自歌曲`笑颜百景`，我很喜欢这首歌。

### 示例
[markdown_demo](./web_demo/markdown_demo.md) 是 hyakkei 支持的标签示例，实际效果为 [hyakkei 博客模板示例](https://blog.lamlake.com/hyakkei-blog-example.html) 。

### 目标与功能
- [x] 为文章和自定义页生成 `.html` 静态文件；
- [x] 具有后台管理功能；
- [x] 简单的 markdown 编辑器，能够预览，使用自定义一些 markdown 标签；
- [x] 不使用数据库，使用 json 文件存储数据，并自己写个简单缓存加快访问速度以及减少 I/O；
- [x] 可以在后台管理修改博客名称、用户名、登录密码，可以快速 favicon.ico、logo 图片；
- [ ] 评论功能（准备使用 Disqus）；
- [ ] 数据分页，防止数据量过大导致存储文件过大；

### 安装与使用
本项目基于 golang 编写，所以原则上不需要安装任何依赖环境，只需要可执行的二进制文件和 `templates` 这个模板文件即可开始安装。

但是，我觉得项目还没足够文件，暂不上传对应的文件下载，所以还请想尝试的各位可以先安装一下 [goalng]() ，然后下载源码进行编译执行。

下载源码：
```shell
git clone https://github.com/smallsixFight/hyakkei_blog && cd hyakkei_blog
```
编译生成可执行文件：
```shell
# 请注意：如果你的运行的系统与当前进行编译的系统不同，如编译的系统为 Window/macOS，服务器为 linux 系统的服务器，需要先执行 `set GOOS=linux` 或 `export GOOS=linux`，之后再进行编译。
go build -o hyakkei_blog main.go
```

安装：
```shell
./hyakkei_blog
# 重新生成配置信息可使用 ./hyakkei_blog -install 进入配置设置。
```

启动时，程序会检验相应的文件是否完整，不完整则会自动创建；系统还会依此要求你设置博客名称、登录账户、登录密码、加密的 salt、是否显示 Github 链接（以及 Github 名称）、是否显示书籍页以及静态文件存储路径。

这里需要注意的是静态文件的存储路径，也就是静态页面生成保存的位置，比如我使用了 nginx，那么我的路径就可以设置为 `/usr/local/nginx/html`。

测试：
在安装成功并且服务成功启动后，接着就可以直接打开你设置的静态文件的存储路径，用浏览器打开 index.html 文件打开博客首页，打开 `静态文件存储路径/hyakkei/bs/login.html`，尝试登录后台管理。如果以上操作成功了，那么恭喜你，hyakkei blog 已成功运行了。


### 致谢
- 感谢 [失眠海峡](https://blog.imalan.cn/) 让我扒了很多样式和一些方案，如果多图并列，图片查看等；
- 感谢 [Maverick](https://github.com/AlanDecode/Maverick) ，这个是 `失眠海峡` 博主自己写的静态博客生成器，学习到了很多思路；
- 感谢 [MetingJS](https://github.com/metowolf/MetingJS) 和 [APlayer](https://github.com/DIYgod/APlayer) 提供了音乐播放插件；
- 感谢 [PhotoSwipe](https://github.com/dimsemenov/PhotoSwipe) 提供了图片查看插件；
- 感谢 [highlightJS](https://github.com/highlightjs/highlight.js) 提供了代码高亮插件；
- 感谢 [jquery](https://github.com/jquery/jquery) ;
- 感谢 [goldmark](https://github.com/yuin/goldmark) ，该项目能够十分友好地进行 markdown 解析，并且能够很方便的进行 markdown 标签的扩展解析或者渲染，以及编写自定义的 markdown 标签。
- 感谢 [gin](https://github.com/gin-gonic/gin) ，一个轻量的 golang Web 框架；
- 感谢 [jwt-go](https://github.com/dgrijalva/jwt-go) ，用于生成 JWT 验证。