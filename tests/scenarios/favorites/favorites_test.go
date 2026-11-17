package favorites

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/tests"
)

const (
	photoPath       = "../../../photo.png"
	nonExistentAdID = "00000000-0000-0000-0000-000000000000"
	secondUserEmail = "autotest_helper_user_99@gmail.com"
	secondUserPass  = "1234567A_a"
)

type TestSuite struct {
	suite.Suite
	loginData       authModels.LoginOkResponse
	secondLoginData authModels.LoginOkResponse
}

func TestFavoritesSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	tests.SetupSuite()
	tests.Precondition("Авторизация основного и вспомогательного пользователей")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.secondLoginData = auth.RegisterAndLogin(s.T(), secondUserEmail, secondUserPass)
}

// createAdByToken создает временное объявление от имени владельца токена и возвращает его ID.
func (s *TestSuite) createAdByToken(token string) string {
	price := 500.0
	quantity := 1
	params := advModels.CreateAdvertisementFormParams{
		Title:       "Товар для избранного",
		Description: "Тестовое описание",
		Price:       &price,
		Quantity:    &quantity,
		PhotoPaths:  []string{photoPath},
	}
	respBody := advManager.CreateAdvertisement(s.T(), token, params, http.StatusCreated)
	adID := gjson.Get(respBody, "id").String()
	require.NotEmpty(s.T(), adID, "Не удалось получить ID созданного объявления")
	return adID
}

func (s *TestSuite) TestFavorites_Positive() {
	// Создаем объявление от имени ВТОРОГО пользователя, чтобы основной мог добавить его в избранное
	adID := s.createAdByToken(s.secondLoginData.Token)
	isAdDeleted := false

	defer func() {
		if !isAdDeleted && adID != "" {
			advManager.DeleteAdvertisement(s.T(), s.secondLoginData.Token, adID, http.StatusNoContent)
		}
	}()

	s.Run("Добавление чужого объявления в избранное (201) и проверка в списке GET /favorites", func() {
		favManager.AddToFavorites(s.T(), s.loginData.Token, adID, http.StatusCreated)

		favBody := favManager.GetFavorites(s.T(), s.loginData.Token, http.StatusOK)
		require.True(s.T(), strings.Contains(favBody, adID), "Добавленное объявление должно присутствовать в списке избранного")
	})

	s.Run("Удаление объявления из избранного (204) и проверка, что оно пропало из списка", func() {
		favManager.RemoveFromFavorites(s.T(), s.loginData.Token, adID, http.StatusNoContent)

		favBody := favManager.GetFavorites(s.T(), s.loginData.Token, http.StatusOK)
		require.False(s.T(), strings.Contains(favBody, adID), "Удаленное объявление не должно быть в списке избранного")
	})

	s.Run("При удалении самого объявления через DELETE /advertisement оно исчезает из избранного", func() {
		favManager.AddToFavorites(s.T(), s.loginData.Token, adID, http.StatusCreated)

		// Второй пользователь (владелец) удаляет само объявление
		advManager.DeleteAdvertisement(s.T(), s.secondLoginData.Token, adID, http.StatusNoContent)
		isAdDeleted = true

		favBody := favManager.GetFavorites(s.T(), s.loginData.Token, http.StatusOK)
		require.False(s.T(), strings.Contains(favBody, adID), "После удаления самого объявления оно должно пропасть из избранного")
	})
}

func (s *TestSuite) TestFavorites_Negative() {
	s.Run("Добавление в избранное несуществующего ad_id (400)", func() {
		respBody := favManager.AddToFavorites(s.T(), s.loginData.Token, nonExistentAdID, http.StatusBadRequest)
		require.NotEmpty(s.T(), gjson.Get(respBody, "error").String(), "Должно вернуться сообщение об ошибке")
	})

	s.Run("Удаление из избранного несуществующего ad_id (404)", func() {
		respBody := favManager.RemoveFromFavorites(s.T(), s.loginData.Token, nonExistentAdID, http.StatusNotFound)
		require.NotEmpty(s.T(), gjson.Get(respBody, "error").String(), "Должно вернуться сообщение об ошибке")
	})

	s.Run("Запросы к ручкам избранного без авторизации возвращают 401 Unauthorized", func() {
		respGet := favManager.GetFavorites(s.T(), "invalid_token", http.StatusUnauthorized)
		require.NotEmpty(s.T(), gjson.Get(respGet, "error").String(), "Должно вернуться сообщение об ошибке")

		respPost := favManager.AddToFavorites(s.T(), "invalid_token", nonExistentAdID, http.StatusUnauthorized)
		require.NotEmpty(s.T(), gjson.Get(respPost, "error").String(), "Должно вернуться сообщение об ошибке")

		respDel := favManager.RemoveFromFavorites(s.T(), "invalid_token", nonExistentAdID, http.StatusUnauthorized)
		require.NotEmpty(s.T(), gjson.Get(respDel, "error").String(), "Должно вернуться сообщение об ошибке")
	})

	s.Run("Бизнес-ограничение: попытка добавить своё собственное объявление в избранное возвращает 403 Forbidden", func() {
		ownAdID := s.createAdByToken(s.loginData.Token)
		defer advManager.DeleteAdvertisement(s.T(), s.loginData.Token, ownAdID, http.StatusNoContent)

		respOwn := favManager.AddToFavorites(s.T(), s.loginData.Token, ownAdID, http.StatusForbidden)
		require.NotEmpty(s.T(), gjson.Get(respOwn, "error").String(), "Должно вернуться сообщение об ошибке 403")
	})

	s.Run("Бизнес-ограничение: повторное добавление того же объявления возвращает 409 Conflict", func() {
		otherAdID := s.createAdByToken(s.secondLoginData.Token)
		defer advManager.DeleteAdvertisement(s.T(), s.secondLoginData.Token, otherAdID, http.StatusNoContent)

		favManager.AddToFavorites(s.T(), s.loginData.Token, otherAdID, http.StatusCreated)

		respDuplicate := favManager.AddToFavorites(s.T(), s.loginData.Token, otherAdID, http.StatusConflict)
		require.NotEmpty(s.T(), gjson.Get(respDuplicate, "error").String(), "Должно вернуться сообщение об ошибке 409")
	})
}
