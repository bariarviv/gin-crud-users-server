package db

import (
	"database/sql"
	"fmt"

	"gin_CRUD_server/models"
)

type TestMapOps struct {
	Name  string
	Users map[string]models.User
}

// GetAllUsers gets a list of all the users
func (DB TestMapOps) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if len(DB.Users) == 0 {
		return users, fmt.Errorf("there are no users")
	}
	for _, user := range DB.Users {
		users = append(users, user)
	}
	return users, nil
}

// DeleteUser deletes an existing user in the users map
func (DB TestMapOps) DeleteUser(email string) error {
	if _, ok := DB.Users[email]; !ok {
		return sql.ErrNoRows
	} else {
		delete(DB.Users, email)
	}
	return nil
}

// InsertNewUser inserts a new user into the users map
func (DB TestMapOps) InsertNewUser(user models.User) error {
	if user.Email != "" {
		DB.Users[user.Email] = user
	}
	return nil
}

// UpdateNameAndPassUser updates the name and pass for an existing user in the users map
func (DB TestMapOps) UpdateNameAndPassUser(user models.User) error {
	if val, ok := DB.Users[user.Email]; !ok {
		return sql.ErrNoRows
	} else {
		val.Name = user.Name
		val.Password = user.Password
		DB.Users[user.Email] = val
	}
	return nil
}

// IsExistsInUsersTable checks if the usr exists in the users map
func (DB TestMapOps) IsExistsInUsersTable(email string) (*models.User, error) {
	var user models.User
	if val, ok := DB.Users[email]; !ok {
		return &user, sql.ErrNoRows
	} else {
		return &val, nil
	}
}
