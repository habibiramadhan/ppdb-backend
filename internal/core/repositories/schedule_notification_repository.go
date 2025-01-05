// internal/core/repositories/schedule_notification_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ScheduleNotificationRepository interface {
    Create(notification *models.ScheduleNotification) error
    CreateBatch(notifications []*models.ScheduleNotification) error
    FindBySchedule(scheduleID uuid.UUID, limit, offset int) ([]models.ScheduleNotification, int64, error)
    FindPendingNotifications() ([]models.ScheduleNotification, error)
    UpdateStatus(id uuid.UUID, status string, errorMsg *string) error
    DeleteBySchedule(scheduleID uuid.UUID) error
}

type scheduleNotificationRepository struct {
    db *gorm.DB
}

func NewScheduleNotificationRepository(db *gorm.DB) ScheduleNotificationRepository {
    return &scheduleNotificationRepository{db}
}

func (r *scheduleNotificationRepository) Create(notification *models.ScheduleNotification) error {
    return r.db.Create(notification).Error
}

func (r *scheduleNotificationRepository) CreateBatch(notifications []*models.ScheduleNotification) error {
    return r.db.Create(notifications).Error
}

func (r *scheduleNotificationRepository) FindBySchedule(scheduleID uuid.UUID, limit, offset int) ([]models.ScheduleNotification, int64, error) {
    var notifications []models.ScheduleNotification
    var total int64

    if err := r.db.Model(&models.ScheduleNotification{}).
        Where("schedule_id = ?", scheduleID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    err := r.db.Preload("Schedule").
        Preload("User").
        Where("schedule_id = ?", scheduleID).
        Limit(limit).
        Offset(offset).
        Order("created_at desc").
        Find(&notifications).Error

    return notifications, total, err
}

func (r *scheduleNotificationRepository) FindPendingNotifications() ([]models.ScheduleNotification, error) {
    var notifications []models.ScheduleNotification

    err := r.db.Preload("Schedule").
        Preload("User").
        Where("status = ?", models.NotificationPending).
        Find(&notifications).Error

    return notifications, err
}

func (r *scheduleNotificationRepository) UpdateStatus(id uuid.UUID, status string, errorMsg *string) error {
    updates := map[string]interface{}{
        "status": status,
    }

    if errorMsg != nil {
        updates["error_message"] = *errorMsg
    }

    return r.db.Model(&models.ScheduleNotification{}).
        Where("id = ?", id).
        Updates(updates).Error
}

func (r *scheduleNotificationRepository) DeleteBySchedule(scheduleID uuid.UUID) error {
    return r.db.Where("schedule_id = ?", scheduleID).
        Delete(&models.ScheduleNotification{}).Error
}