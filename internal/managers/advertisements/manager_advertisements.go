package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	client "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

// GetAdvertisementsList получает список объявлений по поисковому запросу
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := client.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.AdvertisementsList
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

// Создание объявления с проверкой статуса
func CreateAdvertisement(t *testing.T, token string, params models.CreateAdvertisementFormParams, expectedStatus int) string {
	raw := client.HttpPostAdvertisement(t, token, params)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус код при создании")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	// Возвращаем сырой ответ (например, там только ID или структура)
	return string(body)
}

// Получение объявления по ID
func GetAdvertisement(t *testing.T, id string, expectedStatus int) models.Advertisement {
	raw := client.HttpGetAdvertisement(t, id)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode)

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.Advertisement
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

// Удаление объявления
func DeleteAdvertisement(t *testing.T, token string, id string, expectedStatus int) {
	raw := client.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode)
}

// GetAdvertisementPhotos возвращает сырой JSON-ответ со списком фотографий объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string, expectedStatus int) string {
	raw := client.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, expectedStatus, raw.StatusCode, "Неверный статус-код при получении фотографий")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	return string(body)
}
