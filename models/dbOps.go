package models

type DBOps interface {
	GetAllUsers() ([]User, error)
	DeleteUser(email string) error
	InsertNewUser(user User) error
	UpdateNameAndPassUser(user User) error
	IsExistsInUsersTable(email string) (*User, error)
}
