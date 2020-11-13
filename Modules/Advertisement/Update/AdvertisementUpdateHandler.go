package Advertisement

import (
	"encoding/json"
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
	addModel "../Add"
	"github.com/mitchellh/mapstructure"
)

func AdvertisementUpdateHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var userMail string
	if isSucessToken, message := util.CheckToken(r); !isSucessToken {
		writeError(message, w)
		return
	} else {
		userMail = message
	}

	var response util.GeneralResponseModel
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		writeError("İlan id bilgisi gönderilmelidir", w)
		return
	}

	data := fb.GetFilteredData("/advertisement", "advertisementID", keys[0])
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}
	var currentAdvertisementData addModel.AdvertisementDataModel
	mapstructure.Decode(data, &currentAdvertisementData)
	if currentAdvertisementData.OwnerUser.PersonEmail != userMail {
		w.WriteHeader(http.StatusForbidden)
		response = util.GeneralResponseModel{
			true, "İzniniz yok", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var advertisementData addModel.AdvertisementAddData
	if err := json.NewDecoder(r.Body).Decode(&advertisementData); err != nil {
		writeError("Gelen veriler hatalı", w)
		return
	}

	if advertisementData.AdvertisementTitle != "" {
		if !controlData(advertisementData.AdvertisementTitle, 5) {
			writeError("Başlık minimum 5 karakterli olmalıdır", w)
			return
		} else {
			currentAdvertisementData.AdvertisementTitle = advertisementData.AdvertisementTitle
		}
	}

	if advertisementData.AdvertisementExplanation != "" {
		if !controlData(advertisementData.AdvertisementExplanation, 15) {
			writeError("Açıklama minimum 15 karakterli olmalıdır", w)
			return
		} else {
			currentAdvertisementData.AdvertisementExplanation = advertisementData.AdvertisementExplanation
		}
	}

	if advertisementData.AdvertisementAnimal.Genre != "" {
		if !controlData(advertisementData.AdvertisementAnimal.Genre, 3) {
			writeError("Hayvan türü minimum 3 karakterli olmalıdır", w)
			return
		} else {
			currentAdvertisementData.AdvertisementAnimal.Genre = advertisementData.AdvertisementAnimal.Genre
		}
	}

	if advertisementData.AdvertisementAnimal.Gender != "" {
		if !controlData(advertisementData.AdvertisementAnimal.Gender, 0) {
			writeError("Hayvan cinsiyeti minimum 1 karakterli olmalıdır", w)
			return
		} else {
			currentAdvertisementData.AdvertisementAnimal.Gender = advertisementData.AdvertisementAnimal.Gender
		}
	}

	if updateError := fb.UpdateFilteredData("/advertisement", "advertisementID", keys[0], currentAdvertisementData); updateError != nil {
		writeError(updateError.Error(), w)
		return
	}

	if updateError := fb.UpdateUserSpesificData("/advertisements", "advertisementID", keys[0], currentAdvertisementData, userMail); updateError != nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, updateError.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Güncelleme Başarılı", nil,
	}
	w.Write(response.ToJson())
}

func controlData(data string, minimum int) bool {
	return len(data) > minimum
}

func writeError(description string, w http.ResponseWriter) {
	var response util.GeneralResponseModel
	w.WriteHeader(http.StatusBadRequest)
	response = util.GeneralResponseModel{
		true, description, nil,
	}
	w.Write(response.ToJson())
}
