package model

type BasePostInfo struct {
	Id       int64    `json:"id,string"`
	Title    string   `json:"title"`
	Slug     string   `json:"slug"`
	Typ      string   `json:"typ"`
	CreateAt string   `json:"create_at"`
	ModifyAt string   `json:"modify_at"`
	Status   PostType `json:"status"`
	Author   string   `json:"author"`
	Tags     []string `json:"tags,omitempty"`
}

type Post struct {
	BasePostInfo
	MarkdownText string `json:"markdown_text"`
	HtmlText     string `json:"html_text"`
}

type PostType string

const (
	Publish PostType = "publish"
	Draft   PostType = "draft"
)
