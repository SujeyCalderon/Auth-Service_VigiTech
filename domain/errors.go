package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrEmailAlreadyExists = errors.New("el email ya está registrado")
	ErrUnauthorized       = errors.New("no autorizado")
)
