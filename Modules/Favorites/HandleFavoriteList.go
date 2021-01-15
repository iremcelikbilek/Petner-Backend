package Favorites

import (
	"net/http"

	addModel "../Advertisement/Add"
	fb "../Firebase"
	util "../Utils"
	"github.com/mitchellh/mapstructure"
)

func HandleFavoriteList(w http.ResponseWriter, r *http.Request) {
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

	var advertisements []addModel.AdvertisementDataModel

	userData := fb.GetFilteredData("/persons", "personEmail", userMail)
	userDataItemsMap := userData.(map[string]interface{})
	userFavorites := userDataItemsMap["favorites"]

	if userFavorites != nil {
		favoritesMap := userFavorites.(map[string]interface{})
		for _, value := range favoritesMap {
			advertisements = append(advertisements, getAdvertisement(value.(string)))
		}
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", advertisements,
	}
	w.Write(response.ToJson())
}

func getAdvertisement(id string) addModel.AdvertisementDataModel {
	data := fb.GetFilteredData("/advertisement", "advertisementID", id)
	var model addModel.AdvertisementDataModel
	mapstructure.Decode(data, &model)
	return model
}
