package models

type UnreadCount struct {
	Diplomacy  string `json:"diplomacy" dynamodbav:"diplomacy"`
	Events     string `json:"events" dynamodbav:"events"`
	GlobalChat string `json:"chat" dynamodbav:"chat"`
}

func (u *UnreadCount) AllZero() bool {
	return u.Diplomacy == "0" && u.Events == "0" && u.GlobalChat == "0"
}

func CompareUnreadCounts(a, b UnreadCount) bool {
	return a.Diplomacy == b.Diplomacy && a.Events == b.Events && a.GlobalChat == b.GlobalChat
}
