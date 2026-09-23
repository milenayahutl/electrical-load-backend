package ds

type Like struct {
	ID uint `gorm:"primaryKey"`

	UserID                uint `gorm:"not null;uniqueIndex:idx_user_consumer"`
	ElectricityConsumerID uint `gorm:"not null;uniqueIndex:idx_user_consumer"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT;"`

	ElectricityConsumer ElectricityConsumer `gorm:"foreignKey:ElectricityConsumerID;constraint:OnDelete:RESTRICT;"`
}
