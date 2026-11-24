package advertisementCreate

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/tests"
)

const photoPath = "photo.png"

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestAdvertisementCreateSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	tests.SetupSuite()
	tests.Precondition("Авторизация пользователя из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TestCreateAdvertisement_Positive() {
	price := 1000.0
	quantity := 5
	params := models.CreateAdvertisementFormParams{
		Title:       "Эксклюзивный велосипед",
		Description: "Почти новый",
		Price:       &price,
		Quantity:    &quantity,
		PhotoPaths:  []string{photoPath},
	}

	var createdAd models.Advertisement

	s.Run("Создание объявления со всеми параметрами", func() {
		respBody := advManager.CreateAdvertisement(s.T(), s.loginData.Token, params, http.StatusCreated)
		err := json.Unmarshal([]byte(respBody), &createdAd)
		require.NoError(s.T(), err, "Не удалось распарсить ответ создания объявления: %s", respBody)
		require.NotEmpty(s.T(), createdAd.ID, "ID созданного объявления не должен быть пустым")
		require.Len(s.T(), createdAd.Photos, len(params.PhotoPaths), "Количество фото в ответе POST не совпадает с отправленным")
		require.NotEmpty(s.T(), createdAd.Photos[0].ID, "ID загруженной фотографии не должен быть пустым")
	})

	defer func() {
		if createdAd.ID != "" {
			advManager.DeleteAdvertisement(s.T(), s.loginData.Token, createdAd.ID, http.StatusNoContent)
		}
	}()

	s.Run("Объявление доступно через GET /advertisement и возвращает переданные параметры", func() {
		ad := advManager.GetAdvertisement(s.T(), createdAd.ID, http.StatusOK)
		require.Equal(s.T(), createdAd.ID, ad.ID, "ID объявления не совпадает")
		require.Equal(s.T(), params.Title, ad.Title, "Title объявления не совпадает")
		require.Equal(s.T(), params.Description, ad.Description, "Description объявления не совпадает")
		require.Equal(s.T(), *params.Price, ad.Price, "Price объявления не совпадает")
		require.Equal(s.T(), *params.Quantity, ad.Quantity, "Quantity объявления не совпадает")
		require.Len(s.T(), ad.Photos, len(createdAd.Photos), "Количество фото в GET не совпадает с POST")
		require.Equal(s.T(), createdAd.Photos[0].ID, ad.Photos[0].ID, "ID фото в GET не совпадает с POST")
	})

	s.Run("Объявление появляется в поиске по полному вхождению названия", func() {
		list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, params.Title)
		var found *models.Advertisement
		for i := range list.Items {
			if list.Items[i].ID == createdAd.ID {
				found = &list.Items[i]
				break
			}
		}
		require.NotNil(s.T(), found, "Созданное объявление с ID %s не найдено в поиске", createdAd.ID)
		require.Equal(s.T(), params.Title, found.Title, "Title найденного объявления не совпадает")
	})

	s.Run("Фотография загружена и доступна по ручке GET /advertisements/{id}/photos", func() {
		photosBody := advManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdAd.ID, http.StatusOK)
		var photos []models.Photo
		err := json.Unmarshal([]byte(photosBody), &photos)
		require.NoError(s.T(), err, "Не удалось распарсить список фотографий: %s", photosBody)
		require.Len(s.T(), photos, len(createdAd.Photos), "Количество фото в /photos не совпадает с созданным")
		require.Equal(s.T(), createdAd.Photos[0].ID, photos[0].ID, "ID фото в /photos не совпадает с ID при создании")
	})
}

func (s *TestSuite) TestCreateAdvertisement_Negative() {
	validPrice := 100.0
	invalidPrice := -50.0
	validQuantity := 1

	cases := []struct {
		name            string
		params          models.CreateAdvertisementFormParams
		expectedStatus  int
		expectedError   string
		expectedMessage string
	}{
		{
			name: "Пустой title",
			params: models.CreateAdvertisementFormParams{
				Title:       "",
				Description: "Описание",
				Price:       &validPrice,
				Quantity:    &validQuantity,
				PhotoPaths:  []string{photoPath},
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "invalid_request", // Уточним после прогона по фактическому ответу сервера
			expectedMessage: "title is required",
		},
		{
			name: "Отрицательная цена",
			params: models.CreateAdvertisementFormParams{
				Title:       "Велосипед",
				Description: "Описание",
				Price:       &invalidPrice,
				Quantity:    &validQuantity,
				PhotoPaths:  []string{photoPath},
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "invalid_request",
			expectedMessage: "price must be between 0 and 1000000 inclusive",
		},
		{
			name: "Слишком длинный title",
			params: models.CreateAdvertisementFormParams{
				Title:       strings.Repeat("А", 300),
				Description: "Описание",
				Price:       &validPrice,
				Quantity:    &validQuantity,
				PhotoPaths:  []string{photoPath},
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   "invalid_request",
			expectedMessage: "title must not exceed 50 characters",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			respBody := advManager.CreateAdvertisement(s.T(), s.loginData.Token, tc.params, tc.expectedStatus)
			assertions.AssertAPIError(s.T(), respBody, tc.expectedError, tc.expectedMessage)
		})
	}

	s.Run("Запрос без авторизации (неверный токен)", func() {
		validParams := models.CreateAdvertisementFormParams{
			Title:       "Велосипед",
			Description: "Описание",
			Price:       &validPrice,
			Quantity:    &validQuantity,
			PhotoPaths:  []string{photoPath},
		}
		respBody := advManager.CreateAdvertisement(s.T(), tests.InvalidToken, validParams, http.StatusUnauthorized)
		assertions.AssertAPIError(s.T(), respBody, "unauthorized", "Invalid or expired token")
	})
}
