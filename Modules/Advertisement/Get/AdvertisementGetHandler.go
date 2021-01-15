package Advertisement

import (
	"net/http"
	"sort"
	"time"

	comment "../../Comment"
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
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}
	itemMap := data.(map[string]interface{})
	if itemMap["isDeleted"] == true {
		w.WriteHeader(http.StatusNotFound)
		response = util.GeneralResponseModel{
			true, "İlan bulunamadı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	userData := fb.GetFilteredData("/persons", "personEmail", userMail)
	userDataItemsMap := userData.(map[string]interface{})
	userFavorites := userDataItemsMap["favorites"]

	isFavorited := false
	if userFavorites != nil {
		favoritesMap := userFavorites.(map[string]interface{})
		for _, value := range favoritesMap {
			if keys[0] == value {
				isFavorited = true
				break
			}
		}
	}

	t, _ := time.Parse(time.RFC3339Nano, itemMap["advEntryDate"].(string))
	itemMap["advEntryDate"] = t.Format("2 January 2006")
	commentsObjects := itemMap["comments"]
	if commentsObjects != nil {
		var commentArray comment.CommentSlice = comment.CommentSlice{}
		commentsMap := commentsObjects.(map[string]interface{})
		for _, data := range commentsMap {
			var comment comment.CommentDbModel
			mapstructure.Decode(data, &comment)
			t, _ := time.Parse(time.RFC3339Nano, comment.Date)
			comment.FullDate = comment.Date
			comment.Date = t.Format("2 January 2006")
			comment.IsDeletable = comment.PersonEmail == userMail
			comment.FavoriteCount = len(comment.Favorites)

			for _, mail := range comment.Favorites {
				if mail == userMail {
					comment.IsFavorited = true
					break
				}
			}
			comment.Favorites = nil
			commentArray = append(commentArray, comment)
		}

		sort.Sort(commentArray)
		itemMap["comments"] = commentArray
	}

	itemMap["isFavorited"] = isFavorited

	response = util.GeneralResponseModel{
		false, "Başarılı", itemMap,
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

	var advertisements AdvSlice
	for _, data := range itemsMap {
		var advertisement addModel.AdvertisementDataModel
		mapstructure.Decode(data, &advertisement)

		if data.(map[string]interface{})["isDeleted"] != true {
			t, _ := time.Parse(time.RFC3339Nano, advertisement.AdvEntryDate)

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
				IsDeleted:         advertisement.Deleted,
				FullDate:          advertisement.AdvEntryDate,
			}
			advertisements = append(advertisements, advertisementListData)
		}
	}

	sort.Sort(advertisements)

	response = util.GeneralResponseModel{
		false, "Veriler başarıyla getirildi", advertisements,
	}
	(*w).Write(response.ToJson())
}
