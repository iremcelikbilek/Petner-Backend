package User

import (
	"encoding/json"
	"net/http"

	signUp "../Auth/SignUp"
	fb "../Firebase"
	util "../Utils"
	"github.com/mitchellh/mapstructure"
)

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	fetchedData := fb.GetFilteredData("/persons", "personEmail", userMail)
	var userDbData signUp.SignUpDbModel
	mapstructure.Decode(fetchedData, &userDbData)

	var updatedData UserUpdateModel

	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if updatedData.PersonName != "" {
		if len(updatedData.PersonName) < 2 {
			response = util.GeneralResponseModel{
				true, "İsim veya Soyisim en az 2 karakterli olmalıdır.", nil,
			}
			w.Write(response.ToJson())
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			userDbData.PersonName = updatedData.PersonName
		}
	}

	if updatedData.PersonLastName != "" {
		if len(updatedData.PersonLastName) < 2 {
			response = util.GeneralResponseModel{
				true, "İsim veya Soyisim en az 2 karakterli olmalıdır.", nil,
			}
			w.Write(response.ToJson())
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			userDbData.PersonLastName = updatedData.PersonLastName
		}
	}

	if updatedData.PersonPhone != "" {
		if len(updatedData.PersonPhone) < 11 {
			response = util.GeneralResponseModel{
				true, "Telefon numarası en az 11 karakterli olmalıdır.", nil,
			}
			w.Write(response.ToJson())
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			userDbData.PersonPhone = updatedData.PersonPhone
		}
	}

	if err := fb.UpdateFilteredData("/persons", "personEmail", userMail, userDbData); err != nil {
		response = util.GeneralResponseModel{
			true, err.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Kullanıcı başarılı şekilde güncellendi", nil,
	}
	w.Write(response.ToJson())
}
