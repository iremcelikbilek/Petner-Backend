package ResetPassword

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	fb "../../Firebase"
	sender "../../MailSender"
	util "../../Utils"
	"github.com/mitchellh/mapstructure"
)

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel
	var resetData ResetPasswordModel

	err := json.NewDecoder(r.Body).Decode(&resetData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, "Gelen veriler hatalı", nil,
		}
		w.Write(response.ToJson())
		return
	}

	if !util.IsEmailValid(resetData.PersonEmail) {
		response = util.GeneralResponseModel{
			true, "eMail geçersiz", nil,
		}
		w.Write(response.ToJson())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response = util.GeneralResponseModel{
		false, "Başarılı", nil,
	}
	w.Write(response.ToJson())

	fb.DeleteAllFilteredDatas("/resetPasswordCodes", "personEmail", resetData.PersonEmail)

	fetchedData := fb.GetFilteredData("/persons", "personEmail", resetData.PersonEmail)
	var result ResetPasswordModel
	mapstructure.Decode(fetchedData, &result)

	if result.PersonEmail == "" {
		return
	}

	code, err := GenerateOTP(6)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var dataModel ResetPasswordDataModel = ResetPasswordDataModel{resetData.PersonEmail, code}

	fbError := fb.PushData("/resetPasswordCodes", dataModel)
	if fbError != nil {
		fmt.Println(fbError.Error())
		return
	}

	go sender.MailSender("Şifre sıfırlama", code, resetData.PersonEmail)
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
