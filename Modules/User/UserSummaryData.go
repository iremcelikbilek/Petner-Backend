package User

import (
	"net/http"

	"github.com/mitchellh/mapstructure"

	fb "../Firebase"
	util "../Utils"
)

func UserSummaryHandler(w http.ResponseWriter, r *http.Request) {
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
	var userData UserSummaryData
	mapstructure.Decode(fetchedData, &userData)

	response = util.GeneralResponseModel{
		false, "Bilgiler başarılı şekilde getirildi", userData,
	}
	w.Write(response.ToJson())
}

type UserSummaryData struct {
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonPhone    string `json:"personPhone"`
	PersonEmail    string `json:"personEmail"`
}
