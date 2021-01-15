package Comment

import (
	"net/http"

	fb "../Firebase"
	util "../Utils"
	"github.com/mitchellh/mapstructure"
)

func CommentFavoriteHandler(w http.ResponseWriter, r *http.Request) {
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

	commentID, ok := r.URL.Query()["commentId"]

	if !ok || len(commentID[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "Yorum id bilgisi gönderilmelidir", nil,
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

	var commentMap = data.(map[string]interface{})["comments"].(map[string]interface{})

	if state[0] == "true" { // Beğenildi
		for index, commentData := range commentMap {
			var commentDB CommentDbModel
			mapstructure.Decode(commentData, &commentDB)

			if commentDB.CommentID == commentID[0] {
				commentDB.Favorites = append(commentDB.Favorites, userMail)
			}

			commentMap[index] = commentDB
		}
	} else {
		for index, commentData := range commentMap {
			var commentDB CommentDbModel
			mapstructure.Decode(commentData, &commentDB)

			var newFavorites []string

			for _, mail := range commentDB.Favorites {
				if mail != userMail {
					newFavorites = append(newFavorites, mail)
				}
			}
			commentDB.Favorites = newFavorites
			commentMap[index] = commentDB
		}
	}

	data.(map[string]interface{})["comments"] = commentMap

	if updateError := fb.UpdateFilteredData("/advertisement", "advertisementID", advID[0], data); updateError != nil {
		response = util.GeneralResponseModel{
			true, "Bir hata oluştu", nil,
		}
		w.Write(response.ToJson())
		return
	}
	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())
}
