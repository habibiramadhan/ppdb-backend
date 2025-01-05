//internal/core/services/auth_service.go
package services

import (
    "errors"
    "ppdb-backend/internal/models"
    "ppdb-backend/internal/core/repositories"
    "ppdb-backend/utils"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
    Register(input RegisterInput) error
    Login(input LoginInput) (string, error)
    ValidateToken(tokenString string) (*jwt.Token, error)
}

type authService struct {
    userRepository repositories.UserRepository
    jwtSecret     string
}

type RegisterInput struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Role     string `json:"role" validate:"required,oneof=admin student"`
}

type LoginInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
    return &authService{
        userRepository: userRepo,
        jwtSecret:     jwtSecret,
    }
}

func (s *authService) Register(input RegisterInput) error {
    // Check if email already exists
    existingUser, _ := s.userRepository.FindByEmail(input.Email)
    if existingUser != nil {
        return errors.New("email already registered")
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return err
    }

    user := &models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: hashedPassword,
        Role:     input.Role,
        Status:   "active",
    }

    return s.userRepository.Create(user)
}

func (s *authService) Login(input LoginInput) (string, error) {
    user, err := s.userRepository.FindByEmail(input.Email)
    if err != nil {
        return "", errors.New("invalid email or password")
    }

    if !utils.CheckPasswordHash(input.Password, user.Password) {
        return "", errors.New("invalid email or password")
    }

    // Generate JWT Token
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.ID
    claims["email"] = user.Email
    claims["role"] = user.Role
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.jwtSecret), nil
    })
}