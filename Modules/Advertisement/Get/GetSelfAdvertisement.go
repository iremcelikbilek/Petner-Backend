package Advertisement

import (
	"net/http"
	"time"

	fb "../../Firebase"
	util "../../Utils"
	addModel "../Add"
	"github.com/mitchellh/mapstructure"
)

func GetSelfAdvertisementHandler(w http.ResponseWriter, r *http.Request) {
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

	data := fb.GetFilteredData("/persons", "personEmail", userMail)
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	itemsMap := data.(map[string]interface{})
	advertisementData := itemsMap["advertisements"]
	if advertisementData == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}
	advertisementsMap := advertisementData.(map[string]interface{})
	var advertisements []AdvertisementGetListData

	for _, data := range advertisementsMap {
		var advertisement addModel.AdvertisementDataModel
		mapstructure.Decode(data, &advertisement)
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
			IsDeleted:         data.(map[string]interface{})["isDeleted"].(bool),
		}
		advertisements = append(advertisements, advertisementListData)
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", advertisements,
	}
	w.Write(response.ToJson())
}
