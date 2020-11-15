package Comment

type CommentModel struct {
	Comment string `json:"comment"`
}

type CommentDbModel struct {
	PersonEmail    string `json:"personEmail"`
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	Comment        string `json:"comment"`
	Date           string `json:"date"`
	IsDeletable    bool   `json:"isDeletable"`
}
