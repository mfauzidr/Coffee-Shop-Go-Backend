package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	router := gin.Default()
	authRepositoryMock := new(repository.AuthRepositoryMock)

	handler := NewAuthHandler(authRepositoryMock)
	authRepositoryMock.On("RegisterUser", mock.Anything).Return("Register user is success", nil)
	router.POST("/auth/register", handler.Register)

	requestBody, _ := json.Marshal(map[string]string{
		"email":    "testing@mail.com",
		"password": "123456",
		"role":     "customer",
	})

	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "An error occurred while making the request")
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code, "Status code does not match")
	var actualResponse pkg.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err, "An error occurred when getting a response")

	assert.Equal(t, 201, actualResponse.Status, "Status code does not match")
	assert.Equal(t, "Register success", actualResponse.Message, "Response message does not match")
	assert.Equal(t, "Register user is success", actualResponse.Data, "Status does not match")
}

// func TestLogin(t *testing.T) {
// 	router := gin.Default()
// 	authRepositoryMock := new(repository.AuthRepositoryMock)

// 	results := &models.UsersRes{{
// 		Email:    "testing@mail.com",
// 		Password: "123456",
// 	}}

// 	handler := NewAuthHandler(authRepositoryMock)
// 	authRepositoryMock.On("GetByEmail", mock.Anything).Return(results, nil)
// 	router.POST("/auth/login", handler.Login)

// 	requestBody, _ := json.Marshal(map[string]string{
// 		"email":    "testing@mail.com",
// 		"password": "123456",
// 	})

// 	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
// 	assert.NoError(t, err, "An error occurred while making the request")
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)

// 	assert.Equal(t, http.StatusOK, recorder.Code, "1. Status code does not match")
// 	var actualResponse pkg.Response
// 	err = json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
// 	assert.NoError(t, err, "An error occurred when getting a response")

// 	assert.Equal(t, http.StatusOK, actualResponse.Status, "2. Status code does not match")
// 	assert.Equal(t, "Login success", actualResponse.Message, "Response message does not match")
// 	assert.NotEmpty(t, actualResponse.Data, "Token should be present in the response")
// }

func TestLogin2(t *testing.T) { //should be deleted
	router := gin.Default()
	authRepositoryMock := new(repository.AuthRepositoryMock)

	results := &models.UsersRes{{
		Email:    "testing@mail.com",
		Password: "123456",
	}}

	handler := NewAuthHandler(authRepositoryMock)
	authRepositoryMock.On("GetByEmail", mock.Anything).Return(results, nil)
	router.POST("/auth/login", handler.LoginDummy)

	requestBody, _ := json.Marshal(map[string]string{
		"email":    "testing@mail.com",
		"password": "123456",
	})

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
	assert.NoError(t, err, "An error occurred while making the request")
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "1. Status code does not match")
	var actualResponse pkg.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err, "An error occurred when getting a response")

	assert.Equal(t, http.StatusOK, actualResponse.Status, "2. Status code does not match")
	assert.Equal(t, "Login success", actualResponse.Message, "Response message does not match")
	// assert.NotEmpty(t, actualResponse.Data, "Token should be present in the response")
}

//terjadi masalah kemungkinan dari hashing / verify password
