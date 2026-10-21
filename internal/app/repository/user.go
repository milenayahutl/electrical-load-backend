
package repository

import (
	"errors"
	"fmt"

	"electricity_consumers/internal/app/ds"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// создание нового пользователя
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

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.User{}, err
	}

	// хешируем пароль перед сохранением
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return ds.User{}, err
	}

	user := ds.User{
		Login:    login,
		Password: string(passwordHash),
	}

	err = r.db.Create(&user).Error

	if err != nil {
		return ds.User{}, err
	}

	return user, nil
}

// поиск по логину
func (r *Repository) GetUserByLogin(
	login string,
) (ds.User, error) {

	var user ds.User

	err := r.db.
		Where("login = ?", login).
		First(&user).Error

	if err != nil {
		return ds.User{}, err
	}

	return user, nil
}
