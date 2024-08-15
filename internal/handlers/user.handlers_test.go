package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"
// 	"github.com/mfauzidr/coffeeshop-go-backend/internal/repository"
// 	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestCreateUser(t *testing.T) {
// 	router := gin.Default()
// 	userRepositoryMock := new(repository.UserRepositoryMock)

// 	results := &models.UsersRes{{
// 		Email:    "testing@mail.com",
// 		Password: "123456",
// 	}}

// 	handler := NewUserRepository(userRepositoryMock)
// 	userRepositoryMock.On("CreateUser", mock.Anything).Return(results, nil)
// 	router.POST("/user/login", handler.CreateUser())

// 	requestBody, _ := json.Marshal(map[string]string{
// 		"email":    "testing@mail.com",
// 		"password": "123456",
// 	})

// 	req, err := http.NewRequest("POST", "/user/login", bytes.NewBuffer(requestBody))
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
