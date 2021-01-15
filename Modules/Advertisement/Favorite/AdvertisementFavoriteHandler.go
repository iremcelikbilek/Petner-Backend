package Favorite

import (
	"fmt"
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
)

func AdvertisementFavoriteHandler(w http.ResponseWriter, r *http.Request) {
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

	advID, ok := r.URL.Query()["id"]

	if !ok || len(advID[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "İlan id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

	state, ok := r.URL.Query()["value"]

	if !ok || len(state[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "Beğeni bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

	data := fb.GetFilteredData("/advertisement", "advertisementID", advID[0])
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	favoriteCount := data.(map[string]interface{})["favoriteCount"].(float64)

	if state[0] == "true" {
		data.(map[string]interface{})["favoriteCount"] = favoriteCount + 1
		if err := fb.PushFilteredData("/persons", "personEmail", userMail, "favorites", data.(map[string]interface{})["advertisementID"]); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		data.(map[string]interface{})["favoriteCount"] = favoriteCount - 1
		if err := fb.DeleteFavoriteAdvertisement(advID[0], userMail); err != nil {
			fmt.Println(err.Error())
		}
	}

	if updateError := fb.UpdateFilteredData("/advertisement", "advertisementID", advID[0], data); updateError != nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, updateError.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())
}
