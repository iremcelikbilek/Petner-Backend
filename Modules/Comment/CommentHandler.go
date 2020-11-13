package Comment

import (
	"encoding/json"
	"net/http"
	"time"

	fb "../Firebase"
	util "../Utils"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
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

	if !ok || len(keys[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "İlan id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

	var comment CommentModel
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	if len(comment.Comment) < 15 {
		response = util.GeneralResponseModel{
			true, "Yorum minimum 15 karakterli olmalıdır", nil,
		}
		w.Write(response.ToJson())
		return
	}

	userData := fb.GetFilteredData("/persons", "personEmail", userMail)
	itemMap := userData.(map[string]interface{})
	userName := itemMap["personName"].(string)
	userLastName := itemMap["personLastName"].(string)

	nowDate, _ := time.Now().MarshalText()
	var commentDbData CommentDbModel = CommentDbModel{
		PersonEmail:    userMail,
		PersonName:     userName,
		PersonLastName: userLastName,
		Comment:        comment.Comment,
		Date:           string(nowDate),
	}

	if err := fb.CommentAdd(keys[0], commentDbData); err != nil {
		response = util.GeneralResponseModel{
			true, "Yorum ekleme başarısız", nil,
		}
		w.Write(response.ToJson())
		return
	}

	response = util.GeneralResponseModel{
		false, "Yorum ekleme başarılı", nil,
	}
	w.Write(response.ToJson())
}
