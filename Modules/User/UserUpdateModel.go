package User

type UserUpdateModel struct {
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonPhone    string `json:"personPhone"`
}
