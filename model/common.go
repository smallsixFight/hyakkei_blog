package model

type CacheVal string

const (
	FriendLinkCount CacheVal = "friend_link_count"
	BookCount       CacheVal = "book_count"
	VisitorCount    CacheVal = "visitor_count"
)

type Reply struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func (reply *Reply) SetCode(code int32) *Reply {
	reply.Code = code
	return reply
}

func (reply *Reply) SetMessage(message string) *Reply {
	reply.Message = message
	return reply
}

func (reply *Reply) SetSuccess(success bool) *Reply {
	reply.Success = success
	return reply
}

func (reply *Reply) SetData(data interface{}) *Reply {
	reply.Data = data
	return reply
}

func (reply *Reply) SetTotal(total int64) *Reply {
	reply.Total = total
	return reply
}

func (reply *Reply) SetToken(token string) *Reply {
	reply.Token = token
	return reply
}
