package main

import (
	"fmt"
	"net/http"
	"os"

	AdvertisementAdd "../Modules/Advertisement/Add"
	AdvertisementGet "../Modules/Advertisement/Get"
	Login "../Modules/Auth/Login"
	NewPassword "../Modules/Auth/NewPassword"
	Reset "../Modules/Auth/ResetPassword"
	Signup "../Modules/Auth/SignUp"
	PhotoUpload "../Modules/PhotoUploader"
	Util "../Modules/Utils"

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
	go http.HandleFunc("/upload-photo", PhotoUpload.HandleUpload)
	go http.HandleFunc("/", handleHome)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
