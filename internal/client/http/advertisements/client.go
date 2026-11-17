package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет поиск объявлений по тексту в названии
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement отправляет multipart/form-data запрос на создание нового объявления.
func HttpPostAdvertisement(t *testing.T, token string, params models.CreateAdvertisementFormParams) *http.Response {
	req := apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", params.Title).
		MultipartFormData("description", params.Description)

	if params.Price != nil {
		req = req.MultipartFormData("price", strconv.FormatFloat(*params.Price, 'f', -1, 64))
	}

	if params.Quantity != nil {
		req = req.MultipartFormData("quantity", strconv.Itoa(*params.Quantity))
	}

	if len(params.PhotoPaths) > 0 {
		req = req.MultipartFile("photos", params.PhotoPaths...)
	}

	return req.Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по ID через query-параметр.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID через query-параметр.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).End().Response
}

// HttpGetAdvertisementPhotos получает список фотографий объявления.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).End().Response
}
