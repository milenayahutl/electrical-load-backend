package repository

import (
	"fmt"

	"electricity_consumers/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) CreateUser(
	login string,
	password string,
) (ds.User, error) {

	var oldUser ds.User

	err := r.db.
		Where("login = ?", login).
		First(&oldUser).Error

	if err == nil {
		return ds.User{},
			fmt.Errorf("пользователь с таким логином уже существует")
	}

	if err != gorm.ErrRecordNotFound {
		return ds.User{}, err
	}

	user := ds.User{
		Login:    login,
		Password: password,
	}

	err = r.db.Create(&user).Error

	if err != nil {
		return ds.User{}, err
	}

	return user, nil
}
