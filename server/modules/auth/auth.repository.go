package auth

import (
	"errors"
	"strings"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
}

type AuthRepositoryStruct struct {
}

func NewAuthRepository() AuthRepository {
	return &AuthRepositoryStruct{}
}

func (a *AuthRepositoryStruct) FindUserByEmail(email string) (*model.User, error) {
	var user model.User

	if err := config.Database.First(&user, "email = ?", strings.ToLower(email)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (a *AuthRepositoryStruct) Create(user *model.User) error {
	result := config.Database.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
