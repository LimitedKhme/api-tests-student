package favorites

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddToFavorites добавляет объявление в избранное и проверяет статус-код ответа.
func AddToFavorites(t *testing.T, token, adID string, expectedStatus int) string {
	req := models.AddFavoriteRequest{AdvertisementID: adID}
	raw := client.HttpPostFavorite(t, token, req)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус-код при добавлении в избранное")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	return string(body)
}

// GetFavorites получает список избранного и проверяет статус-код ответа.
func GetFavorites(t *testing.T, token string, expectedStatus int) string {
	raw := client.HttpGetFavorites(t, token)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус-код при получении избранного")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	return string(body)
}

// RemoveFromFavorites удаляет объявление из избранного и проверяет статус-код ответа.
func RemoveFromFavorites(t *testing.T, token, adID string, expectedStatus int) string {
	raw := client.HttpDeleteFavorite(t, token, adID)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус-код при удалении из избранного")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	return string(body)
}
