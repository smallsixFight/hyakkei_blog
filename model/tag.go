package model

type Tag struct {
	Id          int64  `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	ModifyAt    string `json:"modify_at"`
	//UseCount    int    `json:"use_count"`
}
