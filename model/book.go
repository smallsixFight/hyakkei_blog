package model

type ReadType string

const (
	Wish ReadType = "wish"
	Read ReadType = "read"
)

type Book struct {
	Id           int64    `json:"id,string"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	Author       string   `json:"author"`
	DoubanUrl    string   `json:"douban_url"`
	PicUrl       string   `json:"pic_url"`
	Status       ReadType `json:"status"`
	CreateAt     string   `json:"create_at"`
	ModifyAt     string   `json:"modify_at"`
	ShortComment string   `json:"short_comment"`
	Year         string   `json:"year"`
}
