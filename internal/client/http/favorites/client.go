package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite отправляет запрос на добавление объявления в избранное.
func HttpPostFavorite(t *testing.T, token string, req models.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(req).
		Expect(t).End().Response
}

// HttpGetFavorites получает список избранных объявлений текущего пользователя.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по ad_id.
func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + adID).
		Expect(t).End().Response
}
