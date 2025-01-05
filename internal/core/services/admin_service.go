// internal/core/services/admin_service.go
package services
import (
	"errors"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/models"
	"ppdb-backend/utils"

	"github.com/google/uuid"
)

type AdminService interface {
	GetAllUsers(page, limit int) ([]models.User, *utils.PaginationMeta, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, input UpdateUserInput) error
	DeleteUser(id uuid.UUID) error
	UpdateUserStatus(id uuid.UUID, status string) error
}

type adminService struct {
	userRepository repositories.UserRepository
}

type UpdateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin student"`
	Phone string `json:"phone" validate:"omitempty"`
}

type UpdateStatusInput struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}

func NewAdminService(userRepo repositories.UserRepository) AdminService {
	return &adminService{
		userRepository: userRepo,
	}
}

func (s *adminService) GetAllUsers(page, limit int) ([]models.User, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit
	users, totalCount, err := s.userRepository.FindAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(totalCount) + limit - 1) / limit
	pagination := &utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalCount,
		TotalPage: totalPages,
	}

	return users, pagination, nil
}

func (s *adminService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *adminService) UpdateUser(id uuid.UUID, input UpdateUserInput) error {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Email != input.Email {
		existingUser, _ := s.userRepository.FindByEmail(input.Email)
		if existingUser != nil {
			return errors.New("email already in use")
		}
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Role = input.Role
	user.Phone = input.Phone

	return s.userRepository.Update(user)
}

func (s *adminService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepository.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepository.Delete(id)
}

func (s *adminService) UpdateUserStatus(id uuid.UUID, status string) error {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.Status = status
	return s.userRepository.Update(user)
}
