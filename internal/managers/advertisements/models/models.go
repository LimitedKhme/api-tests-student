package models

// Photo описывает фотографию объявления.
type Photo struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       int     `json:"price"`
	Quantity    int     `json:"quantity"`
	Photos      []Photo `json:"photos"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

// CreateAdvertisementFormParams описывает параметры формы для создания объявления.
type CreateAdvertisementFormParams struct {
	Title       string
	Description string
	Price       *float64
	Quantity    *int
	PhotoPaths  []string
}
