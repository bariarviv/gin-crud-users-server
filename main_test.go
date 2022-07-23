package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin_CRUD_server/db"
	"gin_CRUD_server/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	TestUser = models.NewUser(TestEmail, "bari", "1234")
	MapDB    = db.TestMapOps{Name: "Map DB Test", Users: make(map[string]models.User)}
)

const (
	SlashSeparator = "/"
	KeyFileTest    = "ssl.key"
	CertFileTest   = "ssl.crt"
	ContentType    = "Content-Type"
	TestEmail      = "bari@gmail.com"
)

func Test_createTLSCert(t *testing.T) {
	tests := []struct {
		name     string
		certFile string
		keyFile  string
		port     string
		wantErr  bool
	}{
		{"Failed to create due to empty certFile", "", KeyFileTest, Port, true},
		{"Failed to create due to empty keyFile", CertFileTest, "", Port, true},
		{"Creates TLS cert successfully", CertFileTest, KeyFileTest, Port, false},
		{"Failed to create due to empty port", CertFileTest, KeyFileTest, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, err := createTLSCert(tt.certFile, tt.keyFile, tt.port)
			if ln == nil && ((err != nil) != tt.wantErr) {
				t.Errorf("createTLSCert() ln = %v, wantErr %v\n", err, tt.wantErr)
			}
		})
	}
}

func TestListUsersHandler(t *testing.T) {
	DBApi = MapDB
	tests := []struct {
		name     string
		url      string
		wantCode int
	}{
		{"Failed to get users list due to incorrect URL", URL, http.StatusNotFound},
		{"Gets fail due to empty users map", ListURL, http.StatusInternalServerError},
		{"Gets users list successfully (if there are users in the folder)", ListURL, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates gin router & http test
			respRecorder, router := createRouterAndWriter()
			router.GET(ListURL, ListUsersHandler)
			if tt.wantCode == http.StatusOK {
				DBApi.InsertNewUser(*TestUser)
			}
			// Creates a request
			request, err := createNewRequest(http.MethodGet, tt.url, "", nil)
			if err != nil {
				t.Errorf(err.Error())
			}
			router.ServeHTTP(respRecorder, request)
			assert.Equal(t, respRecorder.Code, tt.wantCode)
			// Prints the users map
			printUsersMap()
		})
	}
}

func TestAddUserHandler(t *testing.T) {
	DBApi = MapDB
	tests := []struct {
		name     string
		user     *models.User
		url      string
		wantCode int
	}{
		{"Adds a new user successfully", TestUser, URL, http.StatusOK},
		{"Adds fail due to incorrect user (nil)", nil, URL, http.StatusBadRequest},
		{"Adds fail due to incorrect URL", TestUser, SlashSeparator, http.StatusNotFound},
		{"Adds fail due to incorrect empty user", &models.User{}, URL, http.StatusBadRequest},
		{"Adds fail due to the invalid email", models.NewUser("abc", "bari", "1234"), URL, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates gin router & http test
			respRecorder, router := createRouterAndWriter()
			router.PUT(URL, AddUserHandler)
			// Prints the users map
			printUsersMap()
			// Creates a request
			request, err := newBindJSONRequest(tt.user, tt.url, http.MethodPut)
			if err != nil {
				t.Errorf("newBindJSONRequest Error: %v\n", err)
			}
			router.ServeHTTP(respRecorder, request)
			assert.Equal(t, respRecorder.Code, tt.wantCode)
			if tt.user != nil {
				printResults(fmt.Sprintf("## %s added successfully:", tt.user.Email), tt.wantCode)
			}
		})
	}
}

func TestGetUserHandler(t *testing.T) {
	DBApi = MapDB
	DBApi.InsertNewUser(*TestUser)

	tests := []struct {
		name     string
		url      string
		email    string
		wantCode int
	}{
		{"Gets fail due to empty email", URL, "", http.StatusBadRequest},
		{"Gets fail due to the invalid email", URL, "abc", http.StatusBadRequest},
		{"Gets an existing user successfully", URL, TestUser.Email, http.StatusOK},
		{"Gets fail due to the user doesn't exist", URL, "a@gmail.com", http.StatusNotFound},
		{"Gets fail due to incorrect URL", SlashSeparator, TestUser.Email, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates gin router & http test
			respRecorder, router := createRouterAndWriter()
			router.GET(URL, GetUserHandler)
			// Prints the users map
			printUsersMap()
			// Performs the request
			request, err := newFileUploadRequest(tt.email, tt.url, http.MethodGet)
			if err != nil {
				t.Errorf("newFileUploadRequest Error: %v\n", err)
			}
			router.ServeHTTP(respRecorder, request)
			assert.Equal(t, respRecorder.Code, tt.wantCode)
			printResults(fmt.Sprintf("## %s gets successfully:", tt.email), tt.wantCode)
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	DBApi = MapDB
	DBApi.InsertNewUser(*TestUser)

	tests := []struct {
		name     string
		user     *models.User
		url      string
		wantCode int
	}{
		{"Updates fail due to incorrect user (nil)", nil, URL, http.StatusBadRequest},
		{"Updates fail due to incorrect URL", TestUser, SlashSeparator, http.StatusNotFound},
		{"Updates fail due to incorrect empty user", &models.User{}, URL, http.StatusBadRequest},
		{"Updates an existing user successfully", models.NewUser(TestEmail, "bari2", "12345"), URL, http.StatusOK},
		{"Updates fail due to incorrect email", models.NewUser(TestEmail+"abc", "bari", "1234"), URL, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates gin router & http test
			respRecorder, router := createRouterAndWriter()
			router.POST(URL, UpdateUserHandler)
			// Prints the users map
			printUsersMap()
			// Creates a request
			request, err := newBindJSONRequest(tt.user, tt.url, http.MethodPost)
			if err != nil {
				t.Errorf("newBindJSONRequest Error: %v\n", err)
			}
			router.ServeHTTP(respRecorder, request)
			assert.Equal(t, respRecorder.Code, tt.wantCode)
			if tt.user != nil {
				printResults(fmt.Sprintf("## %s updated successfully:", tt.user.Email), tt.wantCode)
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	DBApi = MapDB
	DBApi.InsertNewUser(*TestUser)

	tests := []struct {
		name     string
		email    string
		url      string
		wantCode int
	}{
		{"Deletes fail due to empty email", "", URL, http.StatusBadRequest},
		{"Deletes user successfully", TestUser.Email, URL, http.StatusOK},
		{"Deletes fail due to the invalid email", "abc", URL, http.StatusBadRequest},
		{"Deletes fail due to incorrect URL", TestUser.Email, SlashSeparator, http.StatusNotFound},
		{"Deletes fail due to the user doesn't exist", TestUser.Email + "abc", URL, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creates gin router & http test
			respRecorder, router := createRouterAndWriter()
			router.DELETE(URL, DeleteUserHandler)
			// Prints the users map
			printUsersMap()
			// Performs the request
			request, err := newFileUploadRequest(tt.email, tt.url, http.MethodDelete)
			if err != nil {
				t.Errorf("newFileUploadRequest Error: %v\n", err)
			}
			router.ServeHTTP(respRecorder, request)
			assert.Equal(t, respRecorder.Code, tt.wantCode)
			printResults(fmt.Sprintf("## %s deleted successfully:", tt.email), tt.wantCode)
		})
	}
}

func Test_getStatusAndMsgErr(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"Gets status code 404", sql.ErrNoRows, http.StatusNotFound},
		{"Gets status code 500", fmt.Errorf("error"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := getStatusAndMsgErr(tt.err)
			if got != tt.wantStatus {
				t.Errorf("getStatusAndMsg() got = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func Test_setupRouter(t *testing.T) {
	go func() {
		if _, err := setupRouter(CertFileTest, KeyFileTest, ":3001"); err != nil {
			t.Errorf("setupRouter() error = %v\n", err)
		}
	}()
}

// createRouterAndWriter creates gin router & http test record
func createRouterAndWriter() (*httptest.ResponseRecorder, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	respRecorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(respRecorder)
	return respRecorder, router
}

// createNewRequest creates an HTTP request
func createNewRequest(method, url, contentType string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return request, fmt.Errorf("http.NewRequest() Error: %v", err)
	}
	if contentType != "" {
		request.Header.Set(ContentType, contentType)
	}
	return request, nil
}

// newFileUploadRequest creates an HTTP request and sends an email using FormData
func newFileUploadRequest(email, url, method string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField(FieldName, email)
	defer writer.Close()
	// Creates a request
	return createNewRequest(method, url, writer.FormDataContentType(), body)
}

// newFileUploadRequest creates an HTTP request and adds JSON to the body
func newBindJSONRequest(user *models.User, url, method string) (*http.Request, error) {
	buf, err := json.Marshal(&user)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal() Error: %v\n", err)
	}
	// Creates a request
	return createNewRequest(method, url, "application/json", bytes.NewBuffer(buf))
}

// printResults prints the results if the code is 200 OK
func printResults(msg string, code int) {
	if code == http.StatusOK {
		fmt.Println(msg)
		// Prints the users map
		printUsersMap()
	}
}

// printUsersMap prints the current users map
func printUsersMap() {
	if len(MapDB.Users) == 0 {
		fmt.Println("The users map is empty")
		return
	}
	fmt.Println("Current users map:")
	for _, user := range MapDB.Users {
		fmt.Printf("* Email: %s, Username: %s, Password: %s\n", user.Email, user.Name, user.Password)
	}
}
