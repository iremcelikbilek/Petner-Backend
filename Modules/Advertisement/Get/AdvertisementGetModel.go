package Advertisement

import (
	"time"

	addModel "../Add"
)

type AdvertisementGetListData struct {
	AdvertisementID      string                       `json:"advertisementID"`
	AdvertisementTitle   string                       `json:"advertisementTitle"`
	AdvertisementAnimal  Animal                       `json:"advertisementAnimal"`
	AdvertisementAddress Adress                       `json:"advertisementAddress"`
	AdvertisementType    addModel.AdvertisementType   `json:"advertisementType"`
	Status               addModel.AdvertisementStatus `json:"status"`
	Date                 string                       `json:"date"`
	IsDeleted            bool                         `json:"isDeleted"`
	FullDate             string                       `json:"fullDate"`
}

type Adress struct {
	Province string `json:"province"`
	District string `json:"district"`
}

type Animal struct {
	Genre       string `json:"genre"`
	AnimalPhoto string `json:"animalPhoto"`
}

type AdvSlice []AdvertisementGetListData

func (p AdvSlice) Len() int {
	return len(p)
}

func (p AdvSlice) Less(i, j int) bool {
	t1, _ := time.Parse(time.RFC3339Nano, p[i].FullDate)
	t2, _ := time.Parse(time.RFC3339Nano, p[j].FullDate)
	return t1.After(t2)
}

func (p AdvSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
