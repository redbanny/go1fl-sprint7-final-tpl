package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	coumtCafeInCyte := len(cafeList["moscow"])

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
		city  string
	}{
		{0, 0, "moscow"},
		{1, 1, "moscow"},
		{2, 2, "moscow"},
		{100, coumtCafeInCyte, "moscow"},
	}

	for _, w := range requests {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", w.count, w.city), nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		
		cafes := strings.Split(response.Body.String(), ",")		
		var cafeCount int
		if len(cafes) == 1 && cafes[0] == "" {
			cafeCount = 0
		} else {
			cafeCount = len(cafes)
		}
		
		assert.Equal(t, w.want, cafeCount)
	}
}


func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string   // передаваемое значение search 
        wantCount int      // ожидаемое количество кафе в ответе
		city  string
	}{
		{"фасоль", 0, "moscow"},
		{"кофе", 2, "moscow"},
		{"вилка", 1, "moscow"},
	}

	for _, w := range requests {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&search=%s", w.city, w.search), nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		
		cafes := strings.Split(response.Body.String(), ",")		
		var cafeCount int
		if len(cafes) == 1 && cafes[0] == "" {
			cafeCount = 0
		} else {
			cafeCount = len(cafes)
		}
		
		assert.Equal(t, w.wantCount, cafeCount)
	}
}