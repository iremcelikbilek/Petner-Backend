package main

import (
	"net/http"

	AdvertisementAdd "../Modules/Advertisement/Add"
	AdvertisementGet "../Modules/Advertisement/Get"
	Login "../Modules/Auth/Login"
	NewPassword "../Modules/Auth/NewPassword"
	Reset "../Modules/Auth/ResetPassword"
	Signup "../Modules/Auth/SignUp"

	fb "../Modules/Firebase"
)

func main() {
	go fb.ConnectFirebase()
	createServer()

}

func createServer() {
	go http.HandleFunc("/signup", Signup.HandleSignUp)
	go http.HandleFunc("/login", Login.HandleLogin)
	go http.HandleFunc("/resetPassword", Reset.ResetPasswordHandler)
	go http.HandleFunc("/newPassword", NewPassword.NewPassword)

	go http.HandleFunc("/advertisement/add", AdvertisementAdd.AdvertisementAddHandler)
	go http.HandleFunc("/advertisement/get", AdvertisementGet.AdvertisementGetHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
