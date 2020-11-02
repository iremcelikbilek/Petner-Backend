package PhotoUploader

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	util "../Utils"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	util.EnableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response util.GeneralResponseModel

	r.ParseMultipartForm(10 << 20)
	r.ParseForm()
	path := r.FormValue("path")
	path = strings.Replace(path, " ", "", -1)
	os.MkdirAll("/home/petner/upload-images/"+path, os.ModePerm)

	file, _, error := r.FormFile("file")
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, error.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}
	defer file.Close()

	tempFile, error := ioutil.TempFile("/home/petner/upload-images/"+path, "upload-*.png")
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, error.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}
	defer tempFile.Close()
	os.Chmod(tempFile.Name(), 0644)

	fileBytes, error := ioutil.ReadAll(file)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		response = util.GeneralResponseModel{
			true, error.Error(), nil,
		}
		w.Write(response.ToJson())
		return
	}

	tempFile.Write(fileBytes)
	response = util.GeneralResponseModel{
		false, "Fotoğraf başarıyla yüklendi", "https://petner-cdn.yusufozgul.com" + tempFile.Name()[26:],
	}
	w.Write(response.ToJson())
}
