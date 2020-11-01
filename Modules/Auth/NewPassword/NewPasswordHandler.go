package NewPassword

import (
	"encoding/json"
	"fmt"
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
	resetModel "../ResetPassword"
	signUp "../SignUp"
	"github.com/mitchellh/mapstructure"
)

func NewPassword(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	var passwordData NewPasswordModel

	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	if !util.IsEmailValid(passwordData.Mail) {
		response = util.GeneralResponseModel{
			true, "eMail geçersiz", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(passwordData.Password) < 8 {
		response = util.GeneralResponseModel{
			true, "Parola en az 8 karakterli olmalıdır.", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fetchedData := fb.GetFilteredData("/resetPasswordCodes", "personEmail", passwordData.Mail)
	var result resetModel.ResetPasswordDataModel
	mapstructure.Decode(fetchedData, &result)

	if result.Code != passwordData.Code {
		response = util.GeneralResponseModel{
			true, "Kod geçersiz", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userData := fb.GetFilteredData("/persons", "personEmail", passwordData.Mail)
	var userResult signUp.SignUpModel
	mapstructure.Decode(userData, &userResult)
	userResult.Password = util.PasswordHasher(passwordData.Password)
	fb.UpdateFilteredData("/persons", "personEmail", passwordData.Mail, userResult)

	deleteCodeError := fb.Delete("/resetPasswordCodes", "personEmail", passwordData.Mail)
	if deleteCodeError != nil {
		fmt.Println(deleteCodeError.Error())
	}

	response = util.GeneralResponseModel{
		false, "Şifre başarıyla güncellendi", nil,
	}
	w.Write(response.ToJson())
}
