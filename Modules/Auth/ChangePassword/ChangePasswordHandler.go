package ChangePassword

import (
	"encoding/json"
	"fmt"
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
)

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	var userMail string

	if isSucessToken, message := util.CheckToken(r); !isSucessToken {
		response = util.GeneralResponseModel{
			true, message, nil,
		}
		w.Write(response.ToJson())
		return
	} else {
		userMail = message
	}

	var changePassword ChangePasswordModel
	if err := json.NewDecoder(r.Body).Decode(&changePassword); err != nil {
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	fmt.Println(changePassword.NewPassword)
	fmt.Println(changePassword.OldPassword)

	if len(changePassword.NewPassword) < 7 {
		response = util.GeneralResponseModel{
			true, "Şifre minimum 7 karakterli olmalıdır", nil,
		}
		w.Write(response.ToJson())
		return
	}

	userData := fb.GetFilteredData("/persons", "personEmail", userMail)
	itemMap := userData.(map[string]interface{})

	if !util.ComparePasswords(itemMap["password"].(string), changePassword.OldPassword) {
		response = util.GeneralResponseModel{
			true, "Şifreniz hatalı", nil,
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.ToJson())
		return
	}

	itemMap["password"] = util.PasswordHasher(changePassword.NewPassword)

	if err := fb.UpdateFilteredData("/persons", "personEmail", userMail, itemMap); err != nil {
		response = util.GeneralResponseModel{
			true, "Şifreniz güncellenemedi", nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Şifre güncelleme başarılı", nil,
	}

	w.Write(response.ToJson())
}
