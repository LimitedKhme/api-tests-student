package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"
	"github.com/stretchr/testify/require"
)

// ContainsAdvertisementByTitle проверяет наличие объявления с заданным заголовком в списке
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()
	found := false
	for _, item := range items {
		if item.Title == title {
			found = true
			break
		}
	}
	require.True(t, found, "Объявление с заголовком '%s' не найдено в списке", title)
}