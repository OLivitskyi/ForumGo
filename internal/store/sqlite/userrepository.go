package sqlite

import (
	"Forum/internal/model"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) ExistingUser(userName, email string) error {
	queryEmail := "SELECT * FROM users WHERE email = ?"
	rows, err := r.store.Db.Query(queryEmail, email)
	if err != nil {
		return errors.Join(errors.New("email check failed"), err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("email already in use")
	}

	queryName := "SELECT * FROM users WHERE username = ?"
	rows, err = r.store.Db.Query(queryName, userName)
	if err != nil {
		return errors.Join(errors.New("user name check failed"), err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("username already in use")
	}

	return nil
}

func (r *UserRepository) Login(user *model.User) error {
	var hashedPassword string
	err := r.store.Db.QueryRow("SELECT UUID, email, username, password FROM users WHERE email = ?", user.Email).Scan(&user.UUID, &user.Email, &user.UserName, &hashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Register(user *model.User) error {
	queryInsert := "INSERT INTO users(UUID, email, username, password) VALUES(?, ?, ?, ?) "
	_, err := r.store.Db.Exec(queryInsert, user.UUID, user.Email, user.UserName, user.Password)
	if err != nil {
		return err
	}
	return nil
}