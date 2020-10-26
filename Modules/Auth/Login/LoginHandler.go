package Login

import (
	"encoding/json"
	"net/http"
	"time"

	fb "../../Firebase"
	util "../../Utils"
	signUp "../SignUp"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

var JWT_Token = []byte("PETNER_JWT_TOKEN")

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	var loginData LoginModel

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	if !util.IsEmailValid(loginData.PersonEmail) {
		response = util.GeneralResponseModel{
			true, "eMail geçersiz", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fetchedData := fb.GetFilteredData("/persons", "personEmail", loginData.PersonEmail)
	var result signUp.SignUpModel
	mapstructure.Decode(fetchedData, &result)

	if result.PersonEmail == "" {
		response = util.GeneralResponseModel{
			true, "Kullanıcı veya şifre hatalı", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	}

	if !util.ComparePasswords(result.Password, loginData.Password) {
		response = util.GeneralResponseModel{
			true, "Kullanıcı veya şifre hatalı", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	}

	expirationTime := time.Now().Add(6 * time.Hour)
	claims := &Claims{
		Username: loginData.PersonEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_Token)
	if err != nil {
		response = util.GeneralResponseModel{
			true, "Bir Sorun Oluştu", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response = util.GeneralResponseModel{
		false, "Giriş Başarılı", nil,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.Write(response.ToJson())
}
