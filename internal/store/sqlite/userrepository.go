// internal/store/sqlite/userrepository.go
package sqlite

import (
	"Forum/internal/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	store *sqliteStore
}

func (r *UserRepository) ExistingUser(userName, email string) error {
	queryEmail := "SELECT * FROM users WHERE email = ?"
	rows, err := r.store.db.Query(queryEmail, email)
	if err != nil {
		return fmt.Errorf("email check failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("email already in use")
	}

	queryName := "SELECT * FROM users WHERE username = ?"
	rows, err = r.store.db.Query(queryName, userName)
	if err != nil {
		return fmt.Errorf("user name check failed: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("username already in use")
	}

	return nil
}

func (r *UserRepository) Login(user *model.User) error {
	var hashedPassword string
	err := r.store.db.QueryRow("SELECT UUID, email, username, password FROM users WHERE email = ?", user.Email).Scan(&user.UUID, &user.Email, &user.UserName, &hashedPassword)
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
	queryInsert := "INSERT INTO users(UUID, email, username, password, firstname, lastname, age, gender) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.store.db.Exec(queryInsert, user.UUID, user.Email, user.UserName, user.Password, user.FirstName, user.LastName, user.Age, user.Gender)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByUUID(uuid string) (*model.User, error) {
	var u model.User
	if err := r.store.db.QueryRow(
		"SELECT UUID, username, email, password, firstname, lastname, age, gender FROM users WHERE UUID = ?",
		uuid,
	).Scan(&u.UUID, &u.UserName, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Age, &u.Gender); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetUsersOrderedByLastMessageOrAlphabetically() ([]*model.User, error) {
	query := `
        SELECT users.user_id, users.username, users.email, users.first_name, users.last_name, users.age, users.gender, MAX(messages.created_at) AS last_message_time
        FROM users
        LEFT JOIN messages ON users.user_id = messages.sender_id OR users.user_id = messages.receiver_id
        GROUP BY users.user_id
        ORDER BY last_message_time DESC, users.username ASC
    `
	rows, err := r.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UUID, &user.UserName, &user.Email, &user.FirstName, &user.LastName, &user.Age, &user.Gender); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
