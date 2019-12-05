package types

type Message struct {
	Anonymous   Anonymous `json:"anonymous"`
	GroupId     int64     `json:"group_id"`
	Font        int64     `json:"font"`
	Message     string    `json:"message"`
	MessageId   int64     `json:"message_id"`
	MessageType string    `json:"message_type"`
	PostType    string    `json:"post_type"`
	RawMessage  string    `json:"raw_message"`
	SelfId      int64     `json:"self_id"`
	Sender      Sender    `json:"sender"`
	SubType     string    `json:"sub_type"`
	Time        uint64    `json:"time"`
	UserId      int64     `json:"user_id"`
}

type Sender struct {
	Age      int64  `json:"age"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	UserId   int64  `json:"user_id"`
	Role     string `json:"role"`
}

type Anonymous struct {
	Flag string `json:"flag"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
