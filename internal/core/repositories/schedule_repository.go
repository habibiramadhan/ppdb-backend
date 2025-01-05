// internal/core/repositories/schedule_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "time"
)

type ScheduleRepository interface {
    Create(schedule *models.Schedule) error
    FindAll(limit, offset int) ([]models.Schedule, int64, error)
    FindByID(id uuid.UUID) (*models.Schedule, error)
    FindByAcademicYear(yearID uuid.UUID, limit, offset int) ([]models.Schedule, int64, error)
    FindUpcoming(limit int) ([]models.Schedule, error)
    FindOverlapping(startDate, endDate time.Time, excludeID *uuid.UUID) ([]models.Schedule, error)
    Update(schedule *models.Schedule) error
    Delete(id uuid.UUID) error
    SetStatus(id uuid.UUID, isActive bool) error
}

type scheduleRepository struct {
    db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
    return &scheduleRepository{db}
}

func (r *scheduleRepository) Create(schedule *models.Schedule) error {
    return r.db.Create(schedule).Error
}

func (r *scheduleRepository) FindAll(limit, offset int) ([]models.Schedule, int64, error) {
    var schedules []models.Schedule
    var total int64

    if err := r.db.Model(&models.Schedule{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := r.db.Preload("AcademicYear").
        Preload("Creator").
        Preload("Updater").
        Limit(limit).
        Offset(offset).
        Order("start_date asc").
        Find(&schedules).Error

    return schedules, total, err
}

func (r *scheduleRepository) FindByID(id uuid.UUID) (*models.Schedule, error) {
    var schedule models.Schedule
    
    err := r.db.Preload("AcademicYear").
        Preload("Creator").
        Preload("Updater").
        First(&schedule, "id = ?", id).Error

    if err != nil {
        return nil, err
    }

    return &schedule, nil
}

func (r *scheduleRepository) FindByAcademicYear(yearID uuid.UUID, limit, offset int) ([]models.Schedule, int64, error) {
    var schedules []models.Schedule
    var total int64

    if err := r.db.Model(&models.Schedule{}).
        Where("academic_year_id = ?", yearID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := r.db.Preload("AcademicYear").
        Preload("Creator").
        Preload("Updater").
        Where("academic_year_id = ?", yearID).
        Limit(limit).
        Offset(offset).
        Order("start_date asc").
        Find(&schedules).Error

    return schedules, total, err
}

func (r *scheduleRepository) FindUpcoming(limit int) ([]models.Schedule, error) {
    var schedules []models.Schedule

    err := r.db.Preload("AcademicYear").
        Where("end_date > ? AND is_active = ?", time.Now(), true).
        Limit(limit).
        Order("start_date asc").
        Find(&schedules).Error

    return schedules, err
}

func (r *scheduleRepository) FindOverlapping(startDate, endDate time.Time, excludeID *uuid.UUID) ([]models.Schedule, error) {
    var schedules []models.Schedule
    query := r.db.Where(
        "(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
        endDate, startDate, endDate, startDate, startDate, endDate,
    )

    if excludeID != nil {
        query = query.Where("id != ?", *excludeID)
    }

    err := query.Find(&schedules).Error
    return schedules, err
}

func (r *scheduleRepository) Update(schedule *models.Schedule) error {
    return r.db.Save(schedule).Error
}

func (r *scheduleRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.Schedule{}, id).Error
}

func (r *scheduleRepository) SetStatus(id uuid.UUID, isActive bool) error {
    return r.db.Model(&models.Schedule{}).
        Where("id = ?", id).
        Update("is_active", isActive).Error
}