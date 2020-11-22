package Comment

import (
	"time"
)

type CommentModel struct {
	Comment string `json:"comment"`
}

type CommentSlice []CommentDbModel

type CommentDbModel struct {
	CommentID      string `json:"commentID"`
	PersonEmail    string `json:"personEmail"`
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	Comment        string `json:"comment"`
	Date           string `json:"date"`
	IsDeletable    bool   `json:"isDeletable"`
	FullDate       string `json:"fullDate"`
}

func (p CommentSlice) Len() int {
	return len(p)
}

func (p CommentSlice) Less(i, j int) bool {
	t1, _ := time.Parse(time.RFC3339Nano, p[i].FullDate)
	t2, _ := time.Parse(time.RFC3339Nano, p[j].FullDate)
	return t1.After(t2)
}

func (p CommentSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
