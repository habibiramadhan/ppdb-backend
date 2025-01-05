//internal/core/services/auth_service.go
package services
import (
    "crypto/rand"
    "encoding/hex"
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
    userRepo         repositories.UserRepository
    verificationRepo repositories.VerificationRepository
    emailService    EmailService
    jwtSecret       string
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

func NewAuthService(
    userRepo repositories.UserRepository,
    verificationRepo repositories.VerificationRepository,
    emailService EmailService,
    jwtSecret string,
) AuthService {
    return &authService{
        userRepo:         userRepo,
        verificationRepo: verificationRepo,
        emailService:     emailService,
        jwtSecret:        jwtSecret,
    }
}

func generateToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func (s *authService) Register(input RegisterInput) error {
    existingUser, _ := s.userRepo.FindByEmail(input.Email)
    if existingUser != nil {
        return errors.New("email already registered")
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return err
    }

    user := &models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: hashedPassword,
        Role:     input.Role,
        Status:   "inactive",
    }

    if err := s.userRepo.Create(user); err != nil {
        return err
    }

    token, err := generateToken()
    if err != nil {
        return err
    }

    verification := &models.EmailVerification{
        UserID:    user.ID,
        Token:     token,
        SentAt:    time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    if err := s.verificationRepo.Create(verification); err != nil {
        return err
    }

    if err := s.emailService.SendVerificationEmail(user.Email, token, user.Name); err != nil {
        return err
    }

    return nil
}

func (s *authService) Login(input LoginInput) (string, error) {
    user, err := s.userRepo.FindByEmail(input.Email)
    if err != nil {
        return "", errors.New("invalid email or password")
    }

    if user.Status != "active" {
        return "", errors.New("please verify your email first")
    }

    if !utils.CheckPasswordHash(input.Password, user.Password) {
        return "", errors.New("invalid email or password")
    }

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