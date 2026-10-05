package ds

type ElectricityConsumerResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Video       string  `json:"video"`
	PowerKW     float64 `json:"power_kw"`
	CurrentA    float64 `json:"current_a"`

	LikesCount int `json:"likes_count"`

	IsCreator int `json:"is_creator"`
	IsLiked   int `json:"is_liked"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
}

type PublishElectricityConsumerRequest struct {
	Description string  `json:"description"`
	PowerKW     float64 `json:"power_kw"`
}

type LikeRequest struct {
	Value *int `json:"value"`
}

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}


func ToElectricityConsumerResponse(
	consumer ElectricityConsumer,
	currentUserID uint,
) ElectricityConsumerResponse {

	isCreator := 0

	if consumer.CreatorID == currentUserID {
		isCreator = 1
	}

	isLiked := 0

	for _, like := range consumer.Likes {
		if like.UserID == currentUserID {
			isLiked = 1
			break
		}
	}

	return ElectricityConsumerResponse{
		ID:          consumer.ID,
		Name:        consumer.Name,
		Description: consumer.Description,
		Image:       consumer.Image,
		Video:       consumer.Video,
		PowerKW:     consumer.PowerKW,
		CurrentA:    consumer.CurrentA,
		LikesCount:  len(consumer.Likes),
		IsCreator:   isCreator,
		IsLiked:     isLiked,
	}
}

func ToElectricityConsumerResponseList(
	consumers []ElectricityConsumer,
	currentUserID uint,
) []ElectricityConsumerResponse {

	result := make(
		[]ElectricityConsumerResponse,
		0,
		len(consumers),
	)

	for _, consumer := range consumers {
		result = append(
			result,
			ToElectricityConsumerResponse(
				consumer,
				currentUserID,
			),
		)
	}

	return result
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Login: user.Login,
	}
}
