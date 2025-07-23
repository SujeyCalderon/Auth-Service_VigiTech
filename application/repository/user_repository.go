package repository

import "auth-service/domain"

type UserRepository interface {
    GetUserByEmail(email string) (*domain.User, error)
    CreateUser(user *domain.User) error
    GetAllUsers() ([]*domain.User, error)
    UpdateUserStatus(id uint, active bool) error
}
