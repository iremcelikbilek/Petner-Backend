package main

import (
	"net/http"

	Signup "../Modules/Auth/SignUp"

	fb "../Modules/Firebase"
)

func main() {
	go fb.ConnectFirebase()
	createServer()

}

func createServer() {
	go http.HandleFunc("/signup", Signup.HandleSignUp)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		print(err)
	}
}
