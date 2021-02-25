package model

type Dashboard struct {
	Statistics
	InitTime int64 `json:"init_time"`
}

type Statistics struct {
	ArticleInfo     PostRecord `json:"article_info"`
	PageInfo        PostRecord `json:"page_info"`
	FriendLinkCount int32      `json:"friend_link_count"`
	BookCount       int32      `json:"book_count"`
	VisitorCount    int32      `json:"visitor_count"`
}

type PostRecord struct {
	PublishCount int32  `json:"publish_count"`
	DraftCount   int32  `json:"draft_count"`
	LastAdd      string `json:"last_add"`
	LastPublish  string `json:"last_publish"`
}
