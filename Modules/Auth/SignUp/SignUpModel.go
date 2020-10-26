package SignUp

import "github.com/dgrijalva/jwt-go"

type SignUpModel struct {
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonEmail    string `json:"personEmail"`
	PersonPhone    string `json:"personPhone"`
	Password       string `json:"password"`
}

type Claims struct {
	Username string `json:"personEmail"`
	jwt.StandardClaims
}
