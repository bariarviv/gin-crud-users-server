package db

import (
	"gin_CRUD_server/models"
)

type SqlOps struct {
	Name string
}

const (
	DeleteUserQuery    = `DELETE FROM users WHERE email=$1`
	IsExistsUserQuery  = `SELECT * FROM users WHERE email=$1`
	GetAllUsersQuery   = `SELECT * FROM users ORDER BY email DESC`
	UpdateUserQuery    = `UPDATE users SET username=$1, password=$2 WHERE email=$3`
	InsertNewUserQuery = `INSERT INTO users ("email", "username", "password") VALUES ($1, $2, $3)`
)

// GetAllUsers gets a list of all the users
func (DB SqlOps) GetAllUsers() ([]models.User, error) {
	var user models.User
	var users []models.User

	rows, err := Instance.Db.Query(GetAllUsersQuery)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		if err = rows.Scan(&user.Email, &user.Name, &user.Password, &user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// DeleteUser deletes an existing user in the users table
func (DB SqlOps) DeleteUser(email string) error {
	if _, err := Instance.Db.Exec(DeleteUserQuery, email); err != nil {
		return err
	}
	return nil
}

// InsertNewUser inserts a new user into the users table
func (DB SqlOps) InsertNewUser(user models.User) error {
	if _, err := Instance.Db.Exec(InsertNewUserQuery, user.Email, user.Name, user.Password); err != nil {
		return err
	}
	return nil
}

// UpdateNameAndPassUser updates the name and pass for an existing user in the users table
func (DB SqlOps) UpdateNameAndPassUser(user models.User) error {
	if _, err := Instance.Db.Exec(UpdateUserQuery, user.Name, user.Password, user.Email); err != nil {
		return err
	}
	return nil
}

// IsExistsInUsersTable checks if the usr exists in the users table
func (DB SqlOps) IsExistsInUsersTable(email string) (*models.User, error) {
	var user models.User
	row := Instance.Db.QueryRow(IsExistsUserQuery, email)
	err := row.Scan(&user.Email, &user.Name, &user.Password, &user.CreatedAt)
	return &user, err
}
