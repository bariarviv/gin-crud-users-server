# GoGin - CRUD Server
![Coverage](https://img.shields.io/badge/Coverage-64.6%%25-brightgreen)

## Gin Golang
[Gin](https://github.com/gin-gonic/gin) is a web framework written in Go (Golang). It features a martini-like API with performance that is up to 40 times faster thanks to httprouter.

### Advantages:
1. Simple code - can save the developer a lot of time in large projects.
2. High performance.
3. Error management.
4. Easy JSON authentication.
5. Gin has a mode test.

### Disadvantages:
1. Not flexible in development.
2. Can add more complexity to your project and slow down the development time, the infrastructure now depends on the package that other people maintain.

### Alternative Solutions:
1. [net/http:](https://github.com/golang/go) it's easier to use and can handle more cases.
2. [fasthttp:](https://github.com/valyala/fasthttp) was designed for some high-performance edge cases.
3. [fiber:](https://github.com/gofiber/fiber) built on top of the fasthttp HTTP engine for Go, which is the fastest HTTP engine for Go.
4. [echo:](https://github.com/labstack/echo) supports HTTP/2 for faster performance and an overall better user experience and has automatic TLS certificates.


## Details of the application
Backend Golang application that has the following routes (using gin):
* ***PUT    /user  -*** add a new user, you need to add a JSON including email, username, and password in the request body.
* ***GET    /user  -*** to get an existing user, you need to add an email in the request form-data.
* ***POST   /user  -*** to update a username and password for an existing user, you need to add a JSON including email, username, and password in the request body.
* ***DELETE /user  -*** to delete an existing user, you need to add an email in the request form-data.
* ***GET    /users -*** returns a JSON array with the list of users.


## Running Steps
### Step 1 - Build the docker image:
```
docker build -t ginserver .
```

### Step 2 - Run the docker image:
```
docker run -it --rm -p 3000:3000 ginserver
```

#### For example:
```
[baria@ ~]$ sudo docker run -it --rm -p 3000:3000 ginserver 
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] PUT    /user                     --> main.AddUserHandler (3 handlers)
[GIN-debug] GET    /user                     --> main.GetUserHandler (3 handlers)
[GIN-debug] POST   /user                     --> main.UpdateUserHandler (3 handlers)
[GIN-debug] DELETE /user                     --> main.DeleteUserHandler (3 handlers)
[GIN-debug] GET    /users                    --> main.ListUsersHandler (3 handlers)
[GIN-debug] Listening and serving HTTP on listener what's bind with address@[::]:3000
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN] 2022/07/20 - 04:13:34 | 200 |     237.023µs |      172.17.0.1 | GET      "/users"
```


## Unit Tests Output
```
=== RUN   Test_createTLSCert
=== RUN   Test_createTLSCert/Failed_to_create_due_to_empty_certFile
Cannot load TLS certificate from certFile="", keyFile="ssl.key": open : no such file or directory
=== RUN   Test_createTLSCert/Failed_to_create_due_to_empty_keyFile
Cannot load TLS certificate from certFile="ssl.crt", keyFile="": open : no such file or directory
=== RUN   Test_createTLSCert/Creates_TLS_cert_successfully
=== RUN   Test_createTLSCert/Failed_to_create_due_to_empty_port
--- PASS: Test_createTLSCert (0.00s)
    --- PASS: Test_createTLSCert/Failed_to_create_due_to_empty_certFile (0.00s)
    --- PASS: Test_createTLSCert/Failed_to_create_due_to_empty_keyFile (0.00s)
    --- PASS: Test_createTLSCert/Creates_TLS_cert_successfully (0.00s)
    --- PASS: Test_createTLSCert/Failed_to_create_due_to_empty_port (0.00s)
    
=== RUN   TestListUsersHandler
=== RUN   TestListUsersHandler/Failed_to_get_users_list_due_to_incorrect_URL
The users map is empty
=== RUN   TestListUsersHandler/Gets_fail_due_to_empty_users_map
The users map is empty
=== RUN   TestListUsersHandler/Gets_users_list_successfully_(if_there_are_users_in_the_folder)
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
--- PASS: TestListUsersHandler (0.00s)
    --- PASS: TestListUsersHandler/Failed_to_get_users_list_due_to_incorrect_URL (0.00s)
    --- PASS: TestListUsersHandler/Gets_fail_due_to_empty_users_map (0.00s)
    --- PASS: TestListUsersHandler/Gets_users_list_successfully_(if_there_are_users_in_the_folder) (0.00s)
    
=== RUN   TestAddUserHandler
=== RUN   TestAddUserHandler/Adds_a_new_user_successfully
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
## bari@gmail.com added successfully:
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestAddUserHandler/Adds_fail_due_to_incorrect_user_(nil)
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestAddUserHandler/Adds_fail_due_to_incorrect_URL
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestAddUserHandler/Adds_fail_due_to_incorrect_empty_user
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
--- PASS: TestAddUserHandler (0.00s)
    --- PASS: TestAddUserHandler/Adds_a_new_user_successfully (0.00s)
    --- PASS: TestAddUserHandler/Adds_fail_due_to_incorrect_user_(nil) (0.00s)
    --- PASS: TestAddUserHandler/Adds_fail_due_to_incorrect_URL (0.00s)
    --- PASS: TestAddUserHandler/Adds_fail_due_to_incorrect_empty_user (0.00s)
    
=== RUN   TestGetUserHandler
=== RUN   TestGetUserHandler/Gets_fail_due_to_empty_email
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestGetUserHandler/Gets_an_existing_user_successfully
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
## bari@gmail.com gets successfully:
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestGetUserHandler/Gets_fail_due_to_incorrect_URL
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestGetUserHandler/Gets_fail_due_to_the_user_doesn't_exist
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
--- PASS: TestGetUserHandler (0.00s)
    --- PASS: TestGetUserHandler/Gets_fail_due_to_empty_email (0.00s)
    --- PASS: TestGetUserHandler/Gets_an_existing_user_successfully (0.00s)
    --- PASS: TestGetUserHandler/Gets_fail_due_to_incorrect_URL (0.00s)
    --- PASS: TestGetUserHandler/Gets_fail_due_to_the_user_doesn't_exist (0.00s)
    
=== RUN   TestUpdateUserHandler
=== RUN   TestUpdateUserHandler/Updates_fail_due_to_incorrect_user_(nil)
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestUpdateUserHandler/Updates_fail_due_to_incorrect_URL
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestUpdateUserHandler/Updates_fail_due_to_incorrect_empty_user
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestUpdateUserHandler/Updates_an_existing_user_successfully
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
## bari@gmail.com updated successfully:
Current users map:
* Email: bari@gmail.com, Username: bari2, Password: 12345
=== RUN   TestUpdateUserHandler/Updates_fail_due_to_incorrect_email
Current users map:
* Email: bari@gmail.com, Username: bari2, Password: 12345
--- PASS: TestUpdateUserHandler (0.00s)
    --- PASS: TestUpdateUserHandler/Updates_fail_due_to_incorrect_user_(nil) (0.00s)
    --- PASS: TestUpdateUserHandler/Updates_fail_due_to_incorrect_URL (0.00s)
    --- PASS: TestUpdateUserHandler/Updates_fail_due_to_incorrect_empty_user (0.00s)
    --- PASS: TestUpdateUserHandler/Updates_an_existing_user_successfully (0.00s)
    --- PASS: TestUpdateUserHandler/Updates_fail_due_to_incorrect_email (0.00s)
    
=== RUN   TestDeleteUserHandler
=== RUN   TestDeleteUserHandler/Deletes_fail_due_to_empty_email
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
=== RUN   TestDeleteUserHandler/Deletes_user_successfully
Current users map:
* Email: bari@gmail.com, Username: bari, Password: 1234
## bari@gmail.com deleted successfully:
The users map is empty
=== RUN   TestDeleteUserHandler/Deletes_fail_due_to_incorrect_URL
The users map is empty
=== RUN   TestDeleteUserHandler/Deletes_fail_due_to_the_user_doesn't_exist
The users map is empty
--- PASS: TestDeleteUserHandler (0.00s)
    --- PASS: TestDeleteUserHandler/Deletes_fail_due_to_empty_email (0.00s)
    --- PASS: TestDeleteUserHandler/Deletes_user_successfully (0.00s)
    --- PASS: TestDeleteUserHandler/Deletes_fail_due_to_incorrect_URL (0.00s)
    --- PASS: TestDeleteUserHandler/Deletes_fail_due_to_the_user_doesn't_exist (0.00s)
    
=== RUN   Test_getStatusAndMsgErr
=== RUN   Test_getStatusAndMsgErr/Gets_status_code_404
=== RUN   Test_getStatusAndMsgErr/Gets_status_code_500
--- PASS: Test_getStatusAndMsgErr (0.00s)
    --- PASS: Test_getStatusAndMsgErr/Gets_status_code_404 (0.00s)
    --- PASS: Test_getStatusAndMsgErr/Gets_status_code_500 (0.00s)
    
=== RUN   Test_setupRouter
--- PASS: Test_setupRouter (0.00s)
PASS

coverage: 64.6% of statements in ./...
```

<p align="center">
    <img src="coverage.png" width="350"/>
</p>