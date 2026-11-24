package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// AssertAPIError проверяет поля error и message в ответе API при ошибке.
func AssertAPIError(t *testing.T, respBody, expectedError, expectedMessageSubstr string) {
	t.Helper()

	actualError := gjson.Get(respBody, "error").String()
	actualMessage := gjson.Get(respBody, "message").String()

	require.NotEmpty(t, actualError, "Ожидалось непустое поле error. Ответ: %s", respBody)
	if expectedError != "" {
		require.Equal(t, expectedError, actualError, "Неожиданное значение error. Ответ: %s", respBody)
	}

	require.NotEmpty(t, actualMessage, "Ожидалось непустое поле message. Ответ: %s", respBody)
	if expectedMessageSubstr != "" {
		require.Contains(t, actualMessage, expectedMessageSubstr, "Поле message не содержит ожидаемую подстроку. Ответ: %s", respBody)
	}
}
