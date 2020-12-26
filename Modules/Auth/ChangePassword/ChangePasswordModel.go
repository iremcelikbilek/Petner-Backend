package ChangePassword

type ChangePasswordModel struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
