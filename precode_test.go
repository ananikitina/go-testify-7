package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOkAndBodyNotEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка, что код ответа 200
	require.Equal(t, http.StatusOK, responseRecorder.Code)

	// Проверка, что тело ответа не пустое
	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body)
}

func TestMainHandlerWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=unknown_city", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка, что код ответа 400
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	// Проверка, что тело ответа содержит ошибку "wrong city value"
	body := responseRecorder.Body.String()
	assert.Contains(t, body, "wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCafes := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	actualCafes := strings.Split(body, ",")

	// Проверка, что количество кафе не превышено
	require.LessOrEqual(t, len(actualCafes), totalCafes)
	// Проверка, что возвращено общее количество кафе
	assert.Len(t, actualCafes, totalCafes)
	// Проверка, что возвращены все кафе
	expectedCafes := []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"}
	assert.ElementsMatch(t, expectedCafes, actualCafes)
}
