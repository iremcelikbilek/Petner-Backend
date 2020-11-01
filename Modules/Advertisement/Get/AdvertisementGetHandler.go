package Advertisement

import (
	"net/http"
	"time"

	fb "../../Firebase"
	util "../../Utils"
	addModel "../Add"
	"github.com/mitchellh/mapstructure"
)

func AdvertisementGetHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	var advertisements []addModel.AdvertisementDataModel

	allData := fb.ReadData("/advertisement")

	if allData == nil {
		w.WriteHeader(http.StatusInternalServerError)
		response = util.GeneralResponseModel{
			true, "Bir hata oluştu", nil,
		}
		w.Write(response.ToJson())
		return
	}

	itemsMap := allData.(map[string]interface{})
	for _, data := range itemsMap {
		var advertisement addModel.AdvertisementDataModel
		mapstructure.Decode(data, &advertisement)
		t, _ := time.Parse(time.RFC3339Nano, advertisement.AdvEntryDate)
		advertisement.AdvEntryDate = t.Format("2 January 2006")
		advertisements = append(advertisements, advertisement)
	}

	response = util.GeneralResponseModel{
		false, "Veriler başarıyla getirildi", advertisements,
	}
	w.Write(response.ToJson())
}
