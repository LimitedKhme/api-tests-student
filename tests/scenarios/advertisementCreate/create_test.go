package advertisementCreate

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/tests"
)

const photoPath = "../../../photo.png"

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

	var adID string

	s.Run("Создание объявления со всеми параметрами", func() {
		respBody := advManager.CreateAdvertisement(s.T(), s.loginData.Token, params, http.StatusCreated)
		adID = gjson.Get(respBody, "id").String()
		require.NotEmpty(s.T(), adID, "ID созданного объявления не должен быть пустым")
	})

	defer func() {
		if adID != "" {
			advManager.DeleteAdvertisement(s.T(), s.loginData.Token, adID, http.StatusNoContent)
		}
	}()

	s.Run("Объявление доступно через GET /advertisement и возвращает переданные параметры", func() {
		ad := advManager.GetAdvertisement(s.T(), adID, http.StatusOK)
		require.Equal(s.T(), adID, ad.ID, "ID объявления не совпадает")
		require.Equal(s.T(), params.Title, ad.Title, "Title объявления не совпадает")
		require.Equal(s.T(), params.Description, ad.Description, "Description объявления не совпадает")
		require.Equal(s.T(), int(*params.Price), ad.Price, "Price объявления не совпадает")
		require.Equal(s.T(), *params.Quantity, ad.Quantity, "Quantity объявления не совпадает")
	})

	s.Run("Объявление появляется в поиске по полному вхождению названия", func() {
		list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, params.Title)
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, params.Title)
	})

	s.Run("Фотография загружена и доступна по ручке GET /advertisements/{id}/photos", func() {
		photosBody := advManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, adID, http.StatusOK)
		require.NotEmpty(s.T(), photosBody, "Ответ со списком фотографий пуст")
		require.Contains(s.T(), photosBody, "http", "В ответе должна присутствовать ссылка на загруженное фото")
	})
}

func (s *TestSuite) TestCreateAdvertisement_Negative() {
	validPrice := 100.0
	invalidPrice := -50.0
	validQuantity := 1

	cases := []struct {
		name           string
		params         models.CreateAdvertisementFormParams
		expectedStatus int
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
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			respBody := advManager.CreateAdvertisement(s.T(), s.loginData.Token, tc.params, tc.expectedStatus)
			require.NotEmpty(s.T(), gjson.Get(respBody, "error").String(), "Должно вернуться сообщение об ошибке")
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
		respBody := advManager.CreateAdvertisement(s.T(), "invalid_token", validParams, http.StatusUnauthorized)
		require.NotEmpty(s.T(), gjson.Get(respBody, "error").String(), "Должно вернуться сообщение об ошибке авторизации")
	})
}
