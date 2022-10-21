package controller

import (
	"bytes"
	"diary_api/database"
	"diary_api/middleware"
	"diary_api/model"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func router() *gin.Engine {
	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", Register)
	publicRoutes.POST("/login", Login)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("/entry", AddEntry)
	protectedRoutes.GET("/entry", GetAllEntries)

	return router
}

func setup() {

	err := godotenv.Load("../.env.test.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Entry{})
}

func teardown() {
	migrator := database.Database.Migrator()
	migrator.DropTable(&model.User{})
	migrator.DropTable(&model.Entry{})
}

func makeRequest(method, url string, body interface{}, isAuthenticatedRequest bool) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+bearerToken())
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func bearerToken() string {
	user := model.AuthenticationInput{
		Username: "yemiwebby",
		Password: "test",
	}

	writer := makeRequest("POST", "/auth/login", user, false)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["jwt"]
}