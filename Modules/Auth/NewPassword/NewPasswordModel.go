package NewPassword

type NewPasswordModel struct {
	Mail     string `json:"personEmail"`
	Code     string `json:"code"`
	Password string `json:"personPassword"`
}
