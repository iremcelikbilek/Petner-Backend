package ResetPassword

type ResetPasswordModel struct {
	PersonEmail string `json:"personEmail"`
}

type ResetPasswordDataModel struct {
	PersonEmail string `json:"personEmail"`
	Code        string `json:"code"`
}
