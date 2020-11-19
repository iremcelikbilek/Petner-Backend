package Comment

import (
	"net/http"

	fb "../Firebase"
	util "../Utils"
)

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var response util.GeneralResponseModel

	advIDKeys, ok := r.URL.Query()["AdvID"]

	if !ok || len(advIDKeys[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "İlan id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

	commentIDKeys, ok := r.URL.Query()["CommentID"]

	if !ok || len(commentIDKeys[0]) < 1 {
		response = util.GeneralResponseModel{
			true, "Yorum id bilgisi gönderilmelidir", nil,
		}
		w.Write(response.ToJson())
		return
	}

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

	if deleteError := fb.DeleteComment(advIDKeys[0], commentIDKeys[0], userMail); deleteError != nil {
		response = util.GeneralResponseModel{
			true, deleteError.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}
	response = util.GeneralResponseModel{
		false, "Yorum başarıyla silindi", nil,
	}
	w.Write(response.ToJson())
}
