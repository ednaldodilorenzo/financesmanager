package auth_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/server"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct {
}

func (a *MockAuthService) ExecuteAuthentication(username string, password string) (*model.User, error) {
	if username == "test@test.com" && password == "testing" {
		return &model.User{
			ID:        1,
			Name:      "Test",
			Email:     "test@test.com",
			Password:  &password,
			CreatedAt: nil,
			UpdatedAt: nil,
		}, nil
	}

	return nil, errors.New("password does not match")
}

func (a *MockAuthService) RegisterUser(user *model.User) error {
	if user.Name == "Service Error" {
		return errors.New("Test error")
	} else if user.Name == "Duplicate User" {
		return errors.New("duplicate key value violates unique")
	}

	return nil
}

var svr *server.Server

// You can use testing.T, if you want to test the code without benchmarking
func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	svr = &server.Server{
		App: server.InitFiberApplication(),
	}
	//mockAuthService := &MockAuthService{}
	//authController := auth.BuildAuthController(mockAuthService)
	//svr.BasicSetup("/auth", func(router fiber.Router) {
	//	router.Post("/login", authController.SigninUser)
	//	router.Post("/signup", authController.SignUpUser)
	//})

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

func TestSigninUser(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		path         string
		description  string
		method       string
		body         io.Reader
		expectedCode int
	}{
		{
			path:         "/api/auth/login",
			description:  "Test login bad request",
			method:       "POST",
			body:         nil,
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/login",
			description:  "Test login email not registered",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"email": "notregister@test.com", "password": "testing"}`)),
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			path:         "/api/auth/login",
			description:  "Test login wrong password",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"email": "test@test.com", "password": "123"}`)),
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			path:         "/api/auth/login",
			description:  "Test login ok",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"email": "test@test.com", "password": "testing"}`)),
			expectedCode: fiber.StatusOK,
		},
		{
			path:         "/api/auth/logon",
			description:  "Test login invalid address",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"email": "test@test.com", "password": "testing"}`)),
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.path, tt.body)
		req.Header.Set("Content-Type", "application/json")
		res, _ := svr.App.Test(req, -1)
		body, err := ioutil.ReadAll(res.Body)
		log.Println(string(body))
		assert.Nil(t, err, "Error making request")
		assert.Equalf(t, tt.expectedCode, res.StatusCode, tt.description)
		defer res.Body.Close()
	}
}

func TestSignUpUser(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		path         string
		description  string
		method       string
		body         io.Reader
		expectedCode int
	}{
		{
			path:         "/api/auth/signup",
			description:  "Test register ok",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Test", "email": "test@test.com", "password": "testingpass", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusCreated,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register mandatory fields not filled",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"email": "test@test.com", "password": "testingpass", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register mandatory fields not filled",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Test", "password": "testingpass", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register mandatory fields not filled",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Test", "email": "test@test.com", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register mandatory fields not filled",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Test", "email": "test@test.com", "password": "testingpass"}`)),
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register Password diff from confirm",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Test", "email": "test@test.com", "password": "testingpass", "confirmPassword": "differnetpass"}`)),
			expectedCode: fiber.StatusBadRequest,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register service error",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Service Error", "email": "test@test.com", "password": "testingpass", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			path:         "/api/auth/signup",
			description:  "Test register service error",
			method:       "POST",
			body:         bytes.NewReader([]byte(`{"name": "Duplicate User", "email": "test@test.com", "password": "testingpass", "confirmPassword": "testingpass"}`)),
			expectedCode: fiber.StatusConflict,
		},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.path, tt.body)
		req.Header.Set("Content-Type", "application/json")
		res, _ := svr.App.Test(req, -1)
		body, err := ioutil.ReadAll(res.Body)
		log.Println(string(body))
		assert.Nil(t, err, "Error making request")
		assert.Equalf(t, tt.expectedCode, res.StatusCode, tt.description)
		defer res.Body.Close()
	}
}
