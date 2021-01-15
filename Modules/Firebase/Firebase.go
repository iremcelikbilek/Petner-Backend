package Firebase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	database "firebase.google.com/go/db"
)

var ctx = context.Background()
var client *database.Client

func ConnectFirebase() {

	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASECREDENTIONAL")))

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}

	firebaseClient, err := app.DatabaseWithURL(ctx, "https://petner-d1889.firebaseio.com/")
	if err != nil {
		log.Fatal(err)
	}

	client = firebaseClient
}

func WriteData(path string, data interface{}) error {
	err := client.NewRef(path).Set(ctx, data)
	return err

}

func PushData(path string, data interface{}) error {
	_, err := client.NewRef(path).Push(ctx, data)
	return err
}

func PushFilteredData(path string, child string, equal string, newChild string, pushData interface{}) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}
	PushData(path+"/"+dataParentName+"/"+newChild, pushData)
	if err != nil {
		return err
	}
	return nil
}

func GetFilteredData(path string, child string, equal string) interface{} {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		fmt.Println(err)
	}
	itemsMap := data.(map[string]interface{})

	var responseData interface{}
	for _, v := range itemsMap {
		responseData = v
		break
	}
	return responseData
}

func UpdateFilteredData(path string, child string, equal string, updatedData interface{}) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	newData := map[string]interface{}{
		dataParentName: updatedData,
	}

	err = client.NewRef(path).Update(ctx, newData)
	if err != nil {
		return err
	}
	return nil
}

func CommentAdd(advertisementID string, comment interface{}) error {
	var data interface{}
	err := client.NewRef("/advertisement").OrderByChild("advertisementID").EqualTo(advertisementID).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	if dataParentName == "" {
		return errors.New("İlan bulunamadı")
	}

	return PushData("/advertisement/"+dataParentName+"/comments", comment)
}

func UpdateUserSpesificData(path string, child string, equal string, updatedData interface{}, user string) error {
	var data interface{}
	err := client.NewRef("/persons").OrderByChild("personEmail").EqualTo(user).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	return UpdateFilteredData("/persons/"+dataParentName+path, child, equal, updatedData)
}

func DeleteComment(advID string, commentID string, userMail string) error {
	var data interface{}
	err := client.NewRef("/advertisement").OrderByChild("advertisementID").EqualTo(advID).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})
	if itemsMap == nil {
		return errors.New("İlan bulunamadı")
	}

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	if dataParentName == "" {
		return errors.New("İlan bulunamadı")
	}

	advData := itemsMap[dataParentName]
	if advData == nil {
		return errors.New("İlan bulunamadı")
	}
	advDataMap := advData.(map[string]interface{})
	if advDataMap == nil {
		return errors.New("İlan bulunamadı")
	}

	commentsData := advDataMap["comments"]
	if commentsData == nil {
		return errors.New("Yorum bulunamadı")
	}
	commentsMap := commentsData.(map[string]interface{})

	if commentsMap == nil {
		return errors.New("Yorum bulunamadı")
	}
	for i, v := range commentsMap {
		comment := v.(map[string]interface{})
		if comment == nil {
			return errors.New("Yorum bulunamadı")
		}
		if comment["commentID"] == commentID {
			if comment["personEmail"].(string) == userMail {
				return client.NewRef("advertisement/" + dataParentName + "/comments/" + i).Delete(ctx)
			} else {
				return errors.New("Yorumu silmeye yetkiniz yok")
			}
			break
		}
	}
	return errors.New("Yorumu bulunamadı")
}

func Delete(path string, child string, equal string) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	err = client.NewRef(path + "/" + dataParentName).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ReadData(path string) interface{} {
	var data interface{}
	if err := client.NewRef(path).Get(ctx, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func DeleteAllFilteredDatas(path string, child string, equal string) {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		return
	}
	itemsMap := data.(map[string]interface{})

	for i, _ := range itemsMap {
		client.NewRef(path + "/" + i).Delete(ctx)
	}
}

func DeleteFavoriteAdvertisement(id string, user string) error {
	var data interface{}
	err := client.NewRef("/persons").OrderByChild("personEmail").EqualTo(user).Get(ctx, &data)
	if err != nil {
		return err
	}
	itemsMap := data.(map[string]interface{})

	var dataParentName string
	for i, _ := range itemsMap {
		dataParentName = i
		break
	}

	userData := itemsMap[dataParentName].(map[string]interface{})
	favorites := userData["favorites"].(map[string]interface{})

	for index, value := range favorites {
		if value == id {
			client.NewRef("/persons/" + dataParentName + "/favorites/" + index).Delete(ctx)
		}
	}

	return nil
}
