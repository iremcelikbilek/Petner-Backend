package main

import (
	"net/http"

	Login "../Modules/Auth/Login"
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
