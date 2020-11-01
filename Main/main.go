package main

import (
	"net/http"

<<<<<<< Updated upstream
=======
	AdvertisementAdd "../Modules/Advertisement/Add"
	AdvertisementGet "../Modules/Advertisement/Get"
>>>>>>> Stashed changes
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
<<<<<<< Updated upstream

=======
	go http.HandleFunc("/advertisement/add", AdvertisementAdd.AdvertisementAddHandler)
	go http.HandleFunc("/advertisement/get", AdvertisementGet.AdvertisementGetHandler)
>>>>>>> Stashed changes
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
