// handlers/user_handlers.go

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andreevym/gofermart/internal/logger"
	"github.com/andreevym/gofermart/internal/services"
	"go.uber.org/zap"
)

type AuthDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// PostRegisterUser регистрация пользователя
// #### **Регистрация пользователя**
//
// Хендлер: `POST /api/user/register`.
//
// Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным. После успешной регистрации должна происходить автоматическая аутентификация пользователя.
//
// Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок `Authorization`.
//
// Формат запроса:
//
// # Скопировать код
//
// POST /api/user/register HTTP/1.1
// Content-Type: application/json
// ...
//
// {
// "login": "<login>",
// "password": "<password>"
// }
//
// Возможные коды ответа:
//
// *   `200` — пользователь успешно зарегистрирован и аутентифицирован;
// *   `400` — неверный формат запроса;
// *   `409` — логин уже занят;
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) PostRegisterUser(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	a := AuthDTO{}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.authService.Register(r.Context(), a.Login, a.Password)
	if err != nil {
		logger.Logger().Warn("authService.Register", zap.Error(err))
		if errors.Is(err, services.ErrAuthAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
		} else if err != nil && !errors.Is(err, services.ErrAuthAlreadyExists) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PostLoginUser аутентификация пользователя
// #### **Аутентификация пользователя**
//
// Хендлер: `POST /api/user/login`.
//
// Аутентификация производится по паре логин/пароль.
//
// Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок `Authorization`.
//
// Формат запроса:
//
// # Скопировать код
//
// POST /api/user/login HTTP/1.1
// Content-Type: application/json
// ...
//
// {
// "login": "<login>",
// "password": "<password>"
// }
//
// Возможные коды ответа:
//
// *   `200` — пользователь успешно аутентифицирован;
// *   `400` — неверный формат запроса;
// *   `401` — неверная пара логин/пароль;
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) PostLoginUser(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	a := AuthDTO{}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, err := h.authService.Login(r.Context(), a.Login, a.Password)
	if err != nil {
		logger.Logger().Warn("authService.Login", zap.Error(err))
		if errors.Is(err, services.ErrAuthBadCredentials) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", authToken))

	w.WriteHeader(http.StatusOK)
}
