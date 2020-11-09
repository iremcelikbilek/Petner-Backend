package Advertisement

// User request type
type AdvertisementAddData struct {
	AdvertisementTitle       string            `json:"advertisementTitle"`
	AdvertisementExplanation string            `json:"advertisementExplanation"`
	AdvertisementAnimal      Animal            `json:"advertisementAnimal"`
	AdvertisementAddress     Adress            `json:"advertisementAddress"`
	AdvertisementType        AdvertisementType `json:"advertisementType"`
}

type Adress struct {
	Province    string  `json:"province"`
	District    string  `json:"district"`
	FullAddress string  `json:"fullAddress"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
}

type AdvertisementType int

const (
	FoodHelp AdvertisementType = iota
	Ownership
	Vaccination
)

type Animal struct {
	Genre        string   `json:"genre"`
	Age          int32    `json:"age"`
	Gender       string   `json:"gender"`
	AnimalPhotos []string `json:"animalPhotos"`
}

// Database type
type AdvertisementDataModel struct {
	AdvertisementID          string                 `json:"advertisementID"`
	AdvEntryDate             string                 `json:"advEntryDate"`
	OwnerUser                AdvertisementOwnerData `json:"ownerUser"`
	AdvertisementTitle       string                 `json:"advertisementTitle"`
	AdvertisementExplanation string                 `json:"advertisementExplanation"`
	AdvertisementAnimal      Animal                 `json:"advertisementAnimal"`
	AdvertisementAddress     Adress                 `json:"advertisementAddress"`
	AdvertisementType        AdvertisementType      `json:"advertisementType"`
	AdvertisementComments    []string               `json:"advertisementComments"`
	FavoriteCount            int                    `json:"favoriteCount"`
	Status                   AdvertisementStatus    `json:"status"`
}

type AdvertisementOwnerData struct {
	PersonName     string `json:"personName"`
	PersonLastName string `json:"personLastName"`
	PersonEmail    string `json:"personEmail"`
	PersonPhone    string `json:"personPhone"`
}

type AdvertisementStatus int

const (
	Waiting AdvertisementStatus = iota
	Completed
	NeedsMore
)
