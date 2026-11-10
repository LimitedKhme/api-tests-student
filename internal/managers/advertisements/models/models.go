package models

type Advertisement struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}