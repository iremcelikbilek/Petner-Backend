package Advertisement

import (
	"net/http"

	fb "../../Firebase"
	util "../../Utils"
)

func AdvertisementDeleteHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var userMail string
	var response util.GeneralResponseModel
	if isSucessToken, message := util.CheckToken(r); !isSucessToken {
		response = util.GeneralResponseModel{
			true, message, nil,
		}
		w.Write(response.ToJson())
		return
	} else {
		userMail = message
	}

	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "İlan id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
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

	itemMap := data.(map[string]interface{})
	if itemMap["ownerUser"].(map[string]interface{})["personEmail"] != userMail {
		w.WriteHeader(http.StatusForbidden)
		response = util.GeneralResponseModel{
			true, "İzniniz yok", nil,
		}
		w.Write(response.ToJson())
		return
	}
	itemMap["isDeleted"] = true

	if deleteError := fb.UpdateFilteredData("/advertisement", "advertisementID", keys[0], itemMap); deleteError != nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, deleteError.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	if deleteError := fb.UpdateUserSpesificData("/advertisements", "advertisementID", keys[0], itemMap, userMail); deleteError != nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, deleteError.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Silme Başarılı", nil,
	}
	w.Write(response.ToJson())
}
