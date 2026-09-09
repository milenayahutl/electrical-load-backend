package repository

import (
	"fmt"
	"math"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// Service — одна услуга, то есть один электроприбор.
type Service struct {
	ID          int
	Name        string
	PowerKW     float64
	Description string
	Image       string
	Status      string
	Video       string
	LikedBy     []int
}

const (
	StatusDraft     = "черновик"
	StatusPublished = "опубликован"
	StatusDeleted   = "удален"
)

// CurrentA рассчитывает нагрузку на сеть 220 В
// I = P / U
// Мощность переводим из кВт в Вт.
func (s Service) CurrentA() float64 {
	current := s.PowerKW * 1000 / 220

	return math.Round(current*10) / 10
}

// GetServices возвращает все услуги.
func (r *Repository) GetServices() ([]Service, error) {
	const media = "http://localhost:9000/media/"

	services := []Service{
		{
			ID:      1,
			Name:    "Электрический чайник",
			PowerKW: 2.0,
			Description: "Быстро нагревает воду для чая и других горячих напитков. " +
				"Подходит для ежедневного использования дома.",
			Image:   media + "kettle.jpg",
			Video:   media + "kettle.mp4",
			Status:  StatusPublished,
			LikedBy: []int{1, 3, 5, 8},
		},
		{
			ID:      2,
			Name:    "Тостер",
			PowerKW: 0.9,
			Description: "Поджаривает ломтики хлеба до хрустящей корочки. " +
				"Позволяет быстро приготовить горячие тосты к завтраку.",
			Image:   media + "toaster.jpg",
			Video:   media + "toaster.mp4",
			Status:  StatusPublished,
			LikedBy: []int{2, 7},
		},
		{
			ID:      3,
			Name:    "Стиральная машина",
			PowerKW: 2.2,
			Description: "Автоматически стирает одежду и другие текстильные изделия. " +
				"Во время нагрева воды создаёт заметную нагрузку на электросеть.",
			Image:   media + "washing-machine.jpg",
			Video:   media + "washing-machine.mp4",
			Status:  StatusPublished,
			LikedBy: []int{1, 2, 4, 6, 9},
		},
		{
			ID:      4,
			Name:    "Духовой шкаф",
			PowerKW: 3.0,
			Description: "Используется для запекания, выпечки и приготовления горячих блюд. " +
				"Относится к мощным бытовым потребителям электроэнергии.",
			Image:   media + "oven.jpg",
			Video:   media + "oven.mp4",
			Status:  StatusPublished,
			LikedBy: []int{3, 5, 7},
		},
		{
			ID:      5,
			Name:    "Пылесос",
			PowerKW: 1.6,
			Description: "Удаляет пыль и загрязнения с пола и других поверхностей. " +
				"Используется для регулярной уборки помещений.",
			Image:   media + "vacuum.jpg",
			Video:   media + "vacuum.mp4",
			Status:  StatusPublished,
			LikedBy: []int{2, 6, 8},
		},
		{
			ID:      6,
			Name:    "Блендер",
			PowerKW: 0.8,
			Description: "Измельчает и смешивает продукты для напитков и блюд. " +
				"Имеет сравнительно небольшую мощность.",
			Image:   media + "blender.jpg",
			Video:   media + "blender.mp4",
			Status:  StatusDraft,
			LikedBy: []int{1, 9},
		},
		{
			ID:      7,
			Name:    "Утюг",
			PowerKW: 2.0,
			Description: "Разглаживает складки на одежде с помощью нагрева и пара. " +
				"При работе потребляет значительную мощность.",
			Image:   media + "iron.jpg",
			Video:   media + "iron.mp4",
			Status:  StatusDeleted,
			LikedBy: []int{2, 4, 5, 8},
		},
		{
			ID:      8,
			Name:    "Микроволновая печь",
			PowerKW: 1.2,
			Description: "Быстро разогревает и готовит пищу с помощью микроволн. " +
				"Во время работы создаёт дополнительную нагрузку на домашнюю сеть.",
			Image:   media + "microwave.jpg",
			Video:   media + "microwave.mp4",
			Status:  StatusPublished,
			LikedBy: []int{1, 3, 4, 6, 7},
		},
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("массив услуг пустой")
	}

	return services, nil
}

func (r *Repository) GetService(id int) (Service, error) {
	services, err := r.GetPublishedServices()
	if err != nil {
		return Service{}, err
	}

	for _, service := range services {
		if service.ID == id {
			return service, nil
		}
	}

	return Service{}, fmt.Errorf("услуга не найдена")
}

// возвращает только опубликованные услуги
func (r *Repository) GetPublishedServices() ([]Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return nil, err
	}

	var result []Service

	for _, service := range services {
		if service.Status == StatusPublished {
			result = append(result, service)
		}
	}

	return result, nil
}

// GetDraftService возвращает услугу в статусе "черновик".
func (r *Repository) GetDraftService() (Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return Service{}, err
	}

	for _, service := range services {
		if service.Status == StatusDraft {
			return service, nil
		}
	}

	return Service{}, fmt.Errorf("черновик не найден")
}

// Показываем приборы с мощностью НЕ БОЛЬШЕ введённой.
func (r *Repository) GetServicesByMaxPower(maxPower float64) ([]Service, error) {
	services, err := r.GetPublishedServices()
	if err != nil {
		return nil, err
	}

	var result []Service

	for _, service := range services {
		if service.PowerKW <= maxPower {
			result = append(result, service)
		}
	}

	return result, nil
}

func (r *Repository) GetNextServiceID(currentID int) (int, error) {
	services, err := r.GetPublishedServices()
	if err != nil {
		return 0, err
	}

	if len(services) == 0 {
		return 0, fmt.Errorf("нет опубликованных услуг")
	}

	currentExists := false

	minID := services[0].ID
	nextID := -1

	for _, service := range services {

		if service.ID == currentID {
			currentExists = true
		}

		if service.ID < minID {
			minID = service.ID
		}

		// Ищем минимальный существующий ID,
		// который БОЛЬШЕ текущего.
		if service.ID > currentID {

			if nextID == -1 || service.ID < nextID {
				nextID = service.ID
			}
		}
	}

	if !currentExists {
		return 0, fmt.Errorf("услуга не найдена")
	}

	// Если текущая карточка последняя —
	// возвращаемся к первой опубликованной.
	if nextID == -1 {
		return minID, nil
	}

	return nextID, nil
}
