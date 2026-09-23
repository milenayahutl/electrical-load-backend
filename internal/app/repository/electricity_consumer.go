package repository

import (
	"electricity_consumers/internal/app/ds"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// получение одной опубликованной услуги для ленты
func (r *Repository) GetElectricityConsumer(
	id int,
) (ds.ElectricityConsumer, error) {

	var electricityConsumer ds.ElectricityConsumer

	err := r.db.
		Preload("Likes").
		Where(
			"id = ? AND status = ?",
			id,
			ds.StatusPublished,
		).
		First(&electricityConsumer).Error

	if err != nil {
		return ds.ElectricityConsumer{}, err
	}

	return electricityConsumer, nil
}

// получение опубликованных услуг
func (r *Repository) GetPublishedElectricityConsumers() (
	[]ds.ElectricityConsumer,
	error,
) {

	var electricityConsumers []ds.ElectricityConsumer

	err := r.db.
		Preload("Likes").
		Where(
			"status = ?",
			ds.StatusPublished,
		).
		Order("id ASC").
		Find(&electricityConsumers).Error

	if err != nil {
		return nil, err
	}

	return electricityConsumers, nil
}

// получение черновика конкретного пользователя
func (r *Repository) GetDraftElectricityConsumer(
	creatorID uint,
) (*ds.ElectricityConsumer, error) {

	var draft ds.ElectricityConsumer

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			ds.StatusDraft,
		).
		First(&draft).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &draft, nil
}

// поиск опубликованных услуг по мощности
func (r *Repository) GetElectricityConsumersByPowerRange(
	minPower float64,
	maxPower float64,
) ([]ds.ElectricityConsumer, error) {

	var electricityConsumers []ds.ElectricityConsumer

	err := r.db.
		Preload("Likes").
		Where(
			"power_kw >= ? AND power_kw <= ? AND status = ?",
			minPower,
			maxPower,
			ds.StatusPublished,
		).
		Order("id ASC").
		Find(&electricityConsumers).Error

	if err != nil {
		return nil, err
	}

	return electricityConsumers, nil
}

// получение следующей опубликованной услуги
func (r *Repository) GetNextElectricityConsumerID(
	currentID int,
) (int, error) {

	var next ds.ElectricityConsumer

	err := r.db.
		Where(
			"id > ? AND status = ?",
			currentID,
			ds.StatusPublished,
		).
		Order("id ASC").
		First(&next).Error

	if err == nil {
		return int(next.ID), nil
	}

	// Если дошли до конца — возвращаемся к первой
	var first ds.ElectricityConsumer

	err = r.db.
		Where(
			"status = ?",
			ds.StatusPublished,
		).
		Order("id ASC").
		First(&first).Error

	if err != nil {
		return 0, fmt.Errorf(
			"опубликованные услуги не найдены",
		)
	}

	return int(first.ID), nil
}

// создание новой услуги-черновика через ORM
func (r *Repository) CreateDraft(
	creatorID uint,
	name string,
) (ds.ElectricityConsumer, error) {

	// нет ли уже черновика
	var oldDraft ds.ElectricityConsumer

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			ds.StatusDraft,
		).
		First(&oldDraft).Error

	if err == nil {
		return ds.ElectricityConsumer{},
			fmt.Errorf("у пользователя уже есть черновик")
	}

	if err != gorm.ErrRecordNotFound {
		return ds.ElectricityConsumer{}, err
	}

	draft := ds.ElectricityConsumer{
		Name:       name,
		Status:     ds.StatusDraft,
		Image:      "",
		Video:      "",
		PowerKW:    0,
		CurrentA:   0,
		DateCreate: time.Now(),
		CreatorID:  creatorID,
	}

	err = r.db.Create(&draft).Error
	if err != nil {
		return ds.ElectricityConsumer{}, err
	}

	return draft, nil
}

// публикация черновика через ORM
func (r *Repository) PublishDraft(
	id uint,
	creatorID uint,
	description string,
	powerKW float64,
) error {
	currentA := powerKW * 1000 / 220
	currentA = math.Round(currentA*10) / 10

	result := r.db.
		Model(&ds.ElectricityConsumer{}).
		Where(
			"id = ? AND creator_id = ? AND status = ?",
			id,
			creatorID,
			ds.StatusDraft,
		).
		Updates(map[string]interface{}{
			"description":    description,
			"power_kw":       powerKW,
			"current_a":      currentA,
			"status":         ds.StatusPublished,
			"date_formation": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("черновик не найден")
	}

	return nil
}

// логическое удаление через ручной SQL UPDATE
func (r *Repository) DeleteElectricityConsumer(
	id uint,
) error {

	query := `
		UPDATE electricity_consumers
		SET status = ?
		WHERE id = ?
		  AND status = ?
		RETURNING id
	`

	row := r.db.Raw(
		query,
		ds.StatusDeleted,
		id,
		ds.StatusPublished,
	).Row()

	var deletedID uint

	err := row.Scan(&deletedID)
	if err != nil {
		return fmt.Errorf("услуга не найдена: %w", err)
	}

	return nil
}
