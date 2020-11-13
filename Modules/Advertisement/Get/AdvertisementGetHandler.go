package Advertisement

import (
	"fmt"
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
	keys, ok := r.URL.Query()["id"]

	if !ok {
		handleGetList(&w, r)
		return
	} else if len(keys[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "İlan id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

	data := fb.GetFilteredData("/advertisement", "advertisementID", keys[0])
	itemMap := data.(map[string]interface{})
	if data == nil || itemMap["isDeleted"] == true {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		true, "Başarılı", data,
	}
	w.Write(response.ToJson())
}
func handleGetList(w *http.ResponseWriter, r *http.Request) {
	util.EnableCors(&(*w))

	if r.Method == http.MethodOptions {
		(*w).WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	allData := fb.ReadData("/advertisement")
	itemsMap := allData.(map[string]interface{})

	if allData == nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		response = util.GeneralResponseModel{
			true, "Bir hata oluştu", nil,
		}
		(*w).Write(response.ToJson())
		return
	}

	var advertisements []AdvertisementGetListData
	for _, data := range itemsMap {
		var advertisement addModel.AdvertisementDataModel
		mapstructure.Decode(data, &advertisement)
		fmt.Println(advertisement.Deleted)
		if data.(map[string]interface{})["isDeleted"] != true {
			t, _ := time.Parse(time.RFC3339Nano, advertisement.AdvEntryDate)
			advertisement.AdvEntryDate = t.Format("2 January 2006")
			var imageURL string
			if len(advertisement.AdvertisementAnimal.AnimalPhotos) > 0 {
				imageURL = advertisement.AdvertisementAnimal.AnimalPhotos[0]
			}

			var advertisementListData = AdvertisementGetListData{
				AdvertisementID:    advertisement.AdvertisementID,
				AdvertisementTitle: advertisement.AdvertisementTitle,
				AdvertisementAnimal: Animal{
					Genre:       advertisement.AdvertisementAnimal.Genre,
					AnimalPhoto: imageURL,
				},
				AdvertisementAddress: Adress{
					Province: advertisement.AdvertisementAddress.Province,
					District: advertisement.AdvertisementAddress.District,
				},
				AdvertisementType: advertisement.AdvertisementType,
				Status:            advertisement.Status,
				Date:              t.Format("2 January 2006"),
			}
			advertisements = append(advertisements, advertisementListData)
		}
	}

	response = util.GeneralResponseModel{
		false, "Veriler başarıyla getirildi", advertisements,
	}
	(*w).Write(response.ToJson())
}
