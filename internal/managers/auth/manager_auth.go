package auth

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// RegisterAndLogin создает нового пользователя (если он еще не создан) и возвращает данные авторизации.
func RegisterAndLogin(t *testing.T, email, password string) models.LoginOkResponse {
	req := map[string]string{
		"email":      email,
		"password":   password,
		"first_name": "SecondUser",
	}
	raw := auth.HttpPostAuthRegister(t, req)
	_ = raw.Body.Close() // Статус не проверяем: 201 при первом запуске или 400/409, если уже создан ранее

	return Login(t, email, password)
}
