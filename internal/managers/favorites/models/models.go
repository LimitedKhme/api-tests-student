package models

// AddFavoriteRequest описывает тело запроса для добавления объявления в избранное.
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}
