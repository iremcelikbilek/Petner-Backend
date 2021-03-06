package Advertisement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	fb "../../Firebase"
	util "../../Utils"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func AdvertisementAddHandler(w http.ResponseWriter, r *http.Request) {
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
	var advertisementData AdvertisementAddData

	if err := json.NewDecoder(r.Body).Decode(&advertisementData); err != nil {
		writeError("Gelen veriler hatalı", w)
		return
	}

	if !controlData(advertisementData.AdvertisementTitle, 5) {
		writeError("Başlık minimum 5 karakterli olmalıdır", w)
		return
	}

	if !controlData(advertisementData.AdvertisementExplanation, 15) {
		writeError("Açıklama minimum 15 karakterli olmalıdır", w)
		return
	}

	if !(advertisementData.AdvertisementType == 0 || advertisementData.AdvertisementType == 1 || advertisementData.AdvertisementType == 2) {
		writeError("İlan tipi geçersiz", w)
		return
	}

	if !controlData(advertisementData.AdvertisementAnimal.Genre, 3) {
		writeError("Hayvan türü minimum 3 karakterli olmalıdır", w)
		return
	}

	if !controlData(advertisementData.AdvertisementAnimal.Gender, 0) {
		writeError("Hayvan cinsiyeti minimum 1 karakterli olmalıdır", w)
		return
	}

	if advertisementData.AdvertisementAnimal.Age < 0 {
		writeError("Hayvan yaşı 0 ya da daha büyük olmalıdır", w)
		return
	}

	if !controlData(advertisementData.AdvertisementAddress.Province, 3) {
		writeError("İl minimum 3 karakterli olmalıdır", w)
		return
	}

	if !controlData(advertisementData.AdvertisementAddress.District, 3) {
		writeError("İlçe minimum 3 karakterli olmalıdır", w)
		return
	}

	if !controlData(advertisementData.AdvertisementAddress.FullAddress, 10) {
		writeError("Adres detayı minimum 10 karakterli olmalıdır", w)
		return
	}

	if advertisementData.AdvertisementAddress.Latitude == 0.0 || advertisementData.AdvertisementAddress.Longitude == 0.0 {
		writeError("Konum bilgisi hatalı", w)
		return
	}

	fetchedData := fb.GetFilteredData("/persons", "personEmail", userMail)
	var userDbData AdvertisementOwnerData
	mapstructure.Decode(fetchedData, &userDbData)

	nowDate, _ := time.Now().MarshalText()

	dbData := AdvertisementDataModel{
		AdvertisementID:          uuid.New().String(),
		AdvEntryDate:             string(nowDate),
		OwnerUser:                userDbData,
		AdvertisementTitle:       advertisementData.AdvertisementTitle,
		AdvertisementExplanation: advertisementData.AdvertisementExplanation,
		AdvertisementAnimal:      advertisementData.AdvertisementAnimal,
		AdvertisementAddress:     advertisementData.AdvertisementAddress,
		AdvertisementType:        advertisementData.AdvertisementType,
		AdvertisementComments:    []string{},
		FavoriteCount:            0,
		Status:                   0,
		Deleted:                  false,
	}

	var animalPhotos []string
	for _, photo := range dbData.AdvertisementAnimal.AnimalPhotos {
		if photo != "" {
			animalPhotos = append(animalPhotos, photo)
		}
	}
	dbData.AdvertisementAnimal.AnimalPhotos = animalPhotos

	if saveErr := fb.PushData("/advertisement", dbData); saveErr != nil {
		writeError(saveErr.Error(), w)
		return
	}

	response = util.GeneralResponseModel{
		false, "İlan başarıyla oluşturuldu", nil,
	}
	w.Write(response.ToJson())

	if err := fb.PushFilteredData("/persons", "personEmail", userMail, "advertisements", dbData); err != nil {
		fmt.Println(err.Error())
	}
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
