// internal/core/services/schedule_service.go
package services

import (
	"errors"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/models"
	"ppdb-backend/utils"
	"time"

	"github.com/google/uuid"
)

type ScheduleService interface {
	Create(input CreateScheduleInput, userID uuid.UUID) error
	GetAll(page, limit int) ([]models.Schedule, *utils.PaginationMeta, error)
	GetByID(id uuid.UUID) (*models.Schedule, error)
	GetByAcademicYear(yearID uuid.UUID, page, limit int) ([]models.Schedule, *utils.PaginationMeta, error)
	Update(id uuid.UUID, input UpdateScheduleInput, userID uuid.UUID) error
	Delete(id uuid.UUID) error
	SetStatus(id uuid.UUID, isActive bool) error
	GetUpcomingSchedules(limit int) ([]models.Schedule, error)
	CreateNotification(scheduleID uuid.UUID, userIDs []uuid.UUID) error
}

type scheduleService struct {
	scheduleRepo     repositories.ScheduleRepository
	notificationRepo repositories.ScheduleNotificationRepository
	academicYearRepo repositories.AcademicYearRepository
	emailService     EmailService
}

type CreateScheduleInput struct {
	AcademicYearID uuid.UUID `json:"academic_year_id" validate:"required"`
	Title          string    `json:"title" validate:"required"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	ScheduleType   string    `json:"schedule_type" validate:"required"`
	Priority       string    `json:"priority" validate:"required"`
	Location       string    `json:"location"`
	RemindBefore   *int      `json:"remind_before"`
}

type UpdateScheduleInput struct {
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	ScheduleType string    `json:"schedule_type" validate:"required"`
	Priority     string    `json:"priority" validate:"required"`
	Location     string    `json:"location"`
	RemindBefore *int      `json:"remind_before"`
}

func NewScheduleService(
	scheduleRepo repositories.ScheduleRepository,
	notificationRepo repositories.ScheduleNotificationRepository,
	academicYearRepo repositories.AcademicYearRepository,
	emailService EmailService,
) ScheduleService {
	return &scheduleService{
		scheduleRepo:     scheduleRepo,
		notificationRepo: notificationRepo,
		academicYearRepo: academicYearRepo,
		emailService:     emailService,
	}
}

func (s *scheduleService) Create(input CreateScheduleInput, userID uuid.UUID) error {
	academicYear, err := s.academicYearRepo.FindByID(input.AcademicYearID)
	if err != nil {
		return errors.New("tahun ajaran ga ketemu")
	}
	if !academicYear.IsActive {
		return errors.New("tahun ajaran ga aktif")
	}

	if !models.IsValidScheduleType(input.ScheduleType) {
		return errors.New("tipe jadwal ga valid")
	}

	if !models.IsValidPriorityLevel(input.Priority) {
		return errors.New("level prioritas ga valid")
	}

	overlapping, err := s.scheduleRepo.FindOverlapping(input.StartDate, input.EndDate, nil)
	if err != nil {
		return err
	}
	if len(overlapping) > 0 {
		return errors.New("ada jadwal lain di waktu yang sama")
	}

	schedule := &models.Schedule{
		AcademicYearID: input.AcademicYearID,
		Title:          input.Title,
		Description:    input.Description,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		ScheduleType:   models.ScheduleType(input.ScheduleType),
		Priority:       models.PriorityLevel(input.Priority),
		Location:       input.Location,
		RemindBefore:   input.RemindBefore,
		CreatedBy:      userID,
		IsActive:       true,
	}

	return s.scheduleRepo.Create(schedule)
}

func (s *scheduleService) GetAll(page, limit int) ([]models.Schedule, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit

	schedules, total, err := s.scheduleRepo.FindAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(total) + limit - 1) / limit
	pagination := &utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPages,
	}

	return schedules, pagination, nil
}

func (s *scheduleService) GetByID(id uuid.UUID) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("jadwal ga ketemu")
	}
	return schedule, nil
}

func (s *scheduleService) GetByAcademicYear(yearID uuid.UUID, page, limit int) ([]models.Schedule, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit

	schedules, total, err := s.scheduleRepo.FindByAcademicYear(yearID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(total) + limit - 1) / limit
	pagination := &utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPages,
	}

	return schedules, pagination, nil
}

func (s *scheduleService) Update(id uuid.UUID, input UpdateScheduleInput, userID uuid.UUID) error {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("jadwal ga ketemu")
	}

	if !models.IsValidScheduleType(input.ScheduleType) {
		return errors.New("tipe jadwal ga valid")
	}

	if !models.IsValidPriorityLevel(input.Priority) {
		return errors.New("level prioritas ga valid")
	}

	overlapping, err := s.scheduleRepo.FindOverlapping(input.StartDate, input.EndDate, &id)
	if err != nil {
		return err
	}
	if len(overlapping) > 0 {
		return errors.New("ada jadwal lain di waktu yang sama")
	}

	schedule.Title = input.Title
	schedule.Description = input.Description
	schedule.StartDate = input.StartDate
	schedule.EndDate = input.EndDate
	schedule.ScheduleType = models.ScheduleType(input.ScheduleType)
	schedule.Priority = models.PriorityLevel(input.Priority)
	schedule.Location = input.Location
	schedule.RemindBefore = input.RemindBefore
	schedule.UpdatedBy = &userID

	return s.scheduleRepo.Update(schedule)
}

func (s *scheduleService) Delete(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("jadwal ga ketemu")
	}

	if schedule.HasStarted() {
		return errors.New("ga bisa hapus jadwal yang udah mulai")
	}

	if err := s.notificationRepo.DeleteBySchedule(id); err != nil {
		return err
	}

	return s.scheduleRepo.Delete(id)
}

func (s *scheduleService) SetStatus(id uuid.UUID, isActive bool) error {
	_, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("jadwal ga ketemu")
	}
	return s.scheduleRepo.SetStatus(id, isActive)
}

func (s *scheduleService) GetUpcomingSchedules(limit int) ([]models.Schedule, error) {
	return s.scheduleRepo.FindUpcoming(limit)
}

func (s *scheduleService) CreateNotification(scheduleID uuid.UUID, userIDs []uuid.UUID) error {
	_, err := s.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		return errors.New("jadwal ga ketemu")
	}

	notifications := make([]*models.ScheduleNotification, 0, len(userIDs))
	for _, userID := range userIDs {
		notification := &models.ScheduleNotification{
			ScheduleID: scheduleID,
			UserID:     userID,
			Type:       "email",
			Status:     models.NotificationPending,
		}
		notifications = append(notifications, notification)
	}

	return s.notificationRepo.CreateBatch(notifications)
}
