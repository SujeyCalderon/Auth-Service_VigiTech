package db

import (
	"database/sql"
	"log"

	"auth-service/application/repository"
	"auth-service/domain"
)

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	query := "INSERT INTO users (name, email, password, active, role, systemID) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := r.DB.Exec(query,
		user.Name, user.Email, user.Password,
		user.Active, user.Role, user.SystemID,
	)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint(id)
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := "SELECT id, name, email, password, active, role, systemID FROM users WHERE email = ?"
	row := r.DB.QueryRow(query, email)

	var u domain.User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.Password,
		&u.Active, &u.Role, &u.SystemID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		log.Printf("Error al obtener usuario: %v", err)
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetAllUsers() ([]*domain.User, error) {
	query := "SELECT id, name, email, password, active, role, systemID FROM users WHERE role != 'admin'"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Password,
			&u.Active, &u.Role, &u.SystemID,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *userRepository) UpdateUserStatus(id uint, active bool) error {
	_, err := r.DB.Exec("UPDATE users SET active = ? WHERE id = ?", active, id)
	return err
}
