package repository

import "fmt"

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// Service - одна услуга(один электроприбор)
type Service struct {
	ID          int
	Name        string
	PowerKW     float64
	Description string
	Image       string
	Video       string

	// айди пользователей, поставивших лайк
	LikedBy []int
}

// возвращает все электроприборы
func (r *Repository) GetServices() ([]Service, error) {
	services := []Service{
		{
			ID:      1,
			Name:    "Электрический чайник",
			PowerKW: 2.0,
			Description: "Быстро нагревает воду для чая и других горячих напитков. " +
				"Подходит для ежедневного использования дома.",
			Image:   "/static/img/kettle.jpg",
			Video:   "/static/video/kettle.mp4",
			LikedBy: []int{1, 3, 5, 8},
		},
		{
			ID:      2,
			Name:    "Тостер",
			PowerKW: 0.9,
			Description: "Поджаривает ломтики хлеба до хрустящей корочки. " +
				"Позволяет быстро приготовить горячие тосты к завтраку.",
			Image:   "/static/img/toaster.jpg",
			Video:   "/static/video/toaster.mp4",
			LikedBy: []int{2, 7},
		},
		{
			ID:      3,
			Name:    "Стиральная машина",
			PowerKW: 2.2,
			Description: "Автоматически стирает одежду и другие текстильные изделия. " +
				"Во время нагрева воды создаёт заметную нагрузку на электросеть.",
			Image:   "/static/img/washing-machine.jpg",
			Video:   "/static/video/washing-machine.mp4",
			LikedBy: []int{1, 2, 4, 6, 9},
		},
		{
			ID:      4,
			Name:    "Духовой шкаф",
			PowerKW: 3.0,
			Description: "Используется для запекания, выпечки и приготовления горячих блюд. " +
				"Относится к мощным бытовым потребителям электроэнергии.",
			Image:   "/static/img/oven.jpg",
			Video:   "/static/video/oven.mp4",
			LikedBy: []int{3, 5, 7},
		},
		{
			ID:      5,
			Name:    "Пылесос",
			PowerKW: 1.6,
			Description: "Удаляет пыль и загрязнения с пола и других поверхностей. " +
				"Используется для регулярной уборки помещений.",
			Image:   "/static/img/vacuum.jpg",
			Video:   "/static/video/vacuum.mp4",
			LikedBy: []int{2, 6, 8},
		},
		{
			ID:      6,
			Name:    "Блендер",
			PowerKW: 0.8,
			Description: "Измельчает и смешивает продукты для напитков и блюд. " +
				"Имеет сравнительно небольшую мощность среди бытовых приборов.",
			Image:   "/static/img/blender.jpg",
			Video:   "/static/video/blender.mp4",
			LikedBy: []int{1, 9},
		},
		{
			ID:      7,
			Name:    "Утюг",
			PowerKW: 2.0,
			Description: "Разглаживает складки на одежде с помощью нагрева и пара. " +
				"При работе нагревательного элемента потребляет значительную мощность.",
			Image:   "/static/img/iron.jpg",
			Video:   "/static/video/iron.mp4",
			LikedBy: []int{2, 4, 5, 8},
		},
		{
			ID:      8,
			Name:    "Микроволновая печь",
			PowerKW: 1.2,
			Description: "Быстро разогревает и готовит пищу с помощью микроволн. " +
				"Во время работы создаёт дополнительную нагрузку на домашнюю сеть.",
			Image:   "/static/img/microwave.jpg",
			Video:   "/static/video/microwave.mp4",
			LikedBy: []int{1, 3, 4, 6, 7},
		},
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("массив услуг пустой")
	}

	return services, nil
}

// фильтрация: показываем приборы с мощностью >= введённой
func (r *Repository) GetServicesByMaxPower(maxPower float64) ([]Service, error) {
	services, err := r.GetServices()
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

// пользователь нажал на чайник с ID=1,
// тогда чайник становится первым элементом ленты,
// а после него идут остальные
func (r *Repository) GetServicesFromID(id int) ([]Service, error) {
	services, err := r.GetServices()
	if err != nil {
		return nil, err
	}

	startIndex := -1

	for i, service := range services {
		if service.ID == id {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return nil, fmt.Errorf("услуга не найдена")
	}

	// берутся элементы начиная с выбранного
	result := append([]Service{}, services[startIndex:]...)

	// чтобы лента оставалась циклической, добавляем начало
	result = append(result, services[:startIndex]...)

	return result, nil
}
