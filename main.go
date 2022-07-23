package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"

	"gin_CRUD_server/db"
	"gin_CRUD_server/models"
	"github.com/gin-gonic/gin"
)

const (
	DBPort    = "5432"
	Port      = ":3000"
	FieldName = "email"
	URL       = "/user"
	ListURL   = "/users"
	Host      = "database"
	CertFile  = "/etc/ssl/certs/ssl.crt"
	KeyFile   = "/etc/ssl/certs/ssl.key"
)

var (
	DBApi models.DBOps
)

func main() {
	// Setups the DB instance
	setupDB(Host, DBPort)
	if _, err := setupRouter(CertFile, KeyFile, Port); err != nil {
		fmt.Println(err)
		return
	}
}

// setupDB setups the DB instance
func setupDB(host, port string) {
	dbName := os.Getenv("POSTGRES_DB")
	dbSSL := os.Getenv("POSTGRES_SSL")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	config := db.NewConfig(host, port, dbUser, dbPassword, dbName, dbSSL)
	db.GetInstance(config)
	DBApi = db.SqlOps{Name: "SQL Server"}
}

// setupRouter setups the server and the routers according to the HTTP requests
func setupRouter(certFile, keyFile, port string) (*gin.Engine, error) {
	// Setups the server and the routers according to the HTTP requests
	router := gin.Default()
	router.PUT(URL, AddUserHandler)
	router.GET(URL, GetUserHandler)
	router.POST(URL, UpdateUserHandler)
	router.DELETE(URL, DeleteUserHandler)
	router.GET(ListURL, ListUsersHandler)

	// Creates tls certificate
	ln, err := createTLSCert(certFile, keyFile, port)
	if err != nil || ln == nil {
		return router, err
	}
	// Starts server with https/ssl enabled on http://localhost:Port
	log.Fatal(router.RunListener(*ln))
	return router, nil
}

// createTLSCert creates tls certificate
func createTLSCert(certFile, keyFile, port string) (*net.Listener, error) {
	// Creates tls certificate
	certs, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Printf("Cannot load TLS certificate from certFile=%q, keyFile=%q: %s\n", certFile, keyFile, err)
		return nil, err
	}
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	// Creates custom listener
	ln, err := tls.Listen("tcp", port, &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{certs},
	})
	if err != nil {
		return &ln, err
	}
	return &ln, nil
}

// AddUserHandler adds a new user
func AddUserHandler(ctx *gin.Context) {
	user, err := getUserFromBindJSON(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = DBApi.InsertNewUser(*user); err != nil {
		ctx.String(getStatusAndMsgErr(err))
		return
	}
	ctx.String(http.StatusOK, user.Email+" added successfully!\n")
}

// GetUserHandler returns the user according to the email received
func GetUserHandler(ctx *gin.Context) {
	// Gets the email from the form-data
	email, err := getEmail(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	user, err := DBApi.IsExistsInUsersTable(email)
	if err != nil {
		ctx.String(getStatusAndMsgErr(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateUserHandler updates username & password of an existing user
func UpdateUserHandler(ctx *gin.Context) {
	user, err := getUserFromBindJSON(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = DBApi.UpdateNameAndPassUser(*user); err != nil {
		ctx.String(getStatusAndMsgErr(err))
		return
	}
	ctx.String(http.StatusOK, user.Email+" updated successfully!\n")
}

// DeleteUserHandler deletes an existing user
func DeleteUserHandler(ctx *gin.Context) {
	// Gets the email from the form-data
	email, err := getEmail(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = DBApi.DeleteUser(email); err != nil {
		ctx.String(getStatusAndMsgErr(err))
		return
	}
	ctx.String(http.StatusOK, email+" deleted successfully!\n")
}

// ListUsersHandler returns a JSON array with the list of all the users
func ListUsersHandler(ctx *gin.Context) {
	users, err := DBApi.GetAllUsers()
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s\n", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// getUserFromBindJSON binds the received JSON to user
func getUserFromBindJSON(ctx *gin.Context) (*models.User, error) {
	user := models.User{}
	if err := ctx.BindJSON(&user); err != nil {
		return &user, fmt.Errorf("ctx.BindJSON() Error: %s\n", err.Error())
	}
	if user.Email == "" || user.Name == "" || user.Password == "" {
		return &user, fmt.Errorf("Please try again and enter email, username and password\n")
	}
	return &user, nil
}

// getEmail returns the email from the form-data
func getEmail(ctx *gin.Context) (string, error) {
	email := ctx.PostForm(FieldName)
	if email == "" {
		return email, fmt.Errorf("Please add an email to the form-data request\n")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return email, err
	}
	return email, nil
}

// getStatusAndMsgErr returns the status code and message error according to the error received
func getStatusAndMsgErr(err error) (int, string) {
	var status int
	if err == sql.ErrNoRows {
		status = http.StatusNotFound
	} else {
		status = http.StatusInternalServerError
	}
	return status, fmt.Sprintf("Error: %s\n", err.Error())
}
