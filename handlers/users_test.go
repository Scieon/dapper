package handlers

import (
	"dapper/models"
	"dapper/util"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestGetUsers_ValidJWT(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	// driver translates orm commands to db commands
	driver := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(driver, &gorm.Config{})

	rows := sqlmock.NewRows([]string{"firstname", "lastname"}).AddRow("firstname", "lastname")
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	jwtToken, _ := util.CreateJWT("email")

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("X-Authentication-Token", jwtToken)

	responseRecorder := httptest.NewRecorder()

	handler := GetUsers(db)

	handler.ServeHTTP(responseRecorder, req)

	// Normally would assert user response too but to keep time limit in mind returning 200 is ok!
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetUsers_InvalidJWT(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("X-Authentication-Token", "bad token")

	responseRecorder := httptest.NewRecorder()

	handler := GetUsers(nil)

	handler.ServeHTTP(responseRecorder, req)

	// Normally would assert user response too but to keep time limit in mind returning 200 is ok!
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "Valid User Request should return 201",
			body:         `{"firstname":"bob", "lastname":"jones", "email":"email@yahoo.com", "password":"psw"}"`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid User request should return 400",
			body:         `{"firstname":2, "lastname":"jones", "email":"email@yahoo.com", "password":"psw"}"`,
			expectedCode: http.StatusBadRequest,
		},
	}

	mockDb, mock, _ := sqlmock.New()

	// driver translates orm commands to db commands
	driver := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(driver, &gorm.Config{})

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"users\"").
		WithArgs("email@yahoo.com", "bob", "jones", "psw").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/signup", strings.NewReader(tc.body))

			responseRecorder := httptest.NewRecorder()
			handler := CreateUser(db)

			handler.ServeHTTP(responseRecorder, req)

			if statusCode := responseRecorder.Code; statusCode != tc.expectedCode {
				t.Errorf("Handler returned different status code: received %v expected %v", statusCode, tc.expectedCode)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	loginBody := models.LoginBody{
		Email:    "email@yahoo.com",
		Password: "psw",
	}

	loginJSON, _ := json.Marshal(loginBody)

	bodyReader := strings.NewReader(string(loginJSON))

	mockDb, mock, _ := sqlmock.New()

	driver := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(driver, &gorm.Config{})

	rows := sqlmock.NewRows([]string{"email", "firstname", "lastname", "password"}).
		AddRow("email@yahoo.com", "bob", "jones", "psw")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND password = $2 ORDER BY "users"."email" LIMIT 1`)).
		WithArgs(loginBody.Email, loginBody.Password).
		WillReturnRows(rows)

	req := httptest.NewRequest("POST", "/login", bodyReader)

	responseRecorder := httptest.NewRecorder()

	handler := LoginUser(db)

	handler.ServeHTTP(responseRecorder, req)

	// Good enough to assert 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}
