package Advertisement

import addModel "../Add"

type AdvertisementGetListData struct {
	AdvertisementID      string                       `json:"advertisementID"`
	AdvertisementTitle   string                       `json:"advertisementTitle"`
	AdvertisementAnimal  Animal                       `json:"advertisementAnimal"`
	AdvertisementAddress Adress                       `json:"advertisementAddress"`
	AdvertisementType    addModel.AdvertisementType   `json:"advertisementType"`
	Status               addModel.AdvertisementStatus `json:"status"`
	Date                 string                       `json:"date"`
}

type Adress struct {
	Province string `json:"province"`
	District string `json:"district"`
}

type Animal struct {
	Genre       string `json:"genre"`
	AnimalPhoto string `json:"animalPhoto"`
}
