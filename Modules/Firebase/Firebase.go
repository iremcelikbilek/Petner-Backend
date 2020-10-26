package Firebase

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	database "firebase.google.com/go/db"
)

var ctx = context.Background()
var client *database.Client

func ConnectFirebase() {

	opt := option.WithCredentialsJSON(credentional)

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
		fmt.Println(err)
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
		fmt.Println(err)
		return err
	}
	return nil
}

func Delete(path string, child string, equal string) error {
	var data interface{}
	err := client.NewRef(path).OrderByChild(child).EqualTo(equal).Get(ctx, &data)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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

var credentional = []byte(`{
	"type": "service_account",
	"project_id": "petner-d1889",
	"private_key_id": "f2c4ae1a083d6ef7482c8ed060bb329802275f6d",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC2qAho/4lS6HA/\nMgnvjvUw78W15NhMk2ukZjEKYK3Ke/jL1qktgLjo9kPwGxPOtNKckghXFfcB1gEI\nWpljC2h5El9nuIVeJTU1BQFKYrKqv50Zwnpwbvwys8W4a3+XP3t1Ia3UvFdU0xud\n0++UCtBCOg0K6QCLJn2zxx1uqe15M0tsqRvrAnX0XYanEZeGnGtjextQpeKNCIfb\nLe0bPaRIFqiH/boUq3pLv89EeazJ8H+PAWlH+cWRdKxjPNS9/Scy3wZbqBjHlwfT\nJhds7j5Jixz5SPoHFUXz/3R+mdBHMDTgT+ftMGhxHE94knOOrFgU6iN5evzgxZ0G\nukcElX1NAgMBAAECggEAJhOHs0mnuQyCz62w9AsfUkz9cFJNA8OMgigaa4ElYWsv\n/WAZgsaNZWTm2y/t1F6N8/0eN3c4911C+FiYlpTLeceqcz78MFi9y3hoYTcLazxH\n9dV8fCEqujAPFMd2ANPHOu7jI5CCfQiH/oHudLQ/XzrmOqLBTgCfXiFxIX9TZXsB\nBWydGpsm1qApf0z0XWqm90Vx2QFwPD4Yyhn8DiQ7fG7mYlJ3Nc5Z6WiHQsgsTZlI\nYqrmrIZhF0WREBY7eMqGIpEoKDpHHisR4coGGf7oDibW5B9XhW55LqI6fc67QKnM\ncz4okt+EntMtAmDnccXWXOIwVRQOG4qAYO9ocGIwKQKBgQDwDVr9kuvZYYxK9tvz\nvYfRtpHuTseaQ6DmqN1VW94ZZS3nj+7Jzwwi/AUTRErDVDs3sjn4NCuYNtD/l19V\nrVpLoC4HscnWS8M6kdJ80smttdTKrwRlzXIe/u7z7sJw5KToMtlIvl0mnAfgEOTb\n3EbeuOCXKlqrj3ni75m8aC2t9QKBgQDCyodZk6+UTMquiEPGXD4YF1QCxsVvc44u\n9WQXhJqi8rebM8CjrGgGmGvzOMd/XyU6q+YuL+qvf8aoeDYUzbFxOdZqe2hN1xgZ\nCT6Bp/xXetfzk8xBuVg5b6JjVq4fmx1TfP1Zmg3lAntyvn2CIozSf7WrtbBFtREO\n7z3cMCHi+QKBgH2OHrdebyll0iEreOPFkBJqMW25msDe+ntqe0m4ITSbLSVerQC0\n4J4zvtvS7l+34LlC6PsfHmYg1bO5ks2XPBEuGKVBolYJjnVF7BgJkB7hagkQ/XXZ\nvQTlRkoj6WNu06n3XpqjpskY9y2E6I7uacr4W8/1ATOWeaPuujRHMQ05AoGBALYm\ndVMKi5F+DboPqnD/KQGWLvU5sr55rGe1CJgFZCUkGxWC240yV0Rzm96hJcyxyDqJ\nLIHcRPU/4yD+6HOjtV5P23VPWUYQ8XPX9R+BWrLjKLWZa9O54gozngKOt9zOTCoa\nIz96k6unGpE+GFdsv4rH6bZb/C3zF7SDe7E/QTDhAoGAdBKxoJcNu2YL+q8PxemV\nNsHEF0S3MSx98Cqrok/FAyjceoCxujjneNj+gTcqYgMtJPdaWdD72yyFROXQ2rWj\nKvqvmcLmFfZUuWGT0Z0qij4PtqHDwMuDHpWgd5rGB3u7yRlbJWK4dPC7aJSBhYWr\ntdMS0olHYmKePEqrsszM2Y8=\n-----END PRIVATE KEY-----\n",
	"client_email": "firebase-adminsdk-q8ehl@petner-d1889.iam.gserviceaccount.com",
	"client_id": "107065141860982685176",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-q8ehl%40petner-d1889.iam.gserviceaccount.com"
  }`)
