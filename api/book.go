package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/smallsixFight/hyakkei_blog/logger"
	"github.com/smallsixFight/hyakkei_blog/model"
	"github.com/smallsixFight/hyakkei_blog/service"
	"github.com/smallsixFight/hyakkei_blog/util"
	"net/http"
	"strings"
	"time"
)

func GetBookList(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	page := util.Str2Int(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	list, err := service.GetBooks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	data := make([]model.Book, 0)
	reply.SetData(data)
	if err := json.Unmarshal(list, &data); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析图书数据失败")
		return
	}
	if (page-1)*10 < len(data) {
		last := page * 10
		if last > len(data) {
			last = len(data)
		}
		reply.SetData(data[(page-1)*10 : last])
	}
	reply.SetSuccess(true).SetTotal(int64(len(data))).SetMessage("ok")
}

func SaveBookInfo(ctx *gin.Context) {
	reply := model.Reply{}
	defer ctx.JSON(http.StatusOK, &reply)
	book := model.Book{}
	if err := ctx.BindJSON(&book); err != nil {
		logger.Println("提交参数错误: ", err.Error())
		reply.Message = "提交参数错误"
		return
	}
	if err := handleBookParam(&book); err != nil {
		reply.SetMessage(err.Error())
		return
	}
	book.ModifyAt = time.Now().Format("2006-01-02 15:04:05")
	list, err := service.GetBooks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	books := make([]model.Book, 0)
	if err := json.Unmarshal(list, &books); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析图书列表数据失败")
		return
	}
	// 判断书籍标题是否存在
	if service.BookIsExist(book.Id, book.Title, books) {
		reply.SetMessage("该图书已添加，请修改")
		return
	}
	// 保存
	// id 不为 0 并且书籍存在则只需要更新，不用新增
	var isUpdate bool
	if book.Id != 0 {
		if idx, v := service.FindBookAndIdx(book.Id, books); idx != -1 {
			book.CreateAt = v.CreateAt
			books[idx] = book
			isUpdate = true
		} else {
			reply.SetMessage("更新失败: 没有该图书信息")
			return
		}
	}
	if !isUpdate {
		book.Id = util.GetId()
		book.CreateAt = book.ModifyAt
		books = append(books, book)
	}
	bs, _ := json.Marshal(&books)
	if err := service.SaveBooks(bs); err != nil {
		reply.SetMessage("书籍列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set(model.BookCount, len(books), 0)
	util.Cache.Set("books", bs, 0)

	reply.SetSuccess(true).SetMessage("保存成功")
}

func DeleteBookInfo(ctx *gin.Context) {
	reply := model.Reply{Message: "删除失败，该书籍不存在"}
	defer ctx.JSON(http.StatusOK, &reply)
	id := util.Str2Int64(ctx.Query("id"))
	if id == 0 {
		reply.SetMessage("id 不能为空")
		return
	}
	list, err := service.GetBooks()
	if err != nil {
		reply.SetMessage(err.Error())
		return
	}
	books := make([]model.Book, 0)
	if err := json.Unmarshal(list, &books); err != nil {
		logger.Println(err.Error())
		reply.SetMessage("解析书籍数据失败")
		return
	}
	idx, _ := service.FindBookAndIdx(id, books)
	if idx == -1 {
		return
	}
	books = append(books[:idx], books[idx+1:]...)
	bs, _ := json.Marshal(books)
	if err := service.SaveBooks(bs); err != nil {
		reply.SetMessage("书籍列表保存失败: " + err.Error())
		return
	}
	util.Cache.Set(model.BookCount, len(books), 0)
	util.Cache.Set("books", bs, 0)
	reply.SetMessage("删除成功").SetSuccess(true)
}

func handleBookParam(info *model.Book) error {
	info.Title = strings.TrimSpace(info.Title)
	info.Author = strings.TrimSpace(info.Author)
	info.PicUrl = strings.TrimSpace(info.PicUrl)
	info.DoubanUrl = strings.TrimSpace(info.DoubanUrl)
	info.ShortComment = strings.TrimSpace(info.ShortComment)
	info.Summary = strings.TrimSpace(info.Summary)
	if info.Title == "" {
		return errors.New("书籍名字不能为空")
	} else if info.Author == "" {
		return errors.New("作者名字不能为空")
	} else if info.PicUrl != "" && !util.IsUrl(info.PicUrl) {
		return errors.New("图片链接格式错误")
	} else if info.DoubanUrl != "" && !util.IsUrl(info.DoubanUrl) {
		return errors.New("豆瓣链接格式错误")
	} else if info.Status != model.Wish && info.Status != model.Read {
		return errors.New("书籍状态错误")
	}
	return nil
}
