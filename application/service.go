package application

import (
    "errors"
    "strconv"
    "time"

    "auth-service/application/repository"
    "auth-service/domain"
    "auth-service/infrastructure/token"
    "golang.org/x/crypto/bcrypt"
)

const FixedSystemID = 1255

type Service struct {
    UserRepo     repository.UserRepository
    TokenService token.JWTService
}

func NewService(userRepo repository.UserRepository, tokenService token.JWTService) *Service {
    return &Service{userRepo, tokenService}
}

func hashPassword(password string) (string, error) {
    h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(h), err
}

func checkPassword(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}


func (s *Service) RegisterUser(
    name, email, password, role string,
    systemID *int,
) (*domain.User, error) {

    if existing, _ := s.UserRepo.GetUserByEmail(email); existing != nil {
        return nil, domain.ErrEmailAlreadyExists
    }

    if role == "user" {
        if systemID == nil || *systemID != FixedSystemID {
            return nil, domain.ErrUnauthorized     
        }
    }

    hashed, err := hashPassword(password)
    if err != nil {
        return nil, err
    }

    assignedID := 0
    if systemID != nil {
        assignedID = *systemID        
    }

    user := &domain.User{
        Name:     name,
        Email:    email,
        Password: hashed,
        Active:   true,
        Role:     role,
        SystemID: assignedID,
    }
    if err := s.UserRepo.CreateUser(user); err != nil {
        return nil, err
    }
    return user, nil
}


func (s *Service) Login(email, password string) (string, *domain.User, error) {

    user, err := s.UserRepo.GetUserByEmail(email)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            return "", nil, domain.ErrInvalidCredentials
        }
        return "", nil, err
    }

    if !user.Active {
        return "", nil, errors.New("cuenta desactivada")
    }
    if !checkPassword(password, user.Password) {
        return "", nil, domain.ErrInvalidCredentials
    }

    if user.Role == "user" && user.SystemID != FixedSystemID {
        return "", nil, domain.ErrUnauthorized
    }

    tokenStr, err := s.TokenService.GenerateToken(
        strconv.Itoa(int(user.ID)), 24*time.Hour, user.Role)
    return tokenStr, user, err
}
func (s *Service) GetAllUsers() ([]*domain.User, error) {
    return s.UserRepo.GetAllUsers()
}

func (s *Service) UpdateUserStatus(id uint, active bool) error {
    return s.UserRepo.UpdateUserStatus(id, active)
}
