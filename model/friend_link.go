package model

type FriendLink struct {
	Id          int64  `json:"id,string"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	AvatarUrl   string `json:"avatar_url"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	ModifyAt    string `json:"modify_at"`
}
