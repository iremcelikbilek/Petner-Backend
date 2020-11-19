package Advertisement

import (
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
	addModel "../Add"
	"github.com/mitchellh/mapstructure"
)

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
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

	currentAdvertisementData.Status = 1

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
