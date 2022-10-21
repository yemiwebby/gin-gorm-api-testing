package controller

import (
	"diary_api/model"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRegister(t *testing.T) {
	newUser := model.AuthenticationInput{
		Username: "yemiwebby",
		Password: "test",
	}
	writer := makeRequest("POST", "/auth/register", newUser, false)
	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestLogin(t *testing.T) {
	user := model.AuthenticationInput{
		Username: "yemiwebby",
		Password: "test",
	}

	writer := makeRequest("POST", "/auth/login", user, false)

	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	_, exists := response["jwt"]

	assert.Equal(t, true, exists)
}

func TestIncompleteLoginRequest(t *testing.T) {
	request := map[string]string{"username": "yemiwebby"}
	writer := makeRequest("POST", "/auth/login", request, true)
	assert.Equal(t, http.StatusBadRequest, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	message, exists := response["error"]

	assert.Equal(t, true, exists)
	assert.Equal(t, message, "Password not provided")
}